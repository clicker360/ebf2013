package site

import (
    "appengine"
    "appengine/datastore"
	"html/template"
	"strings"
	"sortutil"
    "net/http"
	"strconv"
    "time"
	"model"
	"sess"
//	"fmt"
)

type FormDataEmp struct {
	IdEmp			string
	RFC				string
	ErrRFC			string
	Nombre			string
	ErrNombre		string
	RazonSoc		string
	ErrRazonSoc		string
	DirCalle		string
	ErrDirCalle		string
	DirCol			string
	ErrDirCol		string
	DirEnt			string
	ErrDirEnt		string
	Entidades		*[]model.Entidad
	Organismos		*[]model.Organismo
	DirMun			string
	ErrDirMun		string
	DirCp			string
	ErrDirCp		string
	NumSuc			string
	ErrNumSuc		string
	OrgEmp			string
	ErrOrgEmp		string
	OrgEmpOtro		string
	ErrOrgEmpOtro	string
	OrgEmpReg		string
	ErrOrgEmpReg	string
	Url				string
	ErrUrl			string
	PartLinea		int
	ExpComer		int
	Desc			string
	ErrDesc			string
}

var pageTpl *template.Template

func init() {
    /* CRUD */
    http.HandleFunc("/r/e/l", RenderList)
    http.HandleFunc("/r/e/c", RenderCreate)
    http.HandleFunc("/r/e/r", Remove)
    http.HandleFunc("/r/e/u", RenderUpdate)
    http.HandleFunc("/r/e/d", RenderDetails)
}

