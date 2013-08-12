package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
	"sortutil"
	"strconv"
    "model"
    "time"
)

type wssucursal struct {
	IdSuc       string	`json:"id"`
	Sucursal    string	`json:"sucursal"`
	DirCalle	string	`json:"calle"`
	DirCol		string	`json:"col"`
	DirEnt		string	`json:"entidad"`
	DirMun		string	`json:"municipio"`
	Latitud		float64	`json:"lat"`
	Longitud	float64	`json:"lng"`
}

func init() {
    http.HandleFunc("/wssucs", ShowEmpSucs)
}

func ShowEmpSucs(w http.ResponseWriter, r *http.Request) {
	var timetolive = 7200 //seconds
	c := appengine.NewContext(r)
	var b []byte
	if item, err := memcache.Get(c, "sucs_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		emsucs := model.GetEmpSucursales(c, r.FormValue("id"))
		wssucs := make([]wssucursal, len(*emsucs), cap(*emsucs))
		for i,v:= range *emsucs {
			wssucs[i].IdSuc = v.IdSuc
			wssucs[i].Sucursal = v.Nombre
			wssucs[i].DirCalle = v.DirCalle
			wssucs[i].DirCol = v.DirCol
			wssucs[i].DirEnt = v.DirEnt
			wssucs[i].Latitud,_ = strconv.ParseFloat(v.Geo1,64)
			wssucs[i].Longitud,_ = strconv.ParseFloat(v.Geo2,64)
		}
		sortutil.AscByField(wssucs, "Sucursal")
		b, _ = json.Marshal(wssucs)
		item := &memcache.Item{
			Key:   "sucs_"+r.FormValue("id"),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Error memcache.Add sucs_idemp : %v", err)
		}
	} else {
		//c.Infof("memcache retrieve sucs_idemp : %v", r.FormValue("id"))
		b = item.Value
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
