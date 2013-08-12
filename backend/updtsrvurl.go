package backend

import (
	"appengine"
    "appengine/datastore"
	"appengine/image"
	"net/http"
	"strconv"
	"model"
)

func init() {
    http.HandleFunc("/backend/updtsrvurl", UpdateServingUrl)
}

func UpdateServingUrl(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	const batch = 300
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
	q := datastore.NewQuery("Oferta").Offset(offset).Order("-FechaHora").Limit(batch)
	n,_ := q.Count(c)
	for i := q.Run(c); ; {
		var e model.Oferta
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}

		// Se crea la URL para servir la oferta desde el CDN, si no se puede
		// se deja en blanco
		var imgprops image.ServingURLOptions
		imgprops.Secure = true
		imgprops.Size = 400
		imgprops.Crop = false
		if e.BlobKey != "none" && e.Codigo == "" {
			if url, err := image.ServingURL(c, e.BlobKey, &imgprops); err != nil {
				c.Errorf("Cannot construct ServingURL : %v", e.IdOft)
				e.Codigo = ""
			} else {
				e.Codigo = url.String()
			}

			//c.Errorf("Get Cta Key; Error al intentar leer key.Parent() de Empresa : %v", e.IdEmp)
			_, err = datastore.Put(c, key, &e)
			if err != nil {
				c.Errorf("PutEmpresa(); Error al intentar actualizar empresa : %v", e.IdEmp)
			}
		}
	}
	c.Infof("UpdateServingUrl() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
	return
}
