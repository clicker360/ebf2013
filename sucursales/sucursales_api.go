package sucursales

import (
    "appengine"
	"encoding/json"
	"sortutil"
    "net/http"
	"sess"
	"model"
	"time"
)

type WsSucursal struct{
	IdOft           string `json:"idoft"`
	IdEmp           string `json:"idemp"`
	IdSuc           string `json:"idsuc"`
	Sucursal        string `json:"sucursal"`
	FechaHora       time.Time `json:"timestamp"`
	Status		    string `json:"status"`
    Tel				string `json:"tel"`
	DirCalle		string `json:"calle"`
	DirCol			string `json:"colonia"`
	DirEnt			string `json:"entidad"`
	DirMun			string `json:"municipio"`
	DirCp			string `json:"cp"`
	GeoUrl			string `json:"geourl"`
	Geo1			string `json:"geo1"`
	Geo2			string `json:"geo2"`
	Geo3			string `json:"geo3"`
	Geo4			string `json:"geo4"`
	Ackn			string `json:"ackn"`
	Err			    *[]model.Errfield `json:"errors"`
}

func init() {
    http.HandleFunc("/r/wss/put", PutSucursal)
    //http.HandleFunc("/r/wss/post", PostSucursal)
    //http.HandleFunc("/r/wss/get", GetSucursal)
    //http.HandleFunc("/r/wss/gets", GetSucursales)
    //http.HandleFunc("/r/wss/del", DelSucursal)
}

func jsonDispatch(w http.ResponseWriter, out WsSucursal) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

func PutSucursal(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
    defer jsonDispatch(w, out)

	if _, ok := sess.IsSess(w, r, c); ok {
		out.Status = "noSession"
        return
    }

    if r.Method != "POST" {
		out.Status = "wrongMethod"
        return
    }

    out.IdEmp = r.FormValue("idemp")
    empresa := model.GetEmpresa(c, out.IdEmp)
    if empresa != nil {
        sucursal := fill(r)
        _, err := empresa.PutSuc(c, &sucursal)
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

func XDelOfSuc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
	out.IdSuc = r.FormValue("idsuc")
	out.IdOft = r.FormValue("idoft")
	err := model.DelOfertaSucursal(c, out.IdOft, out.IdSuc)
	if err != nil {
		out.Status = "notFound"
	} else {
		out.Status = "ok"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

func XShowOfSucursales(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ofsucs, _ := model.GetOfertaSucursales(c, r.FormValue("id"))
	wssucs := make([]WsSucursal, 0 ,len(*ofsucs))
	for i,v:= range *ofsucs {
		wssucs[i].IdOft = v.IdOft
		wssucs[i].IdSuc = v.IdSuc
		wssucs[i].IdEmp = v.IdEmp
		wssucs[i].Sucursal = v.Sucursal
		wssucs[i].FechaHora = v.FechaHora
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}

/*
	Listado de sucursales por empresa
*/
func XShowEmpSucs(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("IdEmp"))
	wssucs := make([]WsSucursal, len(*emsucs), cap(*emsucs))
	for i,v:= range *emsucs {
		wssucs[i].IdOft = ""
		wssucs[i].IdSuc = v.IdSuc
		wssucs[i].IdEmp = v.IdEmp
		wssucs[i].Sucursal = v.Nombre
		wssucs[i].FechaHora = v.FechaHora
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}

/*
	Listado de sucursales por empresa con la oferta marcada
*/
func XShowEmpSucursalOft(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("idemp"))
	ofsucs, _ := model.GetOfertaSucursales(c, r.FormValue("idoft"))
	wssucs := make([]WsSucursal, len(*emsucs), cap(*emsucs))
	for i,es:= range *emsucs {
		for _,os:= range *ofsucs {
			if os.IdSuc == es.IdSuc {
				wssucs[i].IdOft = os.IdOft
			}
		}
		wssucs[i].IdSuc = es.IdSuc
		wssucs[i].IdEmp = es.IdEmp
		wssucs[i].Sucursal = es.Nombre
		wssucs[i].FechaHora = es.FechaHora
	}
	sortutil.AscByField(wssucs, "Sucursal")

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}

/*
func validate(w http.ResponseWriter, r *http.Request, valida bool) (model.Errfield, bool){
    var ef bool
    ef = false
    if fd.Nombre == "" || !model.ValidSimpleText.MatchString(fd.Nombre) {
        fd.ErrNombre = "invalid"
        ef = true
    }
    if fd.Tel != "" && !model.ValidTel.MatchString(fd.Tel) {
        fd.ErrTel = "invalid"
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

    if ef {
        return fd, false
    }
	return fd, true
}
*/
