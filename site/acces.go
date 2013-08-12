package site

import (
    "appengine"
    "appengine/datastore"
	"appengine/urlfetch"
    "appengine/mail"
	"html/template"
	"strings"
	"bytes"
    "net/http"
    "net/url"
    "time"
	"fmt"
	"model"
	"sess"
)

func init() {
    http.HandleFunc("/r/acceso", Acceso)
    http.HandleFunc("/r/recupera", Recover)
    http.HandleFunc("/r/salir", Salir)
}

func Acceso(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	var st sess.Sess
	if _, ok := sess.IsSess(w, r, c); !ok {
		//fmt.Fprintf(w, "u:%s, p:%s", r.FormValue("u"), r.FormValue("p"))
		if(r.FormValue("u") != "" && r.FormValue("p") != "") {
			user := strings.TrimSpace(r.FormValue("u"))
			pass := strings.TrimSpace(r.FormValue("p"))
			/* validar usuario y pass */
			if(model.ValidEmail.MatchString(user) && model.ValidPass.MatchString(pass)) {
				q := datastore.NewQuery("Cta").Filter("Email =", user).Filter("Pass =", pass).Filter("Status =", true)
				if count, _ := q.Count(c); count != 0 {
					for t := q.Run(c); ; {
						var g model.Cta
						key, err := t.Next(&g)
						if err == datastore.Done {
							break
						}
						// Coincide contraseña, se activa una sesión
						_, _, err = sess.SetSess(w, c, key, g.Email, g.Nombre)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						// Redireccion
						http.Redirect(w, r, "/r/dash", http.StatusFound)
						return
					}
				}
			}
			st.User = r.FormValue("u")
			st.ErrMsg = "Usuario o contraseña incorrectos"
			st.ErrClass = "show"
		} else {
			st.User = r.FormValue("u")
			st.ErrMsg = "Proporcione usuario y contraseña"
			st.ErrClass = "show"
		}
	} else {
		// hay sesión
		http.Redirect(w, r, "/r/dash", http.StatusFound)
		return
	}
	tc := make(map[string]interface{})
	tc["Sess"] = st
    accesoErrorTpl.Execute(w, tc)
}

func Salir(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		s.Expiration = now.AddDate(-1,0,0)
		_, err := datastore.Put(c, datastore.NewKey(c, "Sess", s.User, 0, nil), &s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", "", "Thu, 18 Oct 2012 01:01:23 GMT;"))
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", "", "Thu, 18 Oct 2012 01:01:23 GMT;"))
	http.Redirect(w, r, "/r/registro", http.StatusFound)
}

func Recover(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); !ok {
		var email string = strings.TrimSpace(r.FormValue("Email"))
		//var rfc string = strings.TrimSpace(r.FormValue("RFC"))
		if email != "" && model.ValidEmail.MatchString(email) { // && rfc != "" && model.ValidRfc.MatchString(rfc) {
			// intenta buscar en la base un usuario con email y empresa
			if cta, err := model.GetCta(c, email); err == nil {
				if cta.Status {
					if (MailServer=="gmail") {
						var hbody bytes.Buffer
						var sender string
						if (appengine.AppID(c) == "ebfmxorg") {
							sender =  "El Buen Fin <contacto@elbuenfin.org>"
						} else {
							sender =  "El Buen Fin <ahuezo@clicker360.com>"
						}
						if err := mailRecoverTpl.Execute(&hbody, cta); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						// Coincide email y RFC, se manda correo con contraseña
						msg := &mail.Message{
							Sender:		sender,
							To:			[]string{cta.Email},
							Subject:	"Recuperación de contraseña / El Buen Fin",
							HTMLBody:	hbody.String(),
						}
						if err := mail.Send(c, msg); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						} else {
							http.Redirect(w, r, "/recoverok.html", http.StatusFound)
							return
						}
						//fmt.Fprintf(w, mailRecover, cta.Email, cta.Pass)
						return
					} else {
						client := urlfetch.Client(c)
						url := fmt.Sprintf("http://envia-m.mekate.com.mx/?Sender=%s&Tipo=Recupera&Email=%s&Nombre=%s&Pass=%s&AppId=ebfmxorg",
						"registro@elbuenfin.org",
						cta.Email,
						url.QueryEscape(cta.Nombre),
						url.QueryEscape(cta.Pass))
						r1, err := client.Get(url)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						if r1.StatusCode != 200 {
							http.Error(w, "Error de Transporte de Mail", http.StatusInternalServerError)
						}
						defer r1.Body.Close()
						http.Redirect(w, r, "/recoverok.html", http.StatusFound)
						return
					}
				} else {
					http.Redirect(w, r, "/nocta.html", http.StatusFound)
					return
				}
			}
		}
		http.Redirect(w, r, "/nocta.html", http.StatusFound)
		return
	} else {
		http.Redirect(w, r, "/r/dash", http.StatusFound)
	}
}

var accesoErrorTpl = template.Must(template.ParseFiles("templates/acceso_error.html"))
var mailRecoverTpl = template.Must(template.ParseFiles("templates/mail_recover.html"))
