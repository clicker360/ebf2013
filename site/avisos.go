package site

import (
    "appengine"
	"html/template"
    "net/http"
	"sess"
)

func init() {
    http.HandleFunc("/r/avisos", Avisos)
}

func Avisos(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		avisoTpl.Execute(w, tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

var avisoTpl = template.Must(template.ParseFiles("templates/avisos.html"))
