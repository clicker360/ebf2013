package backend

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "model"
	//"strings"
	"strconv"
)


func init() {
    //http.HandleFunc("/backend/updtctaemp", updateCtaEmpresa)
    //http.HandleFunc("/backend/updtctaemp", updateEmpCta)
}

func updateCtaEmpresa(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	const batch = 100

	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
    q := datastore.NewQuery("Cta").Offset(offset).Order("-FechaHora").Limit(batch)
	n,_ := q.Count(c)
	for cursor := q.Run(c); ; {
		var cta model.Cta
		key, err := cursor.Next(&cta)
		if err == datastore.Done  {
			break
		}

		q2 := datastore.NewQuery("Empresa").Ancestor(key).Limit(200)
		for cursor := q2.Run(c); ; {
			var emp model.Empresa
			_, err := cursor.Next(&emp)
			if err == datastore.Done  {
				break
			}
			var ce model.CtaEmpresa
			ce.IdEmp = emp.IdEmp
			ce.Email = cta.Email
			ce.EmailAlt = cta.EmailAlt
			if emp.IdEmp != "" { 
				_, err1 := datastore.Put(c, datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil), &ce)
				if err1 != nil {
					c.Errorf("PutCtaEmpresa(); Error al intentar actualizar CtaEmpresa : %v", ce.IdEmp)
				}
			} else {
				c.Errorf("PutCtaEmpresa(); IdEMP vacio!!! : %v", ce.Email)
			}
		}
	}
	c.Infof("UpdateServingLogoUrl() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
}

func updateEmpCta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	const batch = 300
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
	q := datastore.NewQuery("Empresa").Order("-FechaHora").Offset(offset).Limit(batch)
	n,_ := q.Count(c)
	for i := q.Run(c); ; {
		var e model.Empresa
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		var ua model.Cta
		var ce model.CtaEmpresa
		ce.IdEmp = e.IdEmp
		if err := datastore.Get(c, key.Parent(), &ua); err != nil {
			c.Errorf("Get Cta Key; Error al intentar leer key.Parent() de Empresa : %v", ce.IdEmp)
		} else {
			ce.Email = ua.Email
			ce.EmailAlt = ua.EmailAlt
			if ua.Email != "" {
				_, err1 := datastore.Put(c, datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil), &ce)
				if err1 != nil {
					c.Errorf("PutCtaEmpresa(); Error al intentar actualizar CtaEmpresa : %v", ce.IdEmp)
					}
			} else {
				c.Errorf("PutCtaEmpresa(); Email vacio!!! : %v", ce.IdEmp)
			}
		}
	}
	c.Infof("UpdateServingLogoUrl() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
	return
}
