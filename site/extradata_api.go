// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package site

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/image"
	"encoding/json"
	"model"
	"net/http"
	"sess"
	"sortutil"
	"strings"
	"time"
)

type WsExtraData struct {
	IdEmp     string            `json:"IdEmp"`
	Empresa   string            `json:"Empresa"`
	Desc      string            `json:"Desc"`
	Facebook  string            `json:"Facebook,omitempty"`
	Twitter   string            `json:"Twitter,omitempty"`
	BlobKey   appengine.BlobKey `json:"BlobKey,omitempty"`
	ImageUrl  string            `json:"ImageUrl,omitempty"`
	UploadURL string            `json:"UploadUrl,omitempty"`
	FechaHora time.Time         `json:"FechaHora,omitempty"`
	Sp1       string            `json:"Sp1,omitempty"`
	Sp2       string            `json:"Sp2,omitempty"`
	Sp3       string            `json:"Sp3,omitempty"`
	Sp4       string            `json:"Sp4,omitempty"`
	Status    string            `json:"status"`
	Ackn      string            `json:"ackn,omitempty"`
	Errors    map[string]bool   `json:"errors,omitempty"`
}

func init() {
	//Loc, _ = time.LoadLocation("America/Mexico_City")
	http.HandleFunc("/r/wsed/post", PostExtraData)
	http.HandleFunc("/r/wsed/get", GetExtraData)
	http.HandleFunc("/r/wsed/gets", GetExtraDatas)
	http.HandleFunc("/r/extraimgput", handleUpload)
	http.HandleFunc("/extraimg", handleServe)
	http.HandleFunc("/extraimg400", handleServeImg)
}

/*
	Datos extra de empresa. Regresa un dato extra por id de empresa
*/
func GetExtraData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsExtraData
	defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
		return
	} else {
		u, _ := model.GetCta(c, s.User)
		if r.Method != "GET" {
			out.Status = "wrongMethod"
			return
		}
		empresa, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err != nil {
			out.Status = "notFound"
		} else {
			extra, err := u.GetExtraData(c, empresa.IdEmp)
			if err != nil {
				out.Ackn = "noExtraDataYet"
				// se crea el Extra data
				var eTmp model.ExtraData
				eTmp.IdEmp = empresa.IdEmp
				eTmp.Empresa = empresa.Nombre
				eTmp.FechaHora = time.Now() //.In(Loc)
				if err := u.PutExtraData(c, &eTmp); err != nil {
					out.Status = "writeError"
					return
				}
			} else {
				setWsExtraData(&out, *extra)
			}
			out.IdEmp = empresa.IdEmp
			out.Empresa = empresa.Nombre
			/*
			 * Se crea el url para el form action encargado del upload del blob de imagen
			 */
			if url, err := setUploadUrl(r); err != nil {
				out.Ackn = err.Error()
			} else {
				out.UploadURL = url
			}
			out.Status = "ok"
		}
	}
}

/*
	Datos extras de una empresa
*/
func GetExtraDatas(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsExtraData
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
		model.JsonDispatch(w, &out)
		return
	} else {
		u, _ := model.GetCta(c, s.User)
		if r.Method != "GET" {
			out.Status = "wrongMethod"
			model.JsonDispatch(w, &out)
			return
		}
		e := listExtraData(c, u)
		ws := make([]WsExtraData, len(*e), cap(*e))
		for i, v := range *e {
			setWsExtraData(&ws[i], v)
		}
		out.Status = "ok"
		sortutil.AscByField(ws, "Empresa")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		b, _ := json.Marshal(ws)
		w.Write(b)
	}
}

/* OJO NO HAY PUT PARA ESTO PORQUE SE SUPONE QUE DEBE EXISTIR EMPRESA
func PutExtraData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsExtraData
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    } else {
        // se obtiene el detalle de cta
        u, _ := model.GetCta(c, s.User)

        // POST
        if r.Method != "POST" {
            out.Status = "wrongMethod"
            return
        }

        // se construye la estructura temporal de la empresa nueva
        // y se valida
        eTmp := fillExtra(r)

        out.Errors, out.Status = validateExtraData(eTmp)
        if out.Status != "ok" {
            return
        }

        // Se genera un Id nuevo y se agrega a la estructura de Empresa
        eTmp.IdEmp = model.RandId(20)
        if err := u.PutExtraData(c, &eTmp); err != nil {
            out.Status = "writeError"
            return
        }
        setWsExtraData(&out, eTmp)
    }
    if url, err := setUploadUrl(r); err != nil {
        out.Ackn = err.Error()
    } else {
        out.UploadURL = url
    }
    out.Status = "ok"
    return
}
*/