func RenderList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		tc["Empresa"] = listEmp(c, u)
		if(r.FormValue("type") == "e" || r.FormValue("type") == "") {
		    tc["Page"] = model.PageMeta{"lea", "Empresas"}
            pageTpl = model.PrepareTemplate("empresas", "")
		} else if(r.FormValue("type") == "s") {
            tc["Page"] = model.PageMeta{"les", "Sucursales"}
            pageTpl = model.PrepareTemplate("sucursales", "")
		} else if(r.FormValue("type") == "m") {
            tc["Page"] = model.PageMeta{"lem", "Micrositios"}
            pageTpl = model.PrepareTemplate("micrositios", "")
		} else if(r.FormValue("type") == "o") {
            tc["Page"] = model.PageMeta{"leo", "Ofertas"}
            pageTpl = model.PrepareTemplate("ofertas", "")
		}
        if err := pageTpl.ExecuteTemplate(w, "pageRoot", tc); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func RenderDetails(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		e, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderForm(c, w, &s, e)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func RenderUpdate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		// formato con validación 
		_, valid := formatoEmp(w, r, s, true)
		if !valid {return}

		/*
		 * Se carga un struct con los datos de la forma
		 */
		eTmp := fillEmpresa(r)

		/*
		 * Se requiere leer la estructura para insertar el folio y el status
		 * y otros campos que se deben conservar 
		 */
		e, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err != nil {
			c.Errorf("RenderUpdate() GetEmpresa() Error al intentar actualizar Empresa : %v", e.IdEmp)
		} else {
			eTmp.Folio = e.Folio
			eTmp.Status = e.Status
			eTmp.Benef = e.Benef

			_, err := u.PutEmpresa(c, &eTmp)
			if err == datastore.ErrNoSuchEntity {
				c.Errorf("RenderUpdate() PutEmpresa() Error al intentar actualizar Empresa : %v", e.IdEmp)
			}
		}
		renderForm(c, w, &s, &eTmp)
		//RenderList(w, r)
	} else {
		defer http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func RenderCreate(w http.ResponseWriter, r *http.Request) {
	// formato con validación 
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		_, valid := formatoEmp(w, r, s, true)
		if !valid { return }
		u, _ := model.GetCta(c, s.User)
		eTmp := fillEmpresa(r)

		// Se añade una empresa
		//e, err := u.NewEmpresa(c, &eTmp)
		e, err := u.NewEmpresa(c, &eTmp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err1 := model.PutCtaEmp(c, eTmp.IdEmp, u.Email, u.EmailAlt)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		renderForm(c, w, &s, e)
		//RenderList(w, r)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		if err := u.DelEmpresa(c, r.FormValue("IdEmp")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		RenderList(w, r)
		return
	}
	http.Redirect(w, r, "/r/registro", http.StatusFound)
}

func renderForm(c appengine.Context, w http.ResponseWriter, s *sess.Sess, e *model.Empresa) {
	fd := empresaToForm(*e)
	fd.Entidades = model.ListEnt(c, e.DirEnt)
	fd.Organismos = ListOrg(c, e.OrgEmp)
	tc := make(map[string]interface{})
	tc["Sess"] = s
	tc["FormDataEmp"] = fd
    tc["Page"] = model.PageMeta{"lea", "Empresas"}
    pageTpl = model.PrepareTemplate("empresa", "empresa")
    if err := pageTpl.ExecuteTemplate(w, "pageRoot", tc); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	return
}

func listEmp(c appengine.Context, u *model.Cta) *[]model.Empresa {
	q := datastore.NewQuery("Empresa").Ancestor(u.Key(c))
	n, _ := q.Count(c)
	empresas := make([]model.Empresa, 0, n)
	if _, err := q.GetAll(c, &empresas); err != nil {
		return nil
	}
	sortutil.AscByField(empresas, "Nombre")
	return &empresas
}

func ListOrg(c appengine.Context, siglas string) *[]model.Organismo {
	q := datastore.NewQuery("Organismo")
	orgs := make([]model.Organismo, 0, 32)
	if _, err := q.GetAll(c, &orgs); err != nil {
		return nil
	}
	for i, _ := range orgs {
		if(siglas == orgs[i].Siglas) {
			orgs[i].Selected = `selected`
		}
	}
	return &orgs
}

func formatoEmp(w http.ResponseWriter, r *http.Request, s sess.Sess, valida bool) (FormDataEmp, bool){
	c := appengine.NewContext(r)
	partlinea, _ := strconv.Atoi(r.FormValue("PartLinea"))
	expcomer, _ := strconv.Atoi(r.FormValue("ExpComer"))
	fd := FormDataEmp {
		IdEmp: strings.TrimSpace(r.FormValue("IdEmp")),
		RFC: strings.TrimSpace(r.FormValue("RFC")),
		ErrRFC: "",
		Nombre: strings.TrimSpace(r.FormValue("Nombre")),
		ErrNombre: "",
		RazonSoc: strings.TrimSpace(r.FormValue("RazonSoc")),
		ErrRazonSoc: "",
		DirCalle: strings.TrimSpace(r.FormValue("DirCalle")),
		ErrDirCalle: "",
		DirCol: strings.TrimSpace(r.FormValue("DirCol")),
		ErrDirCol: "",
		DirEnt: strings.TrimSpace(r.FormValue("DirEnt")),
		ErrDirEnt: "",
		DirMun: strings.TrimSpace(r.FormValue("DirMun")),
		ErrDirMun: "",
		DirCp: strings.TrimSpace(r.FormValue("DirCp")),
		ErrDirCp: "",
		NumSuc: strings.TrimSpace(r.FormValue("NumSuc")),
		ErrNumSuc: "",
		OrgEmp: strings.TrimSpace(r.FormValue("OrgEmp")),
		ErrOrgEmp: "",
		OrgEmpOtro: strings.TrimSpace(r.FormValue("OrgEmpOtro")),
		ErrOrgEmpOtro: "",
		OrgEmpReg: strings.TrimSpace(r.FormValue("OrgEmpReg")),
		ErrOrgEmpReg: "",
		Url	: strings.TrimSpace(r.FormValue("Url")),
		ErrUrl:	"",
		PartLinea: partlinea,
		ExpComer: expcomer,
		Desc: strings.TrimSpace(r.FormValue("Desc")),
	}
	if valida {
		var ef bool
		ef = false
		if fd.RFC == "" || !model.ValidRfc.MatchString(fd.RFC) {
			fd.ErrRFC = "invalid"
			ef = true
		}
		if fd.Nombre == "" || !model.ValidSimpleText.MatchString(fd.Nombre) {
			fd.ErrNombre = "invalid"
			ef = true
		}
		if fd.RazonSoc == "" || !model.ValidSimpleText.MatchString(fd.RazonSoc) {
			fd.ErrRazonSoc = "invalid"
			ef = true
		}
		if fd.DirEnt == "" || !model.ValidSimpleText.MatchString(fd.DirEnt) {
			fd.ErrDirEnt = "invalid"
			ef = true
		}
		if fd.DirMun == "" || !model.ValidSimpleText.MatchString(fd.DirMun) {
			fd.ErrDirMun = "invalid"
			ef = true
		}
		if fd.DirCalle == "" || !model.ValidSimpleText.MatchString(fd.DirCalle) {
			fd.ErrDirCalle = "invalid"
			ef = true
		}
		if fd.DirCol == "" || !model.ValidSimpleText.MatchString(fd.DirCol) {
			fd.ErrDirCol = "invalid"
			ef = true
		}
		if fd.DirCp == "" || !model.ValidCP.MatchString(fd.DirCp) {
			fd.ErrDirCp = "invalid"
			ef = true
		}
		if fd.NumSuc != "" && !model.ValidNum.MatchString(fd.NumSuc) {
			fd.ErrNumSuc = "invalid"
			ef = true
		}
		if fd.OrgEmp != "" && !model.ValidSimpleText.MatchString(fd.OrgEmp) {
			fd.ErrOrgEmp = "invalid"
			ef = true
		}
		if fd.OrgEmpOtro != "" && !model.ValidSimpleText.MatchString(fd.OrgEmpOtro) {
			fd.ErrOrgEmpOtro = "invalid"
			ef = true
		}
		if fd.OrgEmpReg != "" && !model.ValidSimpleText.MatchString(fd.OrgEmpReg) {
			fd.ErrOrgEmpReg = "invalid"
			ef = true
		}
		if fd.Url != "" && !model.ValidUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}
		/*
		if fd.Desc != "" && !validSimpleText.MatchString(fd.Desc) {
			fd.ErrDesc = "invalid"
			ef = true
		}
		*/

		if ef {
			fd.Entidades = model.ListEnt(c, strings.TrimSpace(r.FormValue("DirEnt")))
			fd.Organismos = ListOrg(c, fd.OrgEmp)
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormDataEmp"] = fd
            tc["Page"] = model.PageMeta{"lea", "Empresa"}
            pageTpl = model.PrepareTemplate("empresa", "empresa")
            if err := pageTpl.ExecuteTemplate(w, "pageRoot", tc); err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
            }
			return fd, false
		}
	}
	return fd, true
}

