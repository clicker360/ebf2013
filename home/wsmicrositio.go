package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
    "model"
    "time"
)

type micrositio struct {
	IdEmp		string	`json:"idemp"`
	IdImg		string	`json:"idimg"`
	Name		string	`json:"name"`
	Desc		string	`json:"desc"`
	Url         string	`json:"url"`
	Sp1         string	`json:"facebook"`
	Sp2         string	`json:"twitter"`
	Sp4         string	`json:"srvurl"`
}

func init() {
    http.HandleFunc("/wsmicrositio", ShowMicrositio)
}

func ShowMicrositio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var timetolive = 7200 //seconds
	c := appengine.NewContext(r)
	var b []byte
	var m micrositio
	if item, err := memcache.Get(c, "m_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		if e := model.GetEmpresa(c, r.FormValue("id")); e != nil {
			m.IdEmp = e.IdEmp
			m.Name = e.Nombre
			m.Desc = e.Desc
			if imgo := model.GetLogo(c, r.FormValue("id")); imgo != nil {
				m.IdImg = imgo.IdImg
				m.Url = imgo.Url
				m.Sp1 = imgo.Sp1
				m.Sp2 = imgo.Sp2
				m.Sp4 = imgo.Sp4
			}
		}

		b, _ = json.Marshal(m)
		item := &memcache.Item{
			Key:   "m_"+r.FormValue("id"),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Error memcache.Add m_idoft : %v", err)
			w.Write(b)
		}
	} else {
		//c.Infof("memcache retrieve m_idoft : %v", r.FormValue("id"))
		b = item.Value
	}
	w.Write(b)
}