func PostExtraData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsExtraData
	defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
		return
	} else {
		// se obtiene el detalle de cta
		u, _ := model.GetCta(c, s.User)

		// POST
		if r.Method != "POST" {
			out.Status = "wrongMethod"
			return
		}

		// se construye la estructura temporal de la empresa nueva
		// y se valida
		eTmp := fillExtra(r)

		out.Errors, out.Status = validateExtraData(eTmp)
		if out.Status != "ok" {
			return
		}

		// se obtiene la empresa existente a actualizar
		_, err := u.GetExtraData(c, r.FormValue("IdEmp"))
		if err != nil {
			out.Status = err.Error()
			return
		}
		if err := u.PutExtraData(c, &eTmp); err != nil {
			out.Status = "writeError"
			return
		}
		setWsExtraData(&out, eTmp)
	}
	/*
	 * Se crea el url para el form action encargado del upload del blob de imagen
	 */
	if url, err := setUploadUrl(r); err != nil {
		out.Ackn = err.Error()
	} else {
		out.UploadURL = url
	}
	out.Status = "ok"
	return
}

func setUploadUrl(r *http.Request) (string, error) {
	c := appengine.NewContext(r)
	blobOpts := blobstore.UploadURLOptions{
		MaxUploadBytesPerBlob: 1048576,
	}
	if uploadURL, err := blobstore.UploadURL(c, "/r/extraimgput", &blobOpts); err != nil {
		return "", err
	} else {
		c.Infof("TLS: %v", r.TLS)
		if r.TLS != nil {
			return strings.Replace(uploadURL.String(), "http", "https", 1), nil
		} else {
			//return uploadURL.String(), nil
			return strings.Replace(uploadURL.String(), "http", "https", 1), nil
		}
	}
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("id")))
}

func handleServeImg(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "none" {
		c := appengine.NewContext(r)
		var imgprops image.ServingURLOptions
		imgprops.Secure = true
		imgprops.Size = 400
		imgprops.Crop = false
		url, _ := image.ServingURL(c, appengine.BlobKey(r.FormValue("id")), &imgprops)
		http.Redirect(w, r, url.String(), http.StatusFound)
	}
	return
}

func listExtraData(c appengine.Context, u *model.Cta) *[]model.ExtraData {
	q := datastore.NewQuery("Empresa").Ancestor(u.Key(c))
	n, _ := q.Count(c)
	extra := make([]model.ExtraData, 0, n)
	if _, err := q.GetAll(c, &extra); err != nil {
		return nil
	}
	sortutil.AscByField(extra, "Empresa")
	return &extra
}

func setWsExtraData(out *WsExtraData, ed model.ExtraData) {
	out.IdEmp = ed.IdEmp
	out.Empresa = ed.Empresa
	out.Desc = ed.Desc
	out.Facebook = ed.Facebook
	out.Twitter = ed.Twitter
	out.BlobKey = ed.BlobKey
	out.ImageUrl = ed.ImageUrl
	//out.FechaHora = ed.FechaHora
	out.Sp1 = ed.Sp1
	out.Sp2 = ed.Sp2
	out.Sp3 = ed.Sp3
	out.Sp4 = ed.Sp4
}

func fillExtra(r *http.Request) model.ExtraData {
	o := model.ExtraData{
		IdEmp:     strings.TrimSpace(r.FormValue("IdEmp")),
		Empresa:   strings.ToUpper(strings.TrimSpace(r.FormValue("Empresa"))),
		Desc:      strings.TrimSpace(r.FormValue("Desc")),
		Facebook:  strings.TrimSpace(r.FormValue("Facebook")),
		Twitter:   strings.TrimSpace(r.FormValue("Twitter")),
		BlobKey:   appengine.BlobKey(strings.TrimSpace(r.FormValue("BlobKey"))),
		ImageUrl:  strings.TrimSpace(r.FormValue("ImageUrl")),
		FechaHora: time.Now(), //.In(Loc),
		Sp1:       strings.TrimSpace(r.FormValue("Sp1")),
		Sp2:       strings.TrimSpace(r.FormValue("Sp2")),
		Sp3:       strings.TrimSpace(r.FormValue("Sp3")),
		Sp4:       strings.TrimSpace(r.FormValue("Sp4")),
	}
	return o
}

