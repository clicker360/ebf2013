package backend

import (
	"appengine"
    "appengine/datastore"
	"net/http"
	"strconv"
	"model"
)

func init() {
    //http.HandleFunc("/backend/updateempnm", UpdateEmpNm)
    //http.HandleFunc("/backend/updatesearch", RedirUpdateSearch)
}

func UpdateEmpNm(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	const batch = 300
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
	q := datastore.NewQuery("Empresa").Offset(offset).Limit(batch)
	n,_ := q.Count(c)
	for i := q.Run(c); ; {
		var e model.Empresa
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		var ua model.Cta
		if err := datastore.Get(c, key.Parent(), &ua); err != nil {
			c.Errorf("Get Cta Key; Error al intentar leer key.Parent() de Empresa : %v", e.IdEmp)
		} else {
			if _, err := ua.PutEmpresa(c, &e); err != nil {
				c.Errorf("PutEmpresa(); Error al intentar actualizar empresa : %v", e.IdEmp)
			}
		}
	}
	c.Infof("UpdateEmpNm() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
	return
}
