// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"appengine/blobstore"
	"net/http"
	"strings"
	"strconv"
	"model"
	"sess"
	"time"
)

type WsEmpresa struct {
	IdEmp			string `json:"idemp"`
	Nombre			string `json:"nombre"`
	Url				string `json:"url,omitempty"`
	Status		    string `json:"status"`
	Ackn		    string `json:"ackn,omitempty"`
	Errors		    map[string]bool `json:"errors,omitempty"`
	Ofertas	    *[]model.Oferta `json:"ofertas,omitempty"`
}


type WsOferta struct {
	IdOft       string `json:"IdOft,omitempty"`
	IdEmp       string `json:"IdEmp,omitempty"`
	IdCat       int `json:"IdCat,omitempty"`
	Empresa		string `json:"Empresa,omitempty"`
	Oferta		string `json:"Oferta,omitempty"`
	Descripcion	string `json:"Descripcion,omitempty"`
	Codigo      string `json:"Codigo,omitempty"`
	Precio      string `json:"Precio,omitempty"`
	Descuento   string `json:"Descuento,omitempty"`
	Promocion	string `json:"Promocion,omitempty"`
	Enlinea     bool `json:"Enlinea,omitempty"`
	Url         string `json:"Url,omitempty"`
	Meses       string `json:"Meses,omitempty"`
	FechaHoraPub    time.Time `json:"FechaPub,omitempty"`
	StatusPub   bool `json:"StatusPub,omitempty"`
	FechaHora   time.Time `json:"FechaHora,omitempty"`
    BlobKey appengine.BlobKey `json:"BlobKey,omitempty"`
    ImageSmall  string `json:"ImageSmall,omitempty"`
    ImageBig    string `json:"imageBig,omitempty"`

	Ofertas	    *[]model.Oferta `json:"ofertas,omitempty"`
	Categorias	*[]model.Categoria `json:"categorias,omitempty"`
    Sucursales  *[]model.Sucursal `json:"sucursales,omitempty"`
	UploadURL	string `json:"UploadUrl,omitempty"`
	Ackn		string `json:"ackn,omitempty"`
	Status	    string `json:"status,omitempty"`
	Errors		map[string]bool `json:"errors,omitempty"`
}

var Loc *time.Location

func init() {
    Loc, _ = time.LoadLocation("America/Mexico_City")
	http.HandleFunc("/r/wso/put", PutOferta)
	http.HandleFunc("/r/wso/post", PostOferta)
	http.HandleFunc("/r/wso/get", GetOferta)
	http.HandleFunc("/r/wso/gets", GetOfertas)
	http.HandleFunc("/r/wso/del", DelOferta)
}


/*
	Regresa todas las ofertas por empresa
*/
func GetOfertas(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsEmpresa
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    } else {
		// se obtiene el detalle de cta
		u, _ := model.GetCta(c, s.User)

		if r.Method != "GET" {
			out.Status = "wrongMethod"
			return
		}
		if empresa := u.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
			out.IdEmp = empresa.IdEmp
			out.Nombre = empresa.Nombre
			out.Url = empresa.Url
			out.Status = "ok"
			out.Ofertas = empresa.ListOf(c)
		} else {
			out.Status = "notFound"
		}
	}
    return
}

/*
	Regresa un detalle de oferta por id
    Incluye arreglo de categorías con la selecionada
    Incluye enlace de subida para blob de imagen
*/
func GetOferta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsOferta
    defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    if r.Method != "GET" {
		out.Status = "wrongMethod"
        return
    }
	if oferta := model.GetOferta(c, r.FormValue("IdOft")); oferta == nil {
		out.Status = "notFound"
        return
    } else {
        setWsOferta(&out, *oferta)
        out.Categorias = model.ListCat(c, oferta.IdCat);

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
    return
}

/*
    Crea una oferta con una empresa, requiere IdEmp
    La primera vez que se crea una oferta no se integran Sucursales, 
    Palabras clave, Estado x oferta, ni imagen.
    
    Una vez inicializada la oferta, el método de modificación se encarga de eso
    Para lo cual el formato debe entonces ya mostrar todos los campos con la referencia 
    del IdOft creado.
*/
func PutOferta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsOferta
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
    } else {

        // se obtiene el detalle de cta
        u, _ := model.GetCta(c, s.User)

        // PUT
        if r.Method != "POST" {
            out.Status = "wrongMethod"
            return
        }

        // Se obtienen y validan los campos del cgi
        out.IdEmp = r.FormValue("IdEmp")
        oTmp := fill(r)
        oTmp.IdOft = "" // Se inicializa el id pase lo que pase, es un put!
        oTmp.BlobKey = "none"
        if out.Errors, out.Status = validate(oTmp); out.Status != "ok" {
            out.Status = "invalidInput"
            return
        }
        if empresa, err := u.GetEmpresa(c, out.IdEmp); err != nil {
            out.Status = "notFound"
        } else {
            oTmp.IdEmp = empresa.IdEmp
            oTmp.Empresa = strings.ToUpper(empresa.Nombre)
            if _, err := empresa.PutOferta(c, &oTmp); err != nil {
                out.Status = "writeErr"
                return
            }
            // Se pasa a la estructura de salida para JSON
            setWsOferta(&out, oTmp)
            out.Status = "ok"
        }
    }
    return
}

