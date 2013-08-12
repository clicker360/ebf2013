package site

import (
    "appengine"
	"html/template"
    "net/http"
	"sess"
)

func init() {
    http.HandleFunc("/r/mg", materiales)
    http.HandleFunc("/r/dmg", descargamg)
}

func materiales(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		materialesTpl.Execute(w, tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}


func descargamg(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); ok {
		// Stream de materiales
		//w.Header().Set("Content-type", "application/octet-stream")
		//w.Header().Set("Content-disposition", "attachment; filename=test.zip")
		http.Redirect(w, r, "/ElBuenFin_Materiales.zip", http.StatusFound)

	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

var materialesTpl = template.Must(template.ParseFiles("templates/materiales.html"))
