package site

import (
    "appengine"
    "appengine/datastore"
	"appengine/urlfetch"
	"appengine/mail"
    "appengine/user"
    "net/http"
    "net/url"
	"html/template"
	"crypto/md5"
	"bytes"
    "time"
	"fmt"
	"io"
	"sharded_counter"
	"model"
	"sess"
)

type urlCfm struct {
	Md5			string
	Nombre		string
	Apellidos	string
	Email		string
	FechaHora	time.Time
	Llave		string
	AppId		string
}

const MailServer = "gmail"
//const MailServer = "clicker"

func init() {
    http.HandleFunc("/r/registrar", Registrar)
    //http.HandleFunc("/registrar", pendienteVerifica)
    http.HandleFunc("/r/c", ConfirmaCodigo)
}

func Registrar(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	s, ok := sess.IsSess(w, r, c)
	if ok {
		http.Redirect(w, r, "/r/cta", http.StatusFound)
		return
	}
	fd, valid := ctaForm(w, r, s, true, registroTpl)
	if valid {
		u, err := model.GetCta(c, fd.Email)
		ctaFill(r, u)
		if err != nil {
			// No hay Cuenta registrada
			u.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
			u.Status = false

			// Generar código de confirmación distindo cada vez. Md5 del email + fecha-hora
			h := md5.New()
			io.WriteString(h, fmt.Sprintf("%s%s%s%s", time.Now().Add(time.Duration(model.GMTADJ)*time.Second), u.Email, u.Pass, model.RandId(12)))
			u.CodigoCfm = fmt.Sprintf("%x", h.Sum(nil))
		}

		//Si hay estatus es que ya existe
		if(u.Status == false) {

			// Se agrega la cuenta sin activar para realizar el proceso de código de confirmación
			if u, err = model.PutCta(c, u); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			/* No se ha activado, por tanto se inicia el proceso de código de verificación */
			m := urlCfm{
				Md5:		u.CodigoCfm,
				Nombre:		u.Nombre,
				Email:		u.Email,
				FechaHora:	time.Now().Add(time.Duration(model.GMTADJ)*time.Second),
				Llave:		u.Key(c).Encode(),
				AppId:		appengine.AppID(c),
			}
			var hbody bytes.Buffer
			var sender string
			if (appengine.AppID(c) == "ebfmxorg") {
				sender =  "El Buen Fin <contacto@elbuenfin.org>"
			} else {
				sender =  "El Buen Fin <ahuezo@clicker360.com>"
			}
			// Envia código activación 
			if err := mailActivationCodeTpl.Execute(&hbody, m); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if (MailServer == "gmail") {
				msg := &mail.Message{
					Sender:  sender,
					To:      []string{m.Email},
					Subject: "Codigo de Activación de Registro / El Buen Fin en línea",
					HTMLBody: hbody.String(),
				}
				if err := mail.Send(c, msg); err != nil {
					/* Problemas para enviar el correo NOK */
					http.Error(w, err.Error(), http.StatusInternalServerError)
					http.Redirect(w, r, "/", http.StatusFound)
				} else {
					if err := activationMessageTpl.ExecuteTemplate(w, "codesend", m); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				// ************************************************************
				// Si el hay usuario admin se despliega el código de activación
				// ************************************************************
				if gu := user.Current(c); gu != nil {
					if err := mailActivationCodeTpl.Execute(w, m); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			} else {
				client := urlfetch.Client(c)
				url := fmt.Sprintf("http://envia-m.mekate.com.mx/?Sender=%s&Tipo=Codigo&Md5=%s&Llave=%s&Email=%s&Nombre=%s&AppId=ebfmxorg",
				"registro@elbuenfin.org",
				m.Md5,
				url.QueryEscape(m.Llave),
				m.Email,
				url.QueryEscape(m.Nombre))
				r1, err := client.Get(url)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//fmt.Fprintf(w, "HTTP GET returned status %v", r.Status)

				if r1.StatusCode != 200 {
					http.Error(w, "Error de Transporte de Mail", http.StatusInternalServerError)
					return
				} else {
					if err := activationMessageTpl.ExecuteTemplate(w, "codesend", m); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				defer r1.Body.Close()
			}
		} else {
			if err := registroErrorTpl.Execute(w, nil); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func pendienteVerifica(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	s, ok := sess.IsSess(w, r, c)
	if ok {
		http.Redirect(w, r, "/r/cta", http.StatusFound)
		return
	}
	fd, valid := ctaForm(w, r, s, true, registroTpl)
	if valid {
		u, err := model.GetCta(c, fd.Email)
		ctaFill(r, u)
		if err != nil {
			// No hay Cuenta registrada
			u.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
			u.Status = false
			u.CodigoCfm = "Verificar"

			// Generar código de confirmación distindo cada vez. Md5 del email + fecha-hora
			h := md5.New()
			io.WriteString(h, fmt.Sprintf("%s%s%s%s", time.Now().Add(time.Duration(model.GMTADJ)*time.Second), u.Email, u.Pass, model.RandId(12)))
			u.CodigoCfm = fmt.Sprintf("%x", h.Sum(nil))
		}

		u.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
		u.Status = true
		u.CodigoCfm = "Verificar"
		// Se agrega la cuenta sin activar para realizar el proceso de código de confirmación
		if u, err = model.PutCta(c, u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		/* Prende la sesion */
		_, _, err = sess.SetSess(w, c, u.Key(c), u.Email, u.Nombre)

		// avisa del éxito independientemente del correo
		if err := activationMessageTpl.ExecuteTemplate(w, "verify", u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func ConfirmaCodigo(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	md5 := r.FormValue("m")
    key, _ := datastore.DecodeKey(r.FormValue("c"))
	var g model.Cta
    if err := datastore.Get(c, key, &g); err != nil {
		if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        return
    }

	/* Se verifica el código de confirmación */
	if(g.CodigoCfm == md5 && g.Status == false) {
		// Si se confirma el md5 la cuenta admin se le asigna un folio y se activa el status
		if err := sharded_counter.Increment(c, "cuenta_admin"); err == nil {
			if folio, err := sharded_counter.Count(c, "cuenta_admin"); err == nil {
				g.Folio = folio
				g.Status = true
				g.CodigoCfm = "Confirmado"
				_, err := datastore.Put(c, key, &g)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				/* Prende la sesion */
				_, _, err = sess.SetSess(w, c, key, g.Email, g.Nombre)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if (MailServer=="gmail") {
					// Envia código activación 
					var hbody bytes.Buffer
					var sender string
					if (appengine.AppID(c) == "ebfmxorg") {
						sender =  "El Buen Fin <contacto@elbuenfin.org>"
					} else {
						sender =  "El Buen Fin <ahuezo@clicker360.com>"
					}
					if err := mailAvisoActivacionTpl.Execute(&hbody, g); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					msg := &mail.Message{
							Sender:  sender,
							To:      []string{g.Email},
							Subject: "Cuenta Activada / El Buen Fin en línea",
							HTMLBody:	hbody.String(),
					}
					if err := mail.Send(c, msg); err != nil {
						// Problemas para enviar el correo NOK 
						http.Error(w, err.Error(), http.StatusInternalServerError)
						http.Redirect(w, r, "/", http.StatusFound)
					}
					// avisa del éxito independientemente del correo
					if err := activationMessageTpl.ExecuteTemplate(w, "confirm", g); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				} else {
					client := urlfetch.Client(c)
					url := fmt.Sprintf("http://envia-m.mekate.com.mx/?Sender=%s&Tipo=Aviso&Email=%s&Nombre=%s&Pass=%s&AppId=ebfmxorg",
					"registro@elbuenfin.org",
					g.Email,
					url.QueryEscape(g.Nombre),
					url.QueryEscape(g.Pass))
					r1, err := client.Get(url)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if r1.StatusCode != 200 {
						http.Error(w, "Error de Transporte de Mail", http.StatusInternalServerError)
					}
					if err := activationMessageTpl.ExecuteTemplate(w, "confirm", g); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					defer r1.Body.Close()
				}
			} else {
				// El Folio no es seguro, se deshecha la operación o se encola
				if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		} else {
			// El Folio no es seguro, se deshecha la operación o se encola 
			if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	} else {
		if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var registroTpl = template.Must(template.ParseFiles("templates/registro.html")) //, "templates/login.html"))
var registroErrorTpl = template.Must(template.ParseFiles("templates/registro_aviso.html"))
var mailActivationCodeTpl = template.Must(template.ParseFiles("templates/activation_code.html"))
var activationMessageTpl = template.Must(template.ParseFiles("templates/activation_result.html"))
var mailAvisoActivacionTpl = template.Must(template.ParseFiles("templates/mail_aviso_activacion.html"))
