package site

import (
    "appengine"
    "appengine/datastore"
	"encoding/json"
	"strings"
	"sortutil"
    "net/http"
	"strconv"
	"sess"
	"model"
    "time"
//	"fmt"
)

type WsEmpresa struct {
	IdEmp			string `json:"idemp"`
	RFC				string `json:"rfc"`
	Nombre			string `json:"nombre"`
	RazonSoc		string `json:"razonsocial"`
	DirCalle		string `json:"calle"`
	DirCol			string `json:"colonia"`
	DirEnt			string `json:"entidad"`
	DirMun			string `json:"municipio"`
	DirCp			string `json:"cp"`
	NumSuc			string `json:"numsuc"`
	OrgEmp			string `json:"organismo"`
	OrgEmpOtro		string `json:"otro_organismo,omitempty"`
	OrgEmpReg		string `json:"registro_empresarial,omitempty"`
	Url				string `json:"url,omitempty"`
	PartLinea		int `json:"partlinea"`
	ExpComer		int `json:"expcomer"`
	Desc			string `json:"descripcion"`
	Status		    string `json:"status"`
	Ackn		    string `json:"ackn,omitempty"`
	Errors		    map[string]bool `json:"errors,omitempty"`
}

func init() {
    /* CRUD */
    http.HandleFunc("/r/wse/put", PutEmpresa)
    http.HandleFunc("/r/wse/post", PostEmpresa)
    http.HandleFunc("/r/wse/get", GetEmpresa)
    http.HandleFunc("/r/wse/gets", GetEmpresas)
    http.HandleFunc("/r/wse/del", DelEmpresa)
}

func PutEmpresa(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsEmpresa
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    } else {
        // se obtiene el detalle de cta
        u, _ := model.GetCta(c, s.User)

        // PUT
        if r.Method != "POST" {
            out.Status = "wrongMethod"
            return
        }

        // se construye la estructura temporal de la empresa nueva
        // y se valida
        eTmp := fill(r)

        out.Errors, out.Status = validate(eTmp)
        if out.Status != "ok" {
            return
        }

        // Se intenta añadir la estructura nueva
        empresa, err := u.NewEmpresa(c, &eTmp)
        setWsEmpresa(&out, *empresa)
        if err != nil {
            out.Status = err.Error()
            return
        }
        err1 := model.PutCtaEmp(c, eTmp.IdEmp, u.Email, u.EmailAlt)
        if err1 != nil {
            out.Status = err1.Error()
            return
        }
    }
    out.Status = "ok"
    return
}

func PostEmpresa(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsEmpresa
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
        eTmp := fill(r)

        out.Errors, out.Status = validate(eTmp)
        if out.Status != "ok" {
            return
        }

        // se obtiene la empresa existente a actualizar
        _, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
        if err != nil {
            out.Status = err.Error()
            return
        }
        _, err = u.PutEmpresa(c, &eTmp)
        setWsEmpresa(&out, eTmp)
        if err == datastore.ErrNoSuchEntity {
            out.Status = "notFound"
            return
        }
    }
    out.Status = "ok"
    return
}

/*
	Detalle de empresa. Regresa una empresa por id de cuenta
*/
func GetEmpresa(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsEmpresa
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
        e, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
        if err != nil {
            out.Status = "notFound"
        } else {
            setWsEmpresa(&out, *e)
            out.Status = "ok"
        }
    }
    return
}

/*
	Empresas de una cta
*/
func GetEmpresas(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsEmpresa
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
        e :=listEmp(c, u)
        ws := make([]WsEmpresa, len(*e), cap(*e))
        for i,v:= range *e {
            setWsEmpresa(&ws[i], v)
        }
        out.Status = "ok"
        sortutil.AscByField(ws, "Nombre")
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        b, _ := json.Marshal(ws)
        w.Write(b)
    }
}

func DelEmpresa(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsEmpresa
    defer model.JsonDispatch(w, &out)
	if s, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    } else {
	    u, _ := model.GetCta(c, s.User)
        // DELETE
        if r.Method != "GET" {
            out.Status = "wrongMethod"
            return
        }
        if err := u.DelEmpresa(c, r.FormValue("IdEmp")); err != nil {
            out.Status = "notFound"
        }
    }
	out.Status = "ok"
    return
}

func fill(r *http.Request) model.Empresa {
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

func setWsEmpresa(out *WsEmpresa, e model.Empresa) {
    out.IdEmp = e.IdEmp
    out.RFC = e.RFC
    out.Nombre = e.Nombre
    out.RazonSoc = e.RazonSoc
    out.DirCalle = e.DirCalle
    out.DirCol = e.DirCol
    out.DirEnt = e.DirEnt
    out.DirMun = e.DirMun
    out.DirCp = e.DirCp
    out.NumSuc = e.NumSuc
    out.OrgEmp = e.OrgEmp
    out.OrgEmpOtro = e.OrgEmpOtro
    out.OrgEmpReg = e.OrgEmpReg
    out.PartLinea = e.PartLinea
    out.ExpComer = e.ExpComer
    out.Desc = e.Desc
    out.Url = e.Url
}

func validate(e model.Empresa) (map[string]bool, string) {
    errmsg := "ok"
    err := make(map[string]bool)
    if e.RFC == "" || !model.ValidRfc.MatchString(e.RFC) {
        err["RFC"] = false
    }
	if e.Nombre == "" || !model.ValidSimpleText.MatchString(e.Nombre) {
        err["Nombre"] = false
	}
	if e.RazonSoc == "" || !model.ValidSimpleText.MatchString(e.RazonSoc) {
        err["RazonSoc"] = false
	}
    if e.DirEnt == "" || !model.ValidSimpleText.MatchString(e.DirEnt) {
        err["DirEnt"] = false
    }
    if e.DirMun == "" || !model.ValidSimpleText.MatchString(e.DirMun) {
        err["DirMun"] = false
    }
    if e.DirCalle == "" || !model.ValidSimpleText.MatchString(e.DirCalle) {
        err["DirCalle"] = false
    }
    if e.DirCol == "" || !model.ValidSimpleText.MatchString(e.DirCol) {
        err["DirCol"] = false
    }
    if e.DirCp == "" || !model.ValidCP.MatchString(e.DirCp) {
        err["DirCp"] = false
    }
    if e.NumSuc != "" && !model.ValidNum.MatchString(e.NumSuc) {
        err["NumSuc"] = false
    }
    if e.OrgEmp != "" && !model.ValidSimpleText.MatchString(e.OrgEmp) {
        err["OrgEmp"] = false
    }
    if e.OrgEmpOtro != "" && !model.ValidSimpleText.MatchString(e.OrgEmpOtro) {
        err["OrgEmpOtro"] = false
    }
    if e.OrgEmpReg != "" && !model.ValidSimpleText.MatchString(e.OrgEmpReg) {
        err["OrgEmpReg"] = false
    }
    if e.Url != "" && !model.ValidUrl.MatchString(e.Url) {
        err["Url"] = false
    }
    for _, v := range err {
        if v == false {
            errmsg = "invalidInput"
        }
    }
	return err, errmsg
}
