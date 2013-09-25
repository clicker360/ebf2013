// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be fouerrnd in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"appengine/blobstore"
    "appengine/memcache"
	"appengine/image"
	"encoding/json"
	"net/http"
	"sess"
	"model"
	"time"
	//"io"
)

type OfImg struct{
	IdOft		string `json:"idoft"`
	IdBlob		string `json:"blobkey"`
	Status		string `json:"status"`
	UploadURL	string `json:"UploadUrl"`
}

type detalle struct {
	IdEmp		string		`json:"idemp"`
	IdOft		string		`json:"idoft"`
	IdCat		int			`json:"idcat"`
	Oferta		string		`json:"oferta"`
	Empresa		string		`json:"empresa"`
	Descripcion	string		`json:"descripcion"`
	Enlinea     bool		`json:"enlinea"`
	Url         string		`json:"url"`
	BlobKey		appengine.BlobKey	`json:"imgurl"`
}


// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/r/ofimgup", handleUpload)
	// ofimg queda fuera del url seguro /r
	http.HandleFunc("/ofimgb", handleServe)
	http.HandleFunc("/ofimgi", handleServeImg)
	http.HandleFunc("/ofimg", handleServeImgByIdOrBlob)

	//http.HandleFunc("/r/ofimgform", handleRoot)
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

func handleServeImgByIdOrBlob(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "none" {
		c := appengine.NewContext(r)
		var imgprops image.ServingURLOptions
		imgprops.Secure = true
		imgprops.Size = 400
		imgprops.Crop = false
		if model.ValidID.MatchString(r.FormValue("id")) {
			/* Cuando es un Id normal */
			//c.Infof("Blob : %v", r.FormValue("id"))
			//now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
			var timetolive = 900 //seconds
			var b []byte
			var d detalle
			if item, err := memcache.Get(c, "d_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
				oferta := model.GetOferta(c, r.FormValue("id"))
				if(oferta.BlobKey != "none") {
					//if now.After(oferta.FechaHoraPub) {
						d.IdEmp = oferta.IdEmp
						d.IdOft = oferta.IdOft
						d.IdCat = oferta.IdCat
						d.Oferta = oferta.Oferta
						d.Empresa = oferta.Empresa
						d.Descripcion = oferta.Descripcion
						d.Enlinea = oferta.Enlinea
						d.Url = oferta.Url
						d.BlobKey = oferta.BlobKey
					//}

					b, _ = json.Marshal(d)
					item := &memcache.Item{
						Key:   "d_"+r.FormValue("id"),
						Value: b,
						Expiration: time.Duration(timetolive)*time.Second,
					}
					if err := memcache.Add(c, item); err == memcache.ErrNotStored {
						c.Errorf("Memcache.Add d_idoft : %v", err)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			} else {
				c.Infof("memcache retrieve d_idoft : %v", r.FormValue("id"))
				if err := json.Unmarshal(item.Value, &d); err != nil {
					c.Errorf("Unmarshaling ShortLogo item: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			blobstore.Send(w, d.BlobKey)
		} else {
			/* Cuando es un BlobKey */
			c.Infof("Blob : %v", r.FormValue("id"))
			blobstore.Send(w, appengine.BlobKey(r.FormValue("id")))
		}
	}
}

func handleServeImgById(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "none" {
		c := appengine.NewContext(r)
		oft := model.GetOferta(c, r.FormValue("id"))
		if(oft.BlobKey != "none") {
			var imgprops image.ServingURLOptions
			imgprops.Secure = true
			imgprops.Size = 400
			imgprops.Crop = false

			if url, err := image.ServingURL(c, oft.BlobKey, &imgprops); err != nil {
				c.Infof("Cannot construct ServingURL : %v", r.FormValue("id"))
				blobstore.Send(w, oft.BlobKey)
			} else {
				http.Redirect(w, r, url.String(), http.StatusFound)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	return
}

/* 
 * dejamos esto como referencia
 * El envío de la liga de sesión de upload se genera en ofadm.go
 *
func handleRoot(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
		uploadURL, err := blobstore.UploadURL(c, "/ofimgup", nil)
		if err != nil {
		serveError(c, w, err)
		return
	}
	tc := make(map[string]interface{})
	tc["UploadURL"] = uploadURL
	tc["IdOft"] =  r.FormValue("IdOft")
	w.Header().Set("Content-Type", "text/html")
	rootTemplate.ExecuteTemplate(w, "ofupform", tc)
	return
}
*/

func handleUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out OfImg
	blobOpts := blobstore.UploadURLOptions {
		MaxUploadBytesPerBlob: 1048576,
	}
	out.Status = "invalidId"
	out.IdBlob = ""
	if s, ok := sess.IsSess(w, r, c); ok {
        u, _ := model.GetCta(c, s.User)
		blobs, form, err := blobstore.ParseUpload(r)
		file := blobs["image"]
        out.IdOft = form.Get("IdOft")
        out.IdBlob = string(file[0].BlobKey)
        if file[0].ContentType != "image/png" && file[0].ContentType != "image/jpeg" {
			out.Status = "invalidContentType"
        } else {
            if err != nil {
                out.Status = "invalidParseUpload"
                berr := blobstore.Delete(c, file[0].BlobKey)
                model.Check(berr)
            } else {
                oferta := model.GetOferta(c, out.IdOft)
                empresa, _ := u.GetEmpresa(c, oferta.IdEmp)
                if oferta.IdEmp == "none" {
                    out.Status = "invalidIdOft"
                    berr := blobstore.Delete(c, file[0].BlobKey)
                    model.Check(berr)
                } else {
                    out.Status = "ok"
                    if len(file) == 0 {
                        out.Status = "invalidFileSize0"
                        berr := blobstore.Delete(c, file[0].BlobKey)
                        model.Check(berr)
                    } else {
                        var oldblobkey = oferta.BlobKey
                        oferta.BlobKey = file[0].BlobKey
                        out.IdOft = oferta.IdOft

                        // Se crea la URL para servir la oferta desde el CDN, si no se puede
                        var imgprops image.ServingURLOptions
                        imgprops.Secure = true
                        imgprops.Size = 400
                        imgprops.Crop = false
                        if url, err := image.ServingURL(c, oferta.BlobKey, &imgprops); err != nil {
                            c.Errorf("Cannot construct ServingURL : %v", oferta.IdOft)
                            oferta.ImageBig = ""
                        } else {
                            oferta.ImageBig = url.String()
                        }

                        _,err = empresa.PutOferta(c, oferta)
                        if err != nil {
                            out.Status = "invalidUploadWriteErr"
                            berr := blobstore.Delete(c, file[0].BlobKey)
                            model.Check(berr)
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
                    }
                }
            }
        }
	}
    UploadURL, err := blobstore.UploadURL(c, "/r/ofimgup", &blobOpts)
    out.UploadURL = UploadURL.String()
    if err != nil {
        out.Status = "uploadSessionError"
    }
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}