func fillEmpresa(r *http.Request) model.Empresa {
	partlinea, _ := strconv.Atoi(r.FormValue("PartLinea"))
	expcomer, _ := strconv.Atoi(r.FormValue("ExpComer"))
	e := model.Empresa {
		IdEmp:		strings.TrimSpace(r.FormValue("IdEmp")),
		RFC:		strings.TrimSpace(r.FormValue("RFC")),
		Nombre:		strings.TrimSpace(r.FormValue("Nombre")),
		RazonSoc:	strings.TrimSpace(r.FormValue("RazonSoc")),
		DirCalle:	strings.TrimSpace(r.FormValue("DirCalle")),
		DirCol:		strings.TrimSpace(r.FormValue("DirCol")),
		DirEnt:		strings.TrimSpace(r.FormValue("DirEnt")),
		DirMun:		strings.TrimSpace(r.FormValue("DirMun")),
		DirCp:		strings.TrimSpace(r.FormValue("DirCp")),
		NumSuc:		strings.TrimSpace(r.FormValue("NumSuc")),
		OrgEmp:		strings.TrimSpace(r.FormValue("OrgEmp")),
		OrgEmpOtro:	strings.TrimSpace(r.FormValue("OrgEmpOtro")),
		OrgEmpReg:	strings.TrimSpace(r.FormValue("OrgEmpReg")),
		Url:		strings.TrimSpace(r.FormValue("Url")),
		PartLinea:  partlinea,
		ExpComer:	expcomer,
		Desc:		strings.TrimSpace(r.FormValue("Desc")),
		FechaHora:	time.Now().Add(time.Duration(model.GMTADJ)*time.Second),
		Status:		true,
	}
	return e
}

func empresaToForm(e model.Empresa) FormDataEmp {
	fe := FormDataEmp {
		IdEmp:		e.IdEmp,
		RFC:		e.RFC,
		Nombre:		e.Nombre,
		RazonSoc:	e.RazonSoc,
		DirCalle:	e.DirCalle,
		DirCol:		e.DirCol,
		DirEnt:		e.DirEnt,
		DirMun:		e.DirMun,
		DirCp:		e.DirCp,
		NumSuc:		e.NumSuc,
		OrgEmp:		e.OrgEmp,
		OrgEmpOtro:	e.OrgEmpOtro,
		OrgEmpReg:	e.OrgEmpReg,
		PartLinea:  e.PartLinea,
		ExpComer:	e.ExpComer,
		Desc:		e.Desc,
		Url:		e.Url,
	}
	return fe
}