func validateExtraData(s model.ExtraData) (map[string]bool, string) {
	errmsg := "ok"
	err := make(map[string]bool)
	if s.Facebook != "" && !model.ValidSimpleText.MatchString(s.Facebook) {
		err["Facebook"] = false
	}
	if s.Twitter != "" && !model.ValidSimpleText.MatchString(s.Twitter) {
		err["Twitter"] = false
	}
	if s.Desc != "" && !model.ValidSimpleText.MatchString(s.Desc) && len(s.Desc) > 400 {
		err["Desc"] = false
	}
	if s.ImageUrl != "" && !model.ValidUrl.MatchString(s.ImageUrl) {
		err["ImageUrl"] = false
	}
	/* comodines */
	if s.Sp1 != "" && !model.ValidSimpleText.MatchString(s.Sp1) {
		err["Sp1"] = false
	}
	if s.Sp2 != "" && !model.ValidSimpleText.MatchString(s.Sp2) {
		err["Sp2"] = false
	}
	if s.Sp3 != "" && !model.ValidSimpleText.MatchString(s.Sp3) {
		err["Sp3"] = false
	}
	if s.Sp4 != "" && !model.ValidSimpleText.MatchString(s.Sp4) {
		err["Sp4"] = false
	}
	for _, v := range err {
		if v == false {
			errmsg = "invalidInput"
		}
	}
	return err, errmsg
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsExtraData
	defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
		return
	} else {
		u, _ := model.GetCta(c, s.User)
		if r.Method != "POST" {
			out.Status = "wrongMethod"
			return
		}

		blobOpts := blobstore.UploadURLOptions{
			MaxUploadBytesPerBlob: 1048576,
		}

		blobs, form, err := blobstore.ParseUpload(r)
		file := blobs["image"]
		out.BlobKey = file[0].BlobKey
		out.IdEmp = form.Get("IdEmp")
        out.BlobKey = file[0].BlobKey

        UploadURL, err := blobstore.UploadURL(c, "/r/extraimgput", &blobOpts)
        out.UploadURL = UploadURL.String()
        if err != nil {
            out.Status = "uploadSessionError"
            return
        }

        if file[0].ContentType != "image/png" && file[0].ContentType != "image/jpeg" {
			out.Status = "invalidContentType"
        } else {
            if err != nil {
                out.Status = "invalidParseUpload"
                berr := blobstore.Delete(c, file[0].BlobKey)
                model.Check(berr)
            } else {
                if extra, err := u.GetExtraData(c, out.IdEmp); err != nil {
                    out.Status = "notFound"
                    if err := blobstore.Delete(c, file[0].BlobKey); err != nil {
                        out.Ackn = "cantDeleteBlob"
                    }
                    return
                } else {
                    if len(file) == 0 {
                        out.Status = "invalidFileSize0"
                        if err := blobstore.Delete(c, file[0].BlobKey); err != nil {
                            out.Ackn = "cantDeleteBlob"
                        }
                        return
                    }
                    var oldblobkey = extra.BlobKey
                    extra.BlobKey = file[0].BlobKey
                    out.IdEmp = extra.IdEmp

                    // Se crea la URL para servir la extra desde el CDN, si no se puede
                    var imgprops image.ServingURLOptions
                    imgprops.Secure = true
                    imgprops.Size = 400
                    imgprops.Crop = false
                    if url, err := image.ServingURL(c, extra.BlobKey, &imgprops); err != nil {
                        c.Errorf("Cannot construct ServingURL : %v", extra.IdEmp)
                        extra.ImageUrl = ""
                    } else {
                        extra.ImageUrl = url.String()
                    }

                    if err := u.PutExtraData(c, extra); err != nil {
                        out.Status = "invalidUpload"
                        if err := blobstore.Delete(c, file[0].BlobKey); err != nil {
                            out.Ackn = "cantDeleteBlob"
                        }
                        return
                    }
                    /*
                       Se borra el blob anterior, porque siempre crea uno nuevo
                       No se necesita revisar el error
                       Si es el blobkey = none no se borra por obvias razones
                       Se genera una sesion nueva de upload en caso de que quieran
                       cambiar la imágen en la misma pantalla. Esto es debido a que
                       se utiliza un form estático con ajax
                    */
                    if oldblobkey != "none" {
                        blobstore.Delete(c, oldblobkey)
                    }
                    out.Status = "ok"
                }
            }
        }
	}
}
