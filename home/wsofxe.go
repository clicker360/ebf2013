package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
	"sortutil"
    "model"
	"strconv"
    "time"
)

type wsoferta struct {
	IdEmp       string	`json:"idemp"`
	IdOft		string	`json:"idoft"`
	Oferta		string	`json:"oferta"`
	Descripcion	string	`json:"descripcion"`
	Enlinea		bool	`json:"enlinea"`
	Url			string	`json:"url"`
	EmpLogo		string	`json:"EmpLogo"`
	SrvUrl		string	`json:"srvurl"`
	BlobKey	appengine.BlobKey `json:"imgurl"`
}

func init() {
    http.HandleFunc("/wsofxe", ShowEmpOfertas)
}

func ShowEmpOfertas(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
	var timetolive = 3600 //seconds
	var batch = 12 // tamaño de pagina
	c := appengine.NewContext(r)
	var page int
	page,_ = strconv.Atoi(r.FormValue("pagina"))
	page = page-1
	if page < 1 {
		page = 0
	}
	wsbatch := make([]wsoferta, 0, batch)
	offset := page*batch
	var b []byte
	if item, err := memcache.Get(c, "ofxe_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		emofs := model.ListOf(c, r.FormValue("id"))
		wsofs := make([]wsoferta, len(*emofs), cap(*emofs))
		for i,v:= range *emofs {
			if now.After(v.FechaHoraPub) {
				if v.Oferta != "Nueva oferta" {
					wsofs[i].IdEmp = v.IdEmp
					wsofs[i].IdOft = v.IdOft
					wsofs[i].Oferta = v.Oferta
					wsofs[i].Descripcion = v.Descripcion
					wsofs[i].Enlinea = v.Enlinea
					wsofs[i].Url = v.Url
					if v.Promocion != "" {
						wsofs[i].EmpLogo = v.Promocion
					} else {
						wsofs[i].EmpLogo = "http://www.elbuenfin.org/imgs/imageDefault.png"
					}
					wsofs[i].SrvUrl = v.Codigo
					wsofs[i].BlobKey = v.BlobKey
				}
			}
		}
		sortutil.AscByField(wsofs, "Oferta")
		jb, _ := json.Marshal(wsofs)
		item := &memcache.Item{
			Key:   "ofxe_"+r.FormValue("id"),
			Value: jb,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Error memcache.Add ofxe_idemp : %v", err)
		}

		// se pagina la respuesta
		if(offset <= len(*emofs)-batch) {
			b, _ = json.Marshal(wsofs[offset:offset+batch])
		} else {
			if len(wsofs) > offset {
				b, _ = json.Marshal(wsofs[offset:])
			} else {
				b,_ = json.Marshal(make([]wsoferta, 0, batch))
			}
		}
	} else {
		//c.Infof("memcache retrieve sucs_idemp : %v", r.FormValue("id"))
		// se pagina la respuesta
		if err := json.Unmarshal(item.Value, &wsbatch); err != nil {
			c.Errorf("Unmarshaling wsbatch item: %v", err)
		}
		// se pagina la respuesta
		if(offset <= len(wsbatch)-batch) {
			b, _ = json.Marshal(wsbatch[offset:offset+batch])
		} else {
			if len(wsbatch) > offset {
				b, _ = json.Marshal(wsbatch[offset:])
			} else {
				b,_ = json.Marshal(make([]wsoferta, 0, batch))
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