/*
    Modifica una oferta, requiere IdOft
*/
func PostOferta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsOferta
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
    } else {
        // se obtiene el detalle de cta
        u, _ := model.GetCta(c, s.User)

        // POST
        if r.Method != "POST" {
            out.Status = "wrongMethod"
            return
        }

        // Se obtienen y validan los campos del cgi
        oTmp := fill(r)
        if out.Errors, out.Status = validate(oTmp); out.Status != "ok" {
            out.Status = "invalidInput"
            return
        }
        if empresa, err := u.GetEmpresa(c, oTmp.IdEmp); err != nil {
            out.Status = "notFound"
        } else {
            oTmp.IdEmp = empresa.IdEmp
            oTmp.Empresa = strings.ToUpper(empresa.Nombre)
            //oTmp.ImageSmall = emplogo.Sp4
            //oTmp.ImageBig = strings.Replace(emplogo.Sp4, "s180", "s70",1)

            // Se modifica la oferta
            // Se agrega un lock a la oferta en cache
            lock, locked := model.LockItem(r, "Oferta", oTmp.IdOft)
            if locked {
                if keyOferta, err := empresa.PutOferta(c, &oTmp); err != nil {
                    out.Ackn = err.Error()
                    out.Status = "writeError" 
                } else {
                    // Se borran las relaciones oferta-sucursal
                    if err := model.DelOfertaSucursales(c, oTmp.IdOft); err != nil {
                        out.Ackn = err.Error()
                        out.Status = "relationError" 
                    }
                    // Se crea un mapa de Estados para agregar a OfertaEstado
                    edomap := make(map[string]string,32)

                    // Se reconstruyen las Relaciones oferta-sucursal con las solicitadas
                    idsucs := strings.Fields(r.FormValue("schain"))

                    for _, idsuc := range idsucs {
                        suc := model.GetSuc(c, idsuc)

                        lat, _ := strconv.ParseFloat(suc.Geo1, 64)
                        lng, _ := strconv.ParseFloat(suc.Geo2, 64)

                        var ofsuc model.OfertaSucursal
                        ofsuc.IdOft = oTmp.IdOft
                        ofsuc.IdSuc = idsuc
                        ofsuc.IdEmp = oTmp.IdEmp
                        ofsuc.Sucursal = suc.Nombre
                        ofsuc.Lat = lat
                        ofsuc.Lng = lng
                        ofsuc.Empresa = oTmp.Empresa
                        ofsuc.Oferta = oTmp.Oferta
                        ofsuc.Descripcion = oTmp.Descripcion
                        ofsuc.Promocion = oTmp.Promocion
                        ofsuc.Descuento = oTmp.Descuento
                        ofsuc.Url = oTmp.Url
                        ofsuc.StatusPub = oTmp.StatusPub
                        ofsuc.FechaHora = time.Now().In(Loc)

                        // Se añade el estado de la sucursal al mapa de estados
                        edomap[suc.DirEnt] = oTmp.IdOft
                        if err := empresa.PutOfertaSucursal(c, &ofsuc); err != nil {
                            out.Ackn = err.Error()
                            out.Status = "relationError"
                        }
                    }

                    // Se limpia la relación OfertaEstado
                    if err := oTmp.DelOfertaEstado(c); err != nil {
                        out.Ackn = err.Error()
                        out.Status = "relationError"
                    }

                    // Se guarda la relación OfertaEstado
                    if err := empresa.PutOfertaEstado(c, oTmp.IdOft, edomap); err != nil {
                        out.Ackn = err.Error()
                        out.Status = "relationError"
                    }

                    // Se despacha la generación de diccionario de palabras
                    putSearchData(c, oTmp.Empresa+" "+oTmp.Oferta+" "+oTmp.Descripcion+" "+r.FormValue("pchain"), keyOferta, oTmp.IdOft, oTmp.IdCat, oTmp.Enlinea)

                    // Si todo salio bien se desbloquea el item
                    if unlocked := model.UnlockItem(r, lock); !unlocked {
                        out.Ackn = "itemNotUnlocked"
                        c.Infof("WARNING!!!! Unlocked item returned false, verify that memecache key does not exist: lock_%v_%v", lock.Kind, lock.Id)
                    }

                    out.Categorias = model.ListCat(c, oTmp.IdCat);

                    /*
                     * Se crea el url para el form action encargado del upload del blob de imagen
                     */
                    if url, err := setUploadUrl(r); err != nil {
                        out.Ackn = "uploadUrlError"
                    } else {
                        out.UploadURL = url
                    }
                    out.Status = "ok"
                    setWsOferta(&out, oTmp)
                }
            } else {
                out.Status = "itemLocked"
            }
        }
    }
    return
}

