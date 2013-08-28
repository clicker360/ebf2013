// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package micrositio

import (
	"appengine"
	"appengine/datastore"
	"appengine/blobstore"
    "appengine/memcache"
	appimage "appengine/image"
	"crypto/sha1"
	"resize"
	"bytes"
	"strings"
	//"strconv"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // import so we can read PNG files.
	"io"
	"net/http"
	"text/template"
	"sess"
	"model"
)

type WsFormDataImage struct {
	Data	[]byte `json:"-"`
	IdEmp	string `json:"IdEmp"`
	IdImg	string `json:"IdImg"`
	Kind	string `json:"Kind"`
	Name	string `json:"Name"`
	Desc	string `json:"Desc"`
	Facebook string `json:"Facebook"`
	Twitter string `json:"Twitter"`
	Sizepx	int `json:"Sizepx"`
	Sizepy	int `json:"Sizepy"`
	Url		string `json:"Url"`
	Type	string `json:"Type,omitempty"`
	Sp1		string `json:"Sp1,omitempty"`
	Sp2		string `json:"Sp2,omitempty"`
	Sp3		string `json:"Sp3,omitempty"`
	Sp4		string `json:"Sp4,omitempty"`
	Np1		int `json:"Np1,omitempty"`
	Np2		int `json:"Np2,omitempty"`
	Np3		int `json:"Np3,omitempty"`
	Np4		int `json:"Np4,omitempty"`
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/r/wse/extra", wsMicrositio)
	http.HandleFunc("/r/wse/logoup", wsUpload)
	http.HandleFunc("/r/wse/midata", wsModData)

	// simg queda fuera de la ruta segura /r
	http.HandleFunc("/wse/simg", wsImg)
}

func wsMicrositio(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		if emp != nil {
			img := model.GetLogo(c, r.FormValue("IdEmp"))
			if(img == nil) {
				img = new(model.Image)
				img.IdEmp = emp.IdEmp
			}
			fd := imgToForm(*img)
			tc := make(map[string]interface{})
			tc["Sess"] =  s
			tc["Empresa"] = emp
			tc["FormData"] = fd
            out.Status = "ok"
		}
	}
}

