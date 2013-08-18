package site

import (
    "appengine"
	"appengine/memcache"
	"encoding/json"
    "net/http"
	"model"
	"sess"
)

type WsMunicipio struct {
	CveEnt		string `json:"cveent"`
	Entidad		string `json:"entidad"`
    Municipios *[]model.Municipio `json:"municipios"`
	Status		string `json:"status,omitempty"`
	Ackn		string `json:"ackn,omitempty"`
}

func init() {
    http.HandleFunc("/r/wsu/municipios", GetMunicipios)
    http.HandleFunc("/wsu/municipios", GetMunicipios)
}

/*
 Recibe CveMun para desplegar los municipios de un estado
*/
func GetMunicipios(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var b []byte
	var municipios *[]model.Municipio
    var out WsMunicipio
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        model.JsonDispatch(w, &out)
        return
    } else {
        out.Status = "notFound"
		if item, err := memcache.Get(c, "msp_"+r.FormValue("CveEnt")); err == memcache.ErrCacheMiss {
			if entidad, err := model.GetEntidad(c, r.FormValue("CveEnt")); err == nil {
				if municipios, _ = entidad.GetMunicipios(c); err == nil {
                    // Se guarda el json para caché de municipios
                    out.Municipios = municipios
                    out.CveEnt = entidad.CveEnt
                    out.Entidad = entidad.Entidad
                    out.Status = "ok"
                    //sortutil.AscByField(out.Municipios, "CveMun")
                    b, _ = json.Marshal(&out)
                    item := &memcache.Item{
                        Key:   "msp_"+r.FormValue("CveEnt"),
                        Value: b,
                    }
                    if err := memcache.Add(c, item); err == memcache.ErrNotStored {
                        c.Infof("item with key %q already exists", item.Key)
                    } else if err != nil {
                        c.Errorf("Error memcache.Add Municipio : %v", err)
                    }
                    c.Infof("CveMun generado: %v", item.Key)
				}
			}
		} else {
			c.Infof("Memcache activo: %v", item.Key)
            b = item.Value
            out.Status = "ok"
		}
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.Write(b)
	}
}