/*
 * Requiere IdEmp y IdOft
 */
func DelOferta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsSucursal
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    } else {
		// se obtiene el detalle de cta
		u, _ := model.GetCta(c, s.User)

		// DELETE
		if r.Method != "GET" {
			out.Status = "wrongMethod"
			return
		}
		if empresa := u.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
			if err := empresa.DelOferta(c, r.FormValue("IdOft")); err != nil {
				out.Status = "notFound"
			} else {
				out.Status = "ok"
			}
		} else {
			out.Status = "notFound"
		}
	}
    return
}

func setUploadUrl(r *http.Request) (string, error) {
	c := appengine.NewContext(r)
    blobOpts := blobstore.UploadURLOptions{
        MaxUploadBytesPerBlob: 1048576,
    }
    if uploadURL, err := blobstore.UploadURL(c, "/r/ofimgup", &blobOpts); err != nil {
        return "", err
    } else {
        c.Infof("TLS: %v", r.TLS)
        if(r.TLS != nil) {
            return strings.Replace(uploadURL.String(), "http", "https", 1), nil
        } else {
            /* return uploadURL.String(), nil */
            return strings.Replace(uploadURL.String(), "http", "https", 1), nil
        }
    }
}

func setWsOferta(out *WsOferta, oferta model.Oferta) {
        out.IdOft = oferta.IdOft
        out.IdEmp = oferta.IdEmp
        out.IdCat = oferta.IdCat
        out.Empresa = oferta.Empresa
        out.Oferta = oferta.Oferta
	    out.Descripcion = oferta.Descripcion
	    out.Codigo = oferta.Codigo
	    out.Precio = oferta.Precio
	    out.Descuento = oferta.Descuento
	    out.Promocion = oferta.Promocion
	    out.Enlinea = oferta.Enlinea
	    out.Url = oferta.Url
        out.BlobKey = oferta.BlobKey
        out.ImageSmall = oferta.ImageSmall
        out.ImageBig = oferta.ImageBig
	    out.Meses = oferta.Meses
	    out.FechaHoraPub = oferta.FechaHoraPub
	    out.StatusPub = oferta.StatusPub
	    out.FechaHora = time.Now().In(Loc)
}

func fill(r *http.Request) model.Oferta {
    var fh time.Time
	if r.FormValue("FechaHoraPub") != "" {
		fh, _ = time.ParseInLocation("_2 Jan 15:04:05", strings.TrimSpace(r.FormValue("FechaHoraPub"))+" 00:00:00",Loc)
		fh = fh.AddDate(2013,0,0).In(Loc)
	} else {
		fh = time.Now().In(Loc)
	}
	el, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("Enlinea")))
	st, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("StatusPub")))
	ic, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("IdCat")))
	o := model.Oferta {
		IdEmp:		    strings.TrimSpace(r.FormValue("IdEmp")),
		IdOft:		    strings.TrimSpace(r.FormValue("IdOft")),
		IdCat:		    ic,
		Empresa:		strings.ToUpper(strings.TrimSpace(r.FormValue("Empresa"))),
		Oferta:		    strings.TrimSpace(r.FormValue("Oferta")),
		Descripcion:	strings.TrimSpace(r.FormValue("Descripcion")),
		Codigo:		    strings.TrimSpace(r.FormValue("Codigo")),
		Precio:		    strings.TrimSpace(r.FormValue("Precio")),
		Descuento:		strings.TrimSpace(r.FormValue("Descuento")),
		Promocion:		strings.TrimSpace(r.FormValue("Promocion")),
		Enlinea:		el,
		Url:			strings.TrimSpace(r.FormValue("Url")),
		Meses:		    strings.TrimSpace(r.FormValue("Meses")),
		FechaHoraPub:	fh,
		StatusPub:		st,
		FechaHora:	    time.Now().In(Loc),
		BlobKey:		appengine.BlobKey(strings.TrimSpace(r.FormValue("BlobKey"))),
    }
    return o
}

func validate(s model.Oferta) (map[string]bool, string) {
    errmsg := "ok"
    err := make(map[string]bool)
	if s.Oferta == "" || !model.ValidSimpleText.MatchString(s.Oferta) {
        err["Oferta"] = false
    }
	if s.Descripcion == "" || !model.ValidSimpleText.MatchString(s.Descripcion) && len(s.Descripcion) > 200 {
        err["Descripcion"] = false
    }
	if s.Url != "" && !model.ValidUrl.MatchString(s.Url) {
        err["Url"] = false
    }

    for _, v := range err {
        if v == false {
            errmsg = "invalidInput"
        }
    }
	return err, errmsg
}
