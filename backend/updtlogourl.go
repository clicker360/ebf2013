package backend

import (
	"appengine"
    "appengine/datastore"
    "appengine/blobstore"
	"appengine/image"
	"net/http"
	"strconv"
	"model"
)

func init() {
    http.HandleFunc("/backend/updtlogourl", UpdateServingLogoUrl)
    //http.HandleFunc("/backend/updatesearch", RedirUpdateSearch)
}

func UpdateServingLogoUrl(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	const batch = 50
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
	q := datastore.NewQuery("EmpLogo").Offset(offset).Order("IdEmp").Limit(batch)
	n,_ := q.Count(c)
	for i := q.Run(c); ; {
		var e model.Image
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}

		// Se crea la URL para servir la oferta desde el CDN, si no se puede
		// se deja en blanco
		var imgprops image.ServingURLOptions
		imgprops.Secure = true
		imgprops.Size = 180
		imgprops.Crop = false
		if e.Sp4 == "" && e.IdEmp != "" {
			var blobkey appengine.BlobKey
			blob, err := blobstore.Create(c, "image/jpeg")
			if err != nil {
				c.Errorf("blobstore Create: %v", e.IdEmp)
			}
			_, err = blob.Write(e.Data)
			if err != nil {
				c.Errorf("blobstore Write: %v", e.IdEmp)
			}
			err = blob.Close()
			if err != nil {
				c.Errorf("blobstore Close: %v", e.IdEmp)
			}
			blobkey, err = blob.Key()
			if err != nil {
				c.Errorf("blobstore Key Gen: %v", e.IdEmp)
			}
			if url, err := image.ServingURL(c, blobkey, &imgprops); err != nil {
				c.Errorf("Cannot construct EmpLogo ServingURL : %v", e.IdEmp)
			} else {
				e.Sp3 = string(blobkey)
				e.Sp4 = url.String()
			}
			_, err = datastore.Put(c, key, &e)
			if err != nil {
				c.Errorf("PutEmpLogo(); Error al intentar actualizar Emplogo : %v", e.IdEmp)
			}
		}
	}
	c.Infof("UpdateServingLogoUrl() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
	return
}
