package site

import (
    "appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
    "net/http"
	"model"
	"sess"
)

type WsOrganismo struct {
    Organismos *[]model.Organismo `json:"organismos"`
	Status		string `json:"status,omitempty"`
	Ackn		string `json:"ackn,omitempty"`
}

func init() {
    http.HandleFunc("/r/wsu/organismos", GetOrganismos)
    http.HandleFunc("/wsu/organismos", GetOrganismos)
}

/*
 Despliega los organismos
*/
func GetOrganismos(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var b []byte
	var organismos []model.Organismo
    var out WsOrganismo
	if _, ok := sess.IsSess(w, r, c); !ok {
		out.Status = "noSession"
        model.JsonDispatch(w, &out)
        return
    } else {
        out.Status = "notFound"
		if item, err := memcache.Get(c, "organismos"); err == memcache.ErrCacheMiss {
            q := datastore.NewQuery("Organismo")
            organismos = make([]model.Organismo, 0, 32)
            if _, err := q.GetAll(c, &organismos); err == nil {
                // Se guarda el json para caché de municipios
                out.Organismos = &organismos
                out.Status = "ok"
                b, _ = json.Marshal(&out)
                item := &memcache.Item{
                    Key:   "organismos",
                    Value: b,
                }
                if err := memcache.Add(c, item); err == memcache.ErrNotStored {
                    c.Infof("item with key %q already exists", item.Key)
                } else if err != nil {
                    c.Errorf("Error memcache.Add Organismos : %v", err)
                }
                c.Infof("Cache generado: %v", item.Key)
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
