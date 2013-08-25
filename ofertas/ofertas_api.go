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

type WsOferta struct {
	IdOft       string `json:"idoft,omitempty"`
	IdEmp       string `json:"idemp,omitempty"`
	IdCat       int `json:"idcat,omitempty"`
	Empresa		string `json:"empresa,omitempty"`
	Oferta		string `json:"oferta,omitempty"`
	Descripcion	string `json:"descripcion,omitempty"`
	Codigo      string `json:"codigo,omitempty"`
	Precio      string `json:"precio,omitempty"`
	Descuento   string `json:"descuento,omitempty"`
	Promocion	string `json:"promocion,omitempty"`
	Enlinea     bool `json:"enlinea,omitempty"`
	Url         string `json:"url,omitempty"`
	Meses       string `json:"meses,omitempty"`
	FechaHoraPub    time.Time `json:"fechapub,omitempty"`
	StatusPub   bool `json:"publicada,omitempty"`
	FechaHora   time.Time `json:"timestamp,omitempty"`
    BlobKey appengine.BlobKey `json:"blobkey,omitempty"`
    ImageSmall  string `json:"imagesmall,omitempty"`
    ImageBig    string `json:"imagebig,omitempty"`

	Ofertas	    *[]model.Oferta `json:"ofertas,omitempty"`
	Categorias	*[]model.Categoria `json:"categorias,omitempty"`
    Sucursales  *[]model.Sucursal `json:"sucursales,omitempty"`
	UploadURL	string `json:"uploadurl,omitempty"`
	Ackn		string `json:"ackn,omitempty"`
	Status	    string `json:"status,omitempty"`
	Errors		map[string]bool `json:"errors,omitempty"`
}

func init() {
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
	if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
        out.IdEmp = empresa.IdEmp
        out.Empresa = empresa.Nombre
        out.Status = "ok"
        out.Ofertas = model.ListOf(c, empresa.IdEmp)
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
	if oferta, key := model.GetOferta(c, r.FormValue("IdOft")); key == nil {
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
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    // PUT
    if r.Method != "POST" {
		out.Status = "wrongMethod"
        return
    }

    // Se obtienen y validan los campos del cgi
    out.IdEmp = r.FormValue("IdEmp")
    oTmp := fill(r)
    oTmp.IdOft = "" // Se inicializa el id pase lo que pase, es un put!
    if out.Errors, out.Status = validate(oTmp); out.Status != "ok" {
        return
    }
    if empresa := model.GetEmpresa(c, out.IdEmp); empresa != nil {
        oTmp.IdEmp = empresa.IdEmp
        oTmp.BlobKey = "none"
        o, err := empresa.PutOferta(c, &oTmp)
        if err != nil {
            out.Status = "writeErr"
            return
        }
        // Se pasa a la estructura de salida para JSON
        setWsOferta(&out, *o)
        out.Status = "ok"
    } else {
        out.Status = "notFound"
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
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    if r.Method != "POST" {
		out.Status = "wrongMethod"
        return
    }

    _, keyOferta := model.GetOferta(c, r.FormValue("IdOft"))
    if keyOferta == nil {
		out.Status = "notFound"
        return
    }
    // Se obtienen y validan los campos del cgi
    oTmp := fill(r)
    if out.Errors, out.Status = validate(oTmp); out.Status != "ok" {
        return
    }

	if empresa := model.GetEmpresa(c, oTmp.IdEmp); empresa != nil {
        oTmp.IdEmp = empresa.IdEmp
        oTmp.Empresa = strings.ToUpper(empresa.Nombre)
        if emplogo := model.GetLogo(c, empresa.IdEmp); emplogo != nil {
            if(emplogo.Sp4 != "")  {
                oTmp.ImageSmall = emplogo.Sp4
                oTmp.ImageBig = strings.Replace(emplogo.Sp4, "s180", "s70",1)
            }
        }

        // Se modifica la oferta
        // Se agrega un lock a la oferta en cache
        lock, locked := model.LockItem(r, "Oferta", oTmp.IdOft)
        if locked {
            if ofertamod, err := empresa.PutOferta(c, &oTmp); err == nil {
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
                    ofsuc.IdOft = ofertamod.IdOft
                    ofsuc.IdSuc = idsuc
                    ofsuc.IdEmp = ofertamod.IdEmp
                    ofsuc.Sucursal = suc.Nombre
                    ofsuc.Lat = lat
                    ofsuc.Lng = lng
                    ofsuc.Empresa = ofertamod.Empresa
                    ofsuc.Oferta = ofertamod.Oferta
                    ofsuc.Descripcion = ofertamod.Descripcion
                    ofsuc.Promocion = ofertamod.Promocion
                    ofsuc.Descuento = ofertamod.Descuento
                    ofsuc.Url = ofertamod.Url
                    ofsuc.StatusPub = ofertamod.StatusPub
                    ofsuc.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)

                    // Se añade el estado de la sucursal al mapa de estados
                    edomap[suc.DirEnt] = ofertamod.IdOft
                    if err := ofertamod.PutOfertaSucursal(c, &ofsuc); err != nil {
                        out.Ackn = err.Error()
                        out.Status = "relationError"
                    }
                }

                // Se limpia la relación OfertaEstado
                if err := ofertamod.DelOfertaEstado(c); err != nil {
                    out.Ackn = err.Error()
                    out.Status = "relationError"
                }

                // Se guarda la relación OfertaEstado
                if err := ofertamod.PutOfertaEstado(c, edomap); err != nil {
                    out.Ackn = err.Error()
                    out.Status = "relationError"
                }

                // Se despacha la generación de diccionario de palabras
                putSearchData(c, ofertamod.Empresa+" "+ofertamod.Oferta+" "+ofertamod.Descripcion+" "+r.FormValue("pchain"), keyOferta, ofertamod.IdOft, ofertamod.IdCat, ofertamod.Enlinea)

                // Si todo salio bien se desbloquea el item
                if unlocked := model.UnlockItem(r, lock); !unlocked {
                    out.Ackn = "itemNotUnlocked"
                    c.Infof("WARNING!!!! Unlocked item returned false, verify that memecache key does not exist: lock_%v_%v", lock.Kind, lock.Id)
                }

                out.Categorias = model.ListCat(c, ofertamod.IdCat);

                /*
                 * Se crea el url para el form action encargado del upload del blob de imagen
                 */
                if url, err := setUploadUrl(r); err != nil {
                    out.Ackn = "uploadUrlError"
                } else {
                    out.UploadURL = url
                }
                out.Status = "ok"
                setWsOferta(&out, *ofertamod)
            } else {
                out.Ackn = err.Error()
                out.Status = "writeError" 
            }
        } else {
            out.Status = "itemLocked"
        }
	} else {
        out.Status = "notFound"
    }
    return
}