// upload is the HTTP handler for uploading images; it handles "/".
func wsUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var imgprops appimage.ServingURLOptions
	imgprops.Secure = true
	imgprops.Size = 180
	imgprops.Crop = false
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		imgo := model.GetLogo(c, r.FormValue("IdEmp"))
		if imgo == nil {
			imgo = new(model.Image)
			imgo.IdEmp = emp.IdEmp
		}
		fd := imgToForm(*imgo)
		tc := make(map[string]interface{})
		tc["Sess"] =  s
		tc["Empresa"] = emp
		tc["FormData"] = fd
		if r.Method != "POST" {
			// No upload; show the upload form.
			micrositio(w, r)
			return
		}

		idemp := r.FormValue("IdEmp")
		sp1 := r.FormValue("Sp1")
		sp2 := r.FormValue("Sp2")
		f, _, err := r.FormFile("image")
		model.Check(err)
		defer f.Close()

		// Grab the image data
		var buf bytes.Buffer
		io.Copy(&buf, f)
		i, _, err := image.Decode(&buf)
		if err != nil {
			if(r.FormValue("tipo")=="async") {
				fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo, formato no aceptado");
			} else {
                out.Status = "badformat"
			}
			return
		}

		const max = 600
		// We aim for less than max pixels in any dimension.
		if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
			// If it's gigantic, it's more efficient to downsample first
			// and then resize; resizing will smooth out the roughness.
			if b.Dx() > 2*max || b.Dy() > 2*max {
				w, h := max*2, max*2
				if b.Dx() > b.Dy() {
					h = b.Dy() * h / b.Dx()
				} else {
					w = b.Dx() * w / b.Dy()
				}
				i = resize.Resample(i, i.Bounds(), w, h)
				b = i.Bounds()
			}
			w, h := max, max
			if b.Dx() > b.Dy() {
				h = b.Dy() * h / b.Dx()
			} else {
				w = b.Dx() * w / b.Dy()
			}
			i = resize.Resize(i, i.Bounds(), w, h)
		}

		// Encode as a new JPEG image.
		buf.Reset()
		err = jpeg.Encode(&buf, i, nil)
		if err != nil {
			if(r.FormValue("tipo")=="async") {
				fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo, formato no aceptado");
			} else {
                out.Status = "badEncode"
                return
			}
			return
		}
		var blobkey appengine.BlobKey
		blob, err := blobstore.Create(c, "image/jpeg")
		if err != nil {
			c.Errorf("blobstore Create: %v", idemp)
		}
		_, err = blob.Write(buf.Bytes())
		if err != nil {
			c.Errorf("blobstore Write: %v", idemp)
		}
		err = blob.Close()
		if err != nil {
			c.Errorf("blobstore Close: %v", idemp)
		}
		blobkey, err = blob.Key()
		if err != nil {
			c.Errorf("blobstore Key Gen: %v", idemp)
		}
		if url, err := appimage.ServingURL(c, blobkey, &imgprops); err != nil {
			c.Errorf("Cannot construct EmpLogo ServingURL : %v", idemp)
		} else {
			// Save the image under a unique key, a hash of the image.
			img := &model.Image{
				Data: buf.Bytes(), IdEmp: idemp, IdImg: model.RandId(20), 
				Kind: "EmpLogo", Name: imgo.Name, Desc: imgo.Desc, 
				Sizepx: 0, Sizepy: 0, Url: imgo.Url, Type: "",
				Sp1: sp1, Sp2: sp2, Sp3: string(blobkey), Sp4: url.String(),
				Np1: 0, Np2: 0, Np3: 0, Np4: 0,
			}

			_, err = model.PutLogo(c, img)
			if err != nil {
				if(r.FormValue("tipo")=="async") {
					fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo. Sistema en matenimiento, intente en unos minutos");
				} else {
                    out.Status = "cantsave"
				}
				return
			}
		}

		/* 
			se crea icono
		*/
		val := slogores(c, idemp, 70, 0)
		if val != 0 {
            out.Satus = "cantSave"
			return
		}

		if(r.FormValue("tipo")=="async") {
			fmt.Fprintf(w, "<p></p>");
		} else {
			micrositio(w, r)
		}
		return
	} else {
        out.Status = "ok"
	}
}

func modData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		imgo := model.GetLogo(c, r.FormValue("IdEmp"))
		if(imgo == nil) {
			imgo = new(model.Image)
			imgo.IdEmp = emp.IdEmp
		}
		fd := imgToForm(*imgo)
		tc := make(map[string]interface{})
		tc["Sess"] =  s
		tc["Empresa"] = emp
		tc["FormData"] = fd
		if r.Method != "POST" {
			// No upload; show the upload form.
			micrositio(w, r)
			return
		}

		idemp := r.FormValue("IdEmp")
		name := r.FormValue("Name")
		desc := r.FormValue("Desc")
		url := r.FormValue("Url")
		sp1 := r.FormValue("Sp1")
		sp2 := r.FormValue("Sp2")

	//	key := datastore.NewKey(c, "EmpLogo", r.FormValue("id"), 0, nil)
	//	im := new(model.Image)
		// Save the image under a unique key, a hash of the image.
		imgo = &model.Image{
			Data: imgo.Data, IdEmp: idemp, IdImg: imgo.IdImg,
			Kind: "EmpLogo", Name: name, Desc: desc,
			Sizepx: 0, Sizepy: 0, Url: url, Type: "",
			Sp1: sp1, Sp2: sp2, Sp3: imgo.Sp3, Sp4: imgo.Sp4,
			Np1: 0, Np2: 0, Np3: 0, Np4: 0,
		}
		_, err := model.PutLogo(c, imgo)
		if err != nil {
            out.Status = "cantsave"
			return
		}

		micrositio(w, r)
	} else {
        out.Status = "ok"
	}
}

// keyOf returns (part of) the SHA-1 hash of the data, as a hex string.
func keyOf(data []byte) string {
	sha := sha1.New()
	sha.Write(data)
	return fmt.Sprintf("%x", string(sha.Sum(nil))[0:8])
}

