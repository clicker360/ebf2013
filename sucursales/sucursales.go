package sucursales

import (
    "appengine"
    "appengine/datastore"
	"html/template"
    "net/http"
	//"fmt"
	"model"
	"sess"
)

func init() {
    http.HandleFunc("/r/suc", sucursales)
}

func sucursales(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	_, ok := sess.IsSess(w, r, c)
	if ok {
		q := datastore.NewQuery("Sucursal").Filter("IdEmp =", r.FormValue("IdEmp"))
		tpl, _ := template.New("Suc").Parse(SucListTpl)
		//fmt.Fprintf(w, "<div class=\"col-100PR first marg-U20pix\">")
		for i := q.Run(c); ; {
			var s model.Sucursal
			_, err := i.Next(&s)
			if err == datastore.Done {
				break
			}
			//Despliega sucursales
			tpl.Execute(w, s)
		}
		//fmt.Fprintf(w, "</div>")
	}
	return
}

const SucListTpl = `<div class="gridsubRow bg-Gry1"><a href="/sucursal?IdSuc={{.IdSuc}}&IdEmp={{.IdEmp}}">{{.Nombre}}</a><a href="/sucdel?IdSuc={{.IdSuc}}" class="button orange last marg-R5pix"><span>ELIMINAR</span></a> </div>`