func DelOferta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsSucursal
    defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    // DELETE
    if r.Method != "GET" {
		out.Status = "wrongMethod"
        return
    }
	if err := model.DelOferta(c, r.FormValue("IdOft")); err != nil {
		out.Status = "notFound"
    }
	out.Status = "ok"
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
        return strings.Replace(uploadURL.String(), "http", "https", 1), nil
    }
}

func setWsOferta(out *WsOferta, oferta model.Oferta) {
        out.IdEmp = oferta.IdOft
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
	    out.Meses = oferta.Meses
	    out.FechaHoraPub = oferta.FechaHoraPub
	    out.StatusPub = oferta.StatusPub
	    out.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
}

func fill(r *http.Request) model.Oferta {
    loc, _ := time.LoadLocation("America/Mexico_City")
    var fh time.Time
	if r.FormValue("FechaHoraPub") != "" {
		fh, _ = time.ParseInLocation("_2 Jan 15:04:05", strings.TrimSpace(r.FormValue("FechaHoraPub"))+" 00:00:00", loc)
		fh = fh.AddDate(2013,0,0)
	} else {
		fh = time.Now().Add(time.Duration(model.GMTADJ)*time.Second) // 5 horas menos
	}
	el, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("Enlinea")))
	st, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("StatusPub")))
	ic, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("IdCat")))
	o := model.Oferta {
		IdEmp:		    strings.TrimSpace(r.FormValue("IdEmp")),
		IdOft:		    strings.TrimSpace(r.FormValue("IdOft")),
		IdCat:		    ic,
		Empresa:			strings.ToUpper(strings.TrimSpace(r.FormValue("Empresa"))),
		Oferta:		    strings.TrimSpace(r.FormValue("Oferta")),
		Descripcion:		strings.TrimSpace(r.FormValue("Descripcion")),
		Codigo:		    strings.TrimSpace(r.FormValue("Codigo")),
		Precio:		    strings.TrimSpace(r.FormValue("Precio")),
		Descuento:		strings.TrimSpace(r.FormValue("Descuento")),
		Promocion:		strings.TrimSpace(r.FormValue("Promocion")),
		Enlinea:			el,
		Url:				strings.TrimSpace(r.FormValue("Url")),
		Meses:		    strings.TrimSpace(r.FormValue("Meses")),
		FechaHoraPub:	fh,
		StatusPub:		st,
		FechaHora:	    time.Now().Add(time.Duration(model.GMTADJ)*time.Second),
		//BlobKey:		strings.TrimSpace(r.FormValue("IdCat")),
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