// img is the HTTP handler for displaying images;
// it handles "/simg".
func wsImg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/jpeg")
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "EmpLogo", r.FormValue("id"), 0, nil)
	im := new(model.Image)
	if err := datastore.Get(c, key, im); err != nil {
		if item, err := memcache.Get(c, "defaultimg"); err == memcache.ErrCacheMiss {
			/* OJO: Esta es una cochinada, se creo una empresa con el logo por default y es lo que se escupe
			cuando no hay logo para otra empresa */
			dkey := datastore.NewKey(c, "EmpLogo", "oygqgtyayzxqbl", 0, nil)
			if err := datastore.Get(c, dkey, im); err != nil {
				w.WriteHeader(http.StatusNotFound)
				c.Errorf("ImgCtrl No existe id: %v", "EmpLogoDefault")
				return
			}
			item := &memcache.Item{
				Key:   "defaultimg",
				Value: im.Data,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("Memcache.Add defaultimg : %v", err)
			}

			w.Header().Set("Content-type", "image/jpeg")
			w.Write(im.Data)
		} else {
			//c.Infof("memcache retrieve defaultimg : %v", strconv.Itoa(hit))
			w.Header().Set("Content-type", "image/jpeg")
			w.Write(item.Value)
		}
	} else {
		w.Header().Set("Content-type", "image/jpeg")
		w.Write(im.Data)
	}
}


func delimg(w http.ResponseWriter, r *http.Request) {
}

func imgForm(w http.ResponseWriter, r *http.Request, s sess.Sess, valida bool, tpl *template.Template) (FormDataImage, bool){
	fd := FormDataImage {
		IdEmp: strings.TrimSpace(r.FormValue("IdEmp")),
		Name: strings.TrimSpace(r.FormValue("Name")),
		ErrName: "",
		Url: strings.TrimSpace(r.FormValue("Url")),
		ErrUrl: "",
		Sp1: strings.TrimSpace(r.FormValue("Sp1")),
		ErrSp1: "",
		Sp2: strings.TrimSpace(r.FormValue("Sp2")),
		ErrSp2: "",
		Desc: strings.TrimSpace(r.FormValue("Desc")),
		ErrDesc: "",
	}
	if valida {
		var ef bool
		ef = false
		if fd.Name != "" && !model.ValidName.MatchString(fd.Name) {
			fd.ErrName = "invalid"
			ef = true
		}
		if fd.Url != "" && !model.ValidUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}
		if fd.Sp1 != "" && !model.ValidUrl.MatchString(fd.Sp1) {
			fd.ErrSp1 = "invalid"
			ef = true
		}
		if fd.Sp2 != "" && !model.ValidUrl.MatchString(fd.Sp2) {
			fd.ErrSp2 = "invalid"
			ef = true
		}
		if fd.Desc != "" && !model.ValidSimpleText.MatchString(fd.Desc) {
			fd.ErrDesc = "invalid"
			ef = true
		}
		if ef {
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormData"] = fd
			tpl.Execute(w, tc)
			return fd, false
		}
	}
	return fd, true
}

func imgFill(r *http.Request, img *model.Image) {
	img.Name=		strings.TrimSpace(r.FormValue("Name"))
	img.Desc=		strings.TrimSpace(r.FormValue("Desc"))
	img.Url=		strings.TrimSpace(r.FormValue("Url"))
	img.Sp1=		strings.TrimSpace(r.FormValue("Sp1"))
	img.Sp2=		strings.TrimSpace(r.FormValue("Sp2"))
}

func imgToForm(e model.Image) FormDataImage {
	fd := FormDataImage {
		IdEmp:		e.IdEmp,
		IdImg:		e.IdImg,
		Kind:		e.Kind,
		Name:		e.Name,
		Desc:		e.Desc,
		Sizepx:		e.Sizepx,
		Sizepy:		e.Sizepy,
		Url:		e.Url,
		Type:		e.Type,
		Sp1:		e.Sp1,
		Sp2:		e.Sp2,
		Sp3:		e.Sp3,
		Sp4:		e.Sp4,
		Np1:		e.Np1,
		Np2:		e.Np2,
		Np3:		e.Np3,
		Np4:		e.Np4,
	}
	return fd
}
