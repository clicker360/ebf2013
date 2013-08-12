// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	//"appengine/datastore"
	"appengine/blobstore"
	//"appengine/urlfetch"
	"html/template"
	"net/http"
	"strings"
	"strconv"
	//"net/url"
	"model"
	"sess"
	"time"
	//"fmt"
	"io"
)

type FormDataOf struct {
	IdOft			string
	IdEmp			string
	IdCat			int
	Categorias		*[]model.Categoria
	Empresa			string
	Oferta			string
	ErrOferta		string
	Descripcion		string
	ErrDescripcion	string
	Enlinea			bool
	Url				string
	ErrUrl			string
	Meses			string
	ErrMeses		string
	FechaHoraPub    time.Time
	ErrFechaHoraPub string
	StatusPub		bool
	FechaHora		time.Time
	Ackn			string
	Sucursales		string // cadena de id's de sucursales separadas por espacio
	//UploadURL		*url.URL
	UploadURL		string
	BlobKey			appengine.BlobKey
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.

func init() {
	http.HandleFunc("/r/of", model.ErrorHandler(OfShow))
	http.HandleFunc("/r/ofs", model.ErrorHandler(OfShowList))
	http.HandleFunc("/r/ofmod", model.ErrorHandler(OfMod))
	http.HandleFunc("/r/ofdel", model.ErrorHandler(OfDel))
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}

func OfShowList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
			tc["Empresa"] = empresa
			tc["Oferta"] = model.ListOf(c, empresa.IdEmp)
		}
		ofadmTpl.ExecuteTemplate(w, "ofertas", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

/*
 * se pasa eso a model/ofertas.go
func listOf(c appengine.Context, IdEmp string) *[]model.Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdEmp =", IdEmp).Limit(500)
	n, _ := q.Count(c)
	ofertas := make([]model.Oferta, 0, n)
	if _, err := q.GetAll(c, &ofertas); err != nil {
		return nil
	}
	sortutil.AscByField(ofertas, "Oferta")
	return &ofertas
}

func listCat(c appengine.Context, IdCat int) *[]model.Categoria {
	q := datastore.NewQuery("Categoria")
	n, _ := q.Count(c)
	categorias := make([]model.Categoria, 0, n)
	if _, err := q.GetAll(c, &categorias); err != nil {
		return nil
	}
	for i, _ := range categorias {
		if(IdCat == categorias[i].IdCat) {
			categorias[i].Selected = `selected`
		}
	}
	return &categorias
}
*/

func OfShow(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		oferta, _ := model.GetOferta(c, r.FormValue("IdOft"))
		var id string
		if oferta.IdEmp != "none" {
			id = oferta.IdEmp
		} else {
			id = r.FormValue("IdEmp")
		}
		fd := ofToForm(*oferta)
		if empresa := model.GetEmpresa(c, id); empresa != nil {
			tc["Empresa"] = empresa
			fd.IdEmp = empresa.IdEmp
			oferta.Empresa = empresa.Nombre
		}
		fd.Categorias = model.ListCat(c, oferta.IdCat);

		/*
		 * Se crea el form para el upload del blob
		 */
		blobOpts := blobstore.UploadURLOptions{
			MaxUploadBytesPerBlob: 1048576,
		}
		uploadURL, err := blobstore.UploadURL(c, "/r/ofimgup", &blobOpts)
		if err != nil {
			serveError(c, w, err)
			return
		}
		fd.UploadURL = strings.Replace(uploadURL.String(), "http", "https", 1)
		//fd.UploadURL = uploadURL

		tc["FormDataOf"] = fd
		ofadmTpl.ExecuteTemplate(w, "oferta", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

// Modifica si hay, Crea si no hay
// Requiere IdEmp. IdOft es opcional, si no hay lo crea, si hay modifica
func OfMod(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		var fd FormDataOf
		var valid bool
		var ofertamod model.Oferta

		if  r.FormValue("IdOft") == "new" {
			if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
				tc["Empresa"] = empresa
				fd.IdEmp = empresa.IdEmp
				fd.Empresa = empresa.Nombre
				ofertamod.IdEmp = empresa.IdEmp
				ofertamod.Oferta = "Nueva oferta"
				ofertamod.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second) // 5 horas menos
				ofertamod.FechaHoraPub = time.Now().Add(time.Duration(model.GMTADJ)*time.Second) // 5 horas menos
				ofertamod.Empresa = strings.ToUpper(empresa.Nombre)
				ofertamod.BlobKey = "none"
				o, err := model.NewOferta(c, &ofertamod)
				model.Check(err)
				fd = ofToForm(*o)
			fd.Ackn = "Ok";
			} else {
				// redireccionar
				http.Redirect(w, r, "/r/le?d=o", http.StatusFound)
			}
		} else {
			/* 
			 * Se pide un id oferta que en teoría existe, se consulta y se cambia
			 * Se valida y si no existe se informa un error
			 */
			fd, valid =ofForm(w, r, true)

			ofertamod.IdOft = fd.IdOft
			ofertamod.IdEmp = fd.IdEmp
			ofertamod.IdCat = fd.IdCat
			ofertamod.Oferta = fd.Oferta
			ofertamod.Descripcion = fd.Descripcion
			ofertamod.Enlinea =	fd.Enlinea
			ofertamod.Url =	fd.Url
			ofertamod.FechaHoraPub = fd.FechaHoraPub
			ofertamod.StatusPub = fd.StatusPub
			//ofertamod.BlobKey = fd.BlobKey
			ofertamod.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)

			oferta, keyOferta := model.GetOferta(c, ofertamod.IdOft)
			if oferta.IdOft != "none" {
				ofertamod.BlobKey = oferta.BlobKey
				ofertamod.Codigo = oferta.Codigo
				if empresa := model.GetEmpresa(c, ofertamod.IdEmp); empresa != nil {
					tc["Empresa"] = empresa
					fd.IdEmp = empresa.IdEmp
					fd.Empresa = empresa.Nombre
					ofertamod.Empresa = strings.ToUpper(empresa.Nombre)
					emplogo := model.GetLogo(c, empresa.IdEmp)
					if emplogo != nil {
						// Tenga lo que tenga, se pasa Sp4 a Oferta.Promocion
						if(emplogo.Sp4 != "")  {
							ofertamod.Promocion = emplogo.Sp4
							ofertamod.Descuento = strings.Replace(emplogo.Sp4, "s180", "s70",1)
						}
					}
				}
				// TODO
				// es preferible poner un regreso avisando que no existe la empresa
				if valid {
					// Ya existe
					err := model.PutOferta(c, &ofertamod)
					model.Check(err)

					// Se borran las relaciones oferta-sucursal
					err = model.DelOfertaSucursales(c, oferta.IdOft)
					model.Check(err)

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
						edomap[suc.DirEnt] = oferta.IdOft

						errOs := ofertamod.PutOfertaSucursal(c, &ofsuc)
						model.Check(errOs)

					}
					// Se limpia la relación OfertaEstado
					_ = ofertamod.DelOfertaEstado(c)

					// Se guarda la relación OfertaEstado
					errOe := ofertamod.PutOfertaEstado(c, edomap)
					model.Check(errOe)

					var tituloOf string
					tituloOf = ""
					if(strings.ToLower(strings.TrimSpace(ofertamod.Oferta)) != "nueva oferta") {
						tituloOf = ofertamod.Oferta
					}
					putSearchData(c, ofertamod.Empresa+" "+tituloOf+" "+ofertamod.Descripcion+" "+r.FormValue("pchain"), keyOferta, oferta.IdOft, ofertamod.IdCat, ofertamod.Enlinea)

					// Se despacha la generación de diccionario de palabras
					// Se agrega pcves a la descripción
					//fmt.Fprintf(w,"http://movil.%s.appspot.com/backend/generatesearch?kind=Oferta&field=Descripcion&id=%s&value=%s&categoria=%s",
					//appengine.AppID(c), keyOferta.Encode(), ofertamod.Descripcion+" "+r.FormValue("pchain"), strconv.Itoa(ofertamod.IdCat))
					//_ = generatesearch(c, keyOferta, ofertamod.Descripcion+" "+r.FormValue("pchain"), ofertamod.IdCat)

					fd = ofToForm(ofertamod)
					fd.Ackn = "Ok";
				}
			} else {
				// no existe la oferta
			}
		}

		fd.Categorias = model.ListCat(c, ofertamod.IdCat);

		/*
		 * Se crea el form para el upload del blob
		 */
		blobOpts := blobstore.UploadURLOptions{
			MaxUploadBytesPerBlob: 1048576,
		}
		uploadURL, err := blobstore.UploadURL(c, "/r/ofimgup", &blobOpts)
		if err != nil {
			serveError(c, w, err)
			return
		}
		fd.UploadURL = strings.Replace(uploadURL.String(), "http", "https", 1)
		//fd.UploadURL = uploadURL
		tc["FormDataOf"] = fd

		ofadmTpl.ExecuteTemplate(w, "oferta", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func OfDel(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); ok {
		if err := model.DelOferta(c, r.FormValue("IdOft")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		OfShowList(w, r)
		return
	}
	http.Redirect(w, r, "/r/ofertas", http.StatusFound)
}

func ofForm(w http.ResponseWriter, r *http.Request, valida bool) (FormDataOf, bool){
	c := appengine.NewContext(r)
	var fh time.Time
	if r.FormValue("FechaHoraPub") != "" {
		fh, _ = time.Parse("_2 Jan 15:04:05", strings.TrimSpace(r.FormValue("FechaHoraPub"))+" 00:00:00")
		fh = fh.AddDate(2012,0,0)
	} else {
		fh = time.Now().Add(time.Duration(model.GMTADJ)*time.Second) // 5 horas menos
	}
	el, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("Enlinea")))
	st, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("StatusPub")))
	ic, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("IdCat")))
	fd := FormDataOf {
		IdOft:			strings.TrimSpace(r.FormValue("IdOft")),
		IdEmp:			strings.TrimSpace(r.FormValue("IdEmp")),
		IdCat:			ic,
		Oferta:			strings.TrimSpace(r.FormValue("Oferta")),
		ErrOferta:		"",
		Descripcion:	strings.TrimSpace(r.FormValue("Descripcion")),
		ErrDescripcion: "",
		Enlinea:		el,
		Url:			strings.TrimSpace(r.FormValue("Url")),
		ErrUrl:			"",
		ErrMeses:		"",
		FechaHoraPub:	fh,
		ErrFechaHoraPub: strings.TrimSpace(fh.Format("_2 Jan")),
		StatusPub:		st,
	}
	if valida {
		var ef bool
		ef = false
		if fd.Oferta == "" || !model.ValidSimpleText.MatchString(fd.Oferta) {
			fd.ErrOferta = "invalid"
			ef = true
		}
		if fd.Descripcion == "" || !model.ValidSimpleText.MatchString(fd.Descripcion) && len(fd.Descripcion) > 200 {
			fd.ErrDescripcion = "invalid"
			ef = true
		}
		if fd.Url != "" && !model.ValidUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}

		fd.Categorias = model.ListCat(c, ic);
		if ef {
			return fd, false
		}
	}
	return fd, true
}

func ofToForm(e model.Oferta) FormDataOf {
	fd := FormDataOf {
		IdOft:			e.IdOft,
		IdEmp:			e.IdEmp,
		IdCat:			e.IdCat,
		Oferta:			e.Oferta,
		Descripcion:	e.Descripcion,
		Enlinea:		e.Enlinea,
		Url:			e.Url,
		FechaHoraPub:	e.FechaHoraPub,
		ErrFechaHoraPub:	strings.TrimSpace(e.FechaHoraPub.Format("_2 Jan")),
		StatusPub:		e.StatusPub,
		BlobKey:		e.BlobKey,
	}
	return fd
}

/*
func generatesearch(c appengine.Context, oftKey *datastore.Key, description string, idcat int) error {
	client := urlfetch.Client(c)
	descurl := fmt.Sprintf(
	"http://movil.%s.appspot.com/backend/generatesearch?kind=Oferta&field=Descripcion&id=%s&value=%s&categoria=%s",
	appengine.AppID(c), oftKey.Encode(), description, strconv.Itoa(idcat))
	_, err := client.Get(descurl)
	if err != nil {
		return err
	}
	return nil
}
*/
var ofadmTpl = template.Must(template.ParseFiles("templates/ofadm.html"))
