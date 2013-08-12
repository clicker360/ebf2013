package backend

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
	"html/template"
    "net/http"
	"fmt"
	"model"
	"sess"
)

func init() {
    //http.HandleFunc("/r/backend", GaeLogin)
    //http.HandleFunc("/r/listausuarios", ListaUsuarios)
    //http.HandleFunc("/r/listasesiones", ListaSesiones)
    //http.HandleFunc("/r/test", test)
}

func test(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    if u := user.Current(c); u != nil {
		if ck, err := r.Cookie("ebfmex-pub-sessid-ua"); err == nil {
			fmt.Fprintf(w, "Nombre: %q, Valor: %q\n", ck.Name, ck.Value);
		} else {
	     http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cr, err := r.Cookie("ebfmex-pub-sesscontrol-ua"); err == nil {
			fmt.Fprintf(w, "Nombre: %q, Valor: %q\n", cr.Name, cr.Value);
		} else {
	     http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListaUsuarios(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	/* Verifica si el usuario es interno */
    if u := user.Current(c); u != nil {
		q := datastore.NewQuery("Cta").Order("-FechaHora").Limit(10)
		usuarios := make([]model.Cta, 0, 10)
		if _, err := q.GetAll(c, &usuarios); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := listUsersTpl.Execute(w, usuarios); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ListaSesiones(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	/* Verifica si el usuario es interno */
    if u := user.Current(c); u != nil {
		q := datastore.NewQuery("Sess").Limit(10)
		s := make([]sess.Sess, 0, 10)
		if _, err := q.GetAll(c, &s); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := listSessTpl.Execute(w, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GaeLogin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/* Autenticación de usuario interno */
	u := user.Current(c)
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    //http.Redirect(w, r, "/listausuarios", http.StatusFound)
}

var listSessTpl = template.Must(template.ParseFiles("templates/list_sess.html"))
var listUsersTpl = template.Must(template.ParseFiles("templates/list_users.html"))
