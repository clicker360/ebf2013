package backend

import (
	"appengine"
    "appengine/datastore"
    "appengine/blobstore"
	"appengine/image"
	"net/http"
	"strconv"
	"model"
	"io/ioutil"
	"fmt"
)

func init() {
    http.HandleFunc("/backend/scanblob", ScanOfertaBlob)
    //http.HandleFunc("/backend/updatesearch", RedirUpdateSearch)
}

func ScanOfertaBlob(w http.ResponseWriter, r *http.Request) {
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
		if e.BlobKey != "none" {
			reader := blobstore.NewReader(c, e.BlobKey)
			if _, err := ioutil.ReadAll(reader); err != nil {
				fmt.Fprintf(w, "Error en idoft: %s, idemp: %s, blobkey: %v, Fecha: %v\n", e.IdOft, e.IdEmp, string(e.BlobKey), e.FechaHora)
				e.BlobKey = "none"
				_, err = datastore.Put(c, key, &e)
			}
		}
	}
	fmt.Fprintf(w, "Batch: %d, count: %d, from  %d to %d\n", page, n, offset, offset+n)
	return
}
