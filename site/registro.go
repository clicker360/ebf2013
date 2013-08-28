package site

import (
    "appengine"
    "appengine/datastore"
	"appengine/urlfetch"
	"appengine/mail"
    "net/http"
    "net/url"
	"html/template"
	"bytes"
    "time"
	"fmt"
	"model"
	"sess"
)

const MailServer = "gmail"
const MailSender = "El Buen Fin <contacto@elbuenfin.org>"
const MailSubjectConfirmed = "Cuenta Activada / El Buen Fin en línea"

var mailActivationResultTpl = template.Must(template.ParseFiles("layout/reg_activation_result.html"))
var mailAvisoActivacionTpl = template.Must(template.ParseFiles("layout/mail_activation_welcome.html"))

type urlCfm struct {
	Md5			string
	Nombre		string
	Apellidos	string
	Email		string
	FechaHora	time.Time
	Llave		string
	AppId		string
}

func init() {
    http.HandleFunc("/r/c", ConfirmaCodigo)
}

func ConfirmaCodigo(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	md5 := r.FormValue("m")
    key, _ := datastore.DecodeKey(r.FormValue("c"))
	var g model.Cta
    if err := datastore.Get(c, key, &g); err != nil {
		if err := mailActivationResultTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        return
    }

	/* Se verifica el código de confirmación */
	if(g.CodigoCfm == md5 && g.Status == false) {
		// Si se confirma el md5 la cuenta admin se le asigna un folio y se activa el status
        g.Folio = 2013
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
            if err := mailAvisoActivacionTpl.Execute(&hbody, g); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            msg := &mail.Message{
                    Sender:  MailSender,
                    To:      []string{g.Email},
                    Subject: MailSubjectConfirmed,
                    HTMLBody:	hbody.String(),
            }
            if err := mail.Send(c, msg); err != nil {
                // Problemas para enviar el correo NOK 
                http.Error(w, err.Error(), http.StatusInternalServerError)
                http.Redirect(w, r, "/", http.StatusFound)
            }
            // avisa del éxito independientemente del correo
            if err := mailActivationResultTpl.ExecuteTemplate(w, "confirm", g); err != nil {
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
            if err := mailActivationResultTpl.ExecuteTemplate(w, "confirm", g); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            defer r1.Body.Close()
        }
	} else {
		if err := mailActivationResultTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

