package sucursales

import (
    "appengine"
	"encoding/json"
	"sortutil"
    "net/http"
    "strconv"
	"sess"
	"model"
	"time"
    //"fmt"
    //fmt.Fprintf(w, `b`)
)

type WsSucursal struct{
	IdSuc           string `json:"IdSuc"`
	IdEmp           string `json:"IdEmp"`
	Nombre          string `json:"Nombre"`
    Tel				string `json:"Tel"`
	DirCalle		string `json:"DirCalle"`
	DirCol			string `json:"DirCol"`
	DirEnt			string `json:"DirEnt"`
	DirMun			string `json:"DirMun"`
	DirCp			string `json:"DirCp"`
	GeoUrl			string `json:"GeoUrl,omitempty"`
	Geo1			string `json:"Geo1,omitempty"`
	Geo2			string `json:"Geo2,omitempty"`
	Geo3			string `json:"-"`
	Geo4			string `json:"-"`
	FechaHora       time.Time `json:"FechaHora"`
	Latitud		    float64 `json:"Latitud"`
	Longitud	    float64 `json:"Longitud"`
	Status		    string `json:"status"`
	Ackn		    string `json:"ackn,omitempty"`
	Errors		    map[string]bool `json:"errors,omitempty"`
}

func init() {
    /* CRUD */
    http.HandleFunc("/r/wss/put", PutSucursal)
    http.HandleFunc("/r/wss/post", PostSucursal)
    http.HandleFunc("/r/wss/get", GetSucursal)
    http.HandleFunc("/r/wss/gets", GetSucursales)
    http.HandleFunc("/r/wss/del", DelSucursal)
}

func PutSucursal(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
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
    out.IdEmp = r.FormValue("IdEmp")
    sucursal := fill(r)
    out.Errors, out.Status = validate(sucursal)
    if out.Status != "ok" {
        return
    }
    empresa := model.GetEmpresa(c, out.IdEmp)
    if empresa != nil {
        _, err := empresa.PutSuc(c, &sucursal)
        setWsSucursal(&out, sucursal)
        if err != nil {
            out.Status = "writeErr"
        } else {
            out.Status = "ok"
        }
    } else {
        out.Status = "notFound"
    }
    return
}

/* OJO FALTA HACER LA MODIFICACION DE SUCURSAL!!!
*/
func PostSucursal(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
    defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    if r.Method != "POST" {
		out.Status = "wrongMethod"
        return
    }
    out.IdEmp = r.FormValue("IdEmp")
    sucursal := fill(r)
    out.Errors, out.Status = validate(sucursal)
    if out.Status != "ok" {
        return
    }
    empresa := model.GetEmpresa(c, out.IdEmp)
    if empresa != nil {
        _, err := empresa.PutSuc(c, &sucursal)
        setWsSucursal(&out, sucursal)
        if err != nil {
            out.Status = "writeErr"
        } else {
            out.Status = "ok"
        }
    } else {
        out.Status = "notFound"
    }
    return
}

/*
	Detalle de sucursal. Regresa una sucursal por id sucursal
*/
func GetSucursal(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsSucursal
    defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        return
    }
    if r.Method != "GET" {
		out.Status = "wrongMethod"
        return
    }
	s := model.GetSuc(c, r.FormValue("IdSuc"))
    setWsSucursal(&out, *s)
    if s.IdEmp == "none" {
        out.Status = "notFound"
    } else {
        out.Status = "ok"
    }
    return
}

/*
	Regresa todas las sucursales por empresa
*/
func GetSucursales(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsSucursal
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        model.JsonDispatch(w, &out)
        return
    }
    if r.Method != "GET" {
		out.Status = "wrongMethod"
        model.JsonDispatch(w, &out)
        return
    }
	s := model.GetEmpSucursales(c, r.FormValue("IdEmp"))
	ws := make([]WsSucursal, len(*s), cap(*s))
	for i,v:= range *s {
        setWsSucursal(&ws[i], v)
    }
	sortutil.AscByField(ws, "Nombre")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, _ := json.Marshal(ws)
	w.Write(b)
}

func DelSucursal(w http.ResponseWriter, r *http.Request) {
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
    if err := model.DelSuc(c, r.FormValue("idsuc")); err != nil {
		out.Status = "notFound"
    }
	out.Status = "ok"
    return
}

func setWsSucursal(out *WsSucursal, s model.Sucursal) {
    out.IdSuc = s.IdSuc
    out.IdEmp = s.IdEmp
    out.Nombre = s.Nombre
    out.Tel = s.Tel
    out.DirCalle = s.DirCalle
    out.DirCol = s.DirCol
    out.DirEnt = s.DirEnt
    out.DirMun = s.DirMun
    out.DirCp = s.DirCp
    out.GeoUrl = s.GeoUrl
    out.Geo1 = s.Geo1
    out.Geo2 = s.Geo2
    out.Geo3 = s.Geo3
    out.Geo4 = s.Geo4
    out.FechaHora = s.FechaHora
    out.Latitud, _ = strconv.ParseFloat(s.Geo1, 64)
    out.Longitud, _ = strconv.ParseFloat(s.Geo2, 64)
}

func validate(s model.Sucursal) (map[string]bool, string) {
    errmsg := "ok"
    err := make(map[string]bool)
    if s.Nombre == "" || !model.ValidSimpleText.MatchString(s.Nombre) {
        err["Nombre"] = false
    }
    if s.Tel != "" && !model.ValidTel.MatchString(s.Tel) {
        err["Tel"] = false
    }
    if s.DirEnt == "" || !model.ValidSimpleText.MatchString(s.DirEnt) {
        err["DirEnt"] = false
    }
    if s.DirMun == "" || !model.ValidSimpleText.MatchString(s.DirMun) {
        err["DirMun"] = false
    }
    if s.DirCalle == "" || !model.ValidSimpleText.MatchString(s.DirCalle) {
        err["DirCalle"] = false
    }
    if s.DirCol == "" || !model.ValidSimpleText.MatchString(s.DirCol) {
        err["DirCol"] = false
    }
    if s.DirCp == "" || !model.ValidCP.MatchString(s.DirCp) {
        err["DirCp"] = false
    }
    for _, v := range err {
        if v == false {
            errmsg = "invalidInput"
        }
    }
	return err, errmsg
}
