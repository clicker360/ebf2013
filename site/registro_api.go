package site

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"
	"appengine/user"
	"bytes"
    "strings"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"model"
	"net/http"
	"net/url"
	"sess"
	"time"
)

type UrlCfm struct {
	Md5			string		`json:"md5", omitempty`
	Nombre		string		`json:"nombre, omitempty"`
	Apellidos	string		`json:"apellidos, omitempty"`
	Email		string		`json:"email, omitempty"`
	FechaHora	time.Time	`json:"fechahora, omitempty"`
	Llave		string		`json:"llave. omitempty"`
	AppId		string		`json:"appid, omitempty"`
    Pass        string      `json:"pass, omitempty"`
	Status		string		`json:status,omitempty`
}

type WsCta struct {
	Folio		int32		`json:"Folio,omitempty"`
	Nombre		string		`json:"Nombre"`
	Apellidos	string		`json:"Apellidos"`
	Puesto		string		`json:"Puesto"`
	Email		string		`json:"Email"`
	EmailAlt	string		`json:"EmailAlt"`
	Pass		string		`json:"Pass,omitempty"`
	Tel			string		`json:"Tel"`
	Cel			string		`json:"Cel,omitempty"`
	FechaHora	time.Time	`json:"FechaHora"`
	UsuarioInt	string		`json:"UsuarioInt,omitempty"`
	CodigoCfm	string		`json:"-"`
	Pass2		string		`json:"Pass2,omitempty"`
	TermCond	string		`json:"TermCond,omitempty"`
	Status	    bool		`json:"-"`
	StatusMsg	string		`json:"status,omitempty"`
	Ackn	    string		`json:"ackn,omitempty"`
	Errors		map[string]bool `json:"errors,omitempty"`
}

const WsMailServer = "gmail"
const WsMailSender = "El Buen Fin <contacto@elbuenfin.org>"
const WsMailSubject = "Codigo de Activación de Registro / El Buen Fin en línea"
const WsMailSubjectConfirmed = "Cuenta Activada / El Buen Fin en línea"

var wsMailActivationCodeTpl = template.Must(template.ParseFiles("layout/mail_activation_code.html"))
var wsMailAvisoActivacionTpl = template.Must(template.ParseFiles("layout/mail_activation_welcome.html"))

var Loc *time.Location

func init() {
    Loc, _ = time.LoadLocation("America/Mexico_City")
	http.HandleFunc("/r/wsr/put", PutCta)
	http.HandleFunc("/r/wsr/post", PostCta)
	http.HandleFunc("/r/wsr/get", GetCta)
	http.HandleFunc("/r/wsr/del", DelCta)
	http.HandleFunc("/r/wsr/confirm", ConfirmCode)
}

func PutCta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsCta
	defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); ok {
		out.StatusMsg = "alreadyOnSession"
		return
	}
	if r.Method != "POST" {
		out.StatusMsg = "wrongMethod"
		return
	}

	// Se obtienen y validan los campos del cgi
	wsCtaTmp := regFill(r)
	if out.Errors, out.StatusMsg = regValidate(wsCtaTmp, r.FormValue("t")); out.StatusMsg != "ok" {
		return
	}
    // En el caso particular de crear una nueva cta se verifica que ambos campos
    // de password no sean vacíos. Para el Post no es necesario pues se asume que si
    // ambos son vacíos entonces no se solicita un cambio de pass
    if r.FormValue("Pass") == "" && r.FormValue("Pass2") == "" {
        err := make(map[string]bool)
        err["Pass"] = false
        err["Pass2"] = false
        out.Errors = err
        out.StatusMsg = "invalidInput"
        return
    }

    // En el caso particular de crear cuenta se verifica que acepten términos y 
    // condiciones
	if r.FormValue("TermCond") != "1" {
        err := make(map[string]bool)
		err["TermCond"] = false
        out.Errors = err
        out.StatusMsg = "invalidInput"
        return
	}

	cta, err := model.GetCta(c, wsCtaTmp.Email)
	if err != nil {
        cta.Status = false
    }
	// No hay Cuenta registrada
    cta.FechaHora = time.Now().In(Loc)
    cta.Nombre = wsCtaTmp.Nombre
    cta.Apellidos = wsCtaTmp.Apellidos
    cta.Puesto = wsCtaTmp.Puesto
    cta.EmailAlt = wsCtaTmp.EmailAlt
    cta.Tel = wsCtaTmp.Tel
    cta.Cel = wsCtaTmp.Cel
    cta.Pass = wsCtaTmp.Pass
    cta.Folio = wsCtaTmp.Folio
    cta.UsuarioInt = ""

    // Generar código de confirmación distindo cada vez. Md5 del email + fecha-hora
    h := md5.New()
    io.WriteString(h, fmt.Sprintf("%s%s%s%s", time.Now().In(Loc), cta.Email, cta.Pass, model.RandId(12)))
    cta.CodigoCfm = fmt.Sprintf("%x", h.Sum(nil))
	//Si hay estatus es que ya existe
	if cta.Status == false {

		// Se agrega la cuenta sin activar para realizar el proceso de código de confirmación
		if err := model.PutCta(c, cta); err != nil {
			out.StatusMsg = "writeError"
			return
		}

		setWsCta(&out, *cta)
		/* No se ha activado, por tanto se inicia el proceso de código de verificación */
		m := UrlCfm {
			Md5:       cta.CodigoCfm,
			Nombre:    cta.Nombre,
			Email:     cta.Email,
			FechaHora: time.Now().In(Loc),
			Llave:     cta.Key(c).Encode(),
			AppId:     appengine.AppID(c),
		}
		var hbody bytes.Buffer
		// Envia código activación
		if err := wsMailActivationCodeTpl.Execute(&hbody, m); err != nil {
			out.StatusMsg = err.Error()
			return
		}
		if WsMailServer == "gmail" {
			msg := &mail.Message{
				Sender:   WsMailSender,
				To:       []string{m.Email},
				Subject:  WsMailSubject,
				HTMLBody: hbody.String(),
			}
			if err := mail.Send(c, msg); err != nil {
				/* Problemas para enviar el correo NOK */
				out.StatusMsg = err.Error()
				return
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
				out.StatusMsg = err.Error()
				return
			}
			if r1.StatusCode != 200 {
				out.StatusMsg = "errorTransporMail"
				return
			}
			defer r1.Body.Close()
		}
		/*
		 Si hay usuario admin se despliega el código de activación
		*/
		if gu := user.Current(c); gu != nil {
			out.Ackn = cta.CodigoCfm
		}
		out.StatusMsg = "ok"
	} else {
		out.StatusMsg = "alreadyRegistered"
	}
}

func ConfirmCode(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out UrlCfm
	defer model.JsonDispatch(w, &out)
	if _, ok := sess.IsSess(w, r, c); ok {
		out.Status = "alreadyRegistered"
		return
	}
	if r.Method != "GET" {
		out.Status = "wrongMethod"
		return
	}
	md5 := r.FormValue("m")
	key, _ := datastore.DecodeKey(r.FormValue("c"))
	var cta model.Cta
	if err := datastore.Get(c, key, &cta); err != nil {
		out.Status = "notFound"
		return
	}

	/* Se verifica el código de confirmación */
	if cta.CodigoCfm == md5 && cta.Status == false {
		// Si se confirma el md5 la cuenta admin se le asigna un folio y se activa el status
		cta.Folio = 2013 // EL folio ya no se usa, se utilizará por el momento como referencia
		cta.Status = true
		cta.CodigoCfm = "Confirmado"
		if _, err := datastore.Put(c, key, &cta); err != nil {
			out.Status = err.Error()
			return
		}
        // Se llena la estructura de salida json
        out.Md5 = md5
        out.Nombre = cta.Nombre
        out.Apellidos = cta.Apellidos
        out.Email = cta.Email
        out.Pass = cta.Pass
        out.FechaHora = cta.FechaHora
		if gu := user.Current(c); gu != nil {
            out.Llave = cta.Key(c).Encode()
            out.AppId = appengine.AppID(c)
        }
		/* Prende la sesion */
		if _, _, err := sess.SetSess(w, c, key, cta.Email, cta.Nombre); err != nil {
			out.Status = err.Error()
			return
		}

		if WsMailServer == "gmail" {
			// Envia código activación
			var hbody bytes.Buffer
			if err := wsMailAvisoActivacionTpl.Execute(&hbody, cta); err != nil {
				out.Status = err.Error()
			}
			msg := &mail.Message{
				Sender:   WsMailSender,
				To:       []string{cta.Email},
				Subject:  WsMailSubjectConfirmed,
				HTMLBody: hbody.String(),
			}
			if err := mail.Send(c, msg); err != nil {
				// Problemas para enviar el correo NOK
				out.Status = err.Error()
				return
			}
			// avisa del éxito independientemente del correo
			out.Status = "ok"
			return
		} else {
			client := urlfetch.Client(c)
			url := fmt.Sprintf("http://envia-m.mekate.com.mx/?Sender=%s&Tipo=Aviso&Email=%s&Nombre=%s&Pass=%s&AppId=ebfmxorg",
				"registro@elbuenfin.org",
				cta.Email,
				url.QueryEscape(cta.Nombre),
				url.QueryEscape(cta.Pass))
			r1, err := client.Get(url)
			if err != nil {
				out.Status = err.Error()
				return
			}

			if r1.StatusCode != 200 {
				out.Status = "errorTransporMail"
				return
			}
			out.Status = "ok"
			defer r1.Body.Close()
		}
	} else {
		out.Status = "codeError"
	}
}

func GetCta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsCta
    defer model.JsonDispatch(w, &out)
	s, ok := sess.IsSess(w, r, c)
    if !ok {
		out.StatusMsg = "noSession"
        return
    }
    if r.Method != "GET" {
		out.StatusMsg = "wrongMethod"
        return
    }

	if cta, err := model.GetCta(c, s.User); err != nil {
		out.StatusMsg = "notFound"
	} else {
		setWsCta(&out, *cta)
		out.StatusMsg = "ok"
        // no queremos pasar el pass
        out.Pass = ""
        out.Pass2 = ""
	}
}

func PostCta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsCta
    defer model.JsonDispatch(w, &out)
	s, ok := sess.IsSess(w, r, c)
    if !ok {
		out.StatusMsg = "noSession"
        return
    }
    if r.Method != "POST" {
		out.StatusMsg = "wrongMethod"
        return
    }
	wsCtaTmp := regFill(r)
    out.Errors, out.StatusMsg = regValidate(wsCtaTmp, r.FormValue("t"))
    if out.StatusMsg != "ok" {
        return
    }
	cta, err := model.GetCta(c, s.User)
    if err != nil {
		out.StatusMsg = "notFound"
        return
	}
    if !cta.Status {
        out.StatusMsg = "notConfirmedUser"
        return
    }
    cta.Nombre = wsCtaTmp.Nombre
    cta.Apellidos = wsCtaTmp.Apellidos
    cta.EmailAlt = wsCtaTmp.EmailAlt
    cta.Puesto = wsCtaTmp.Puesto
    cta.Tel = wsCtaTmp.Tel
    cta.Cel = wsCtaTmp.Cel
    if wsCtaTmp.Pass != "" {
        cta.Pass = wsCtaTmp.Pass
    }
    cta.Folio = wsCtaTmp.Folio
    cta.UsuarioInt = ""
    setWsCta(&out, *cta)
    if err := model.PutCta(c, cta); err != nil {
        out.StatusMsg = "writeError"
    } else {
        out.StatusMsg = "ok"
    }
}

func DelCta(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    var out WsCta
    defer model.JsonDispatch(w, &out)
	s, ok := sess.IsSess(w, r, c)
    if !ok {
		out.StatusMsg = "noSession"
        return
    }
    if r.Method != "POST" {
		out.StatusMsg = "wrongMethod"
        return
    }
	now := time.Now().In(Loc)
	if cta, err := model.GetCta(c, s.User); err != nil {
		out.StatusMsg = "notFound"
	} else {
		setWsCta(&out, *cta)
		// Sólo desactiva cuenta si no hay empresas dependientes
		empresas := listEmp(c, cta)
		if len(*empresas) != 0 {
			// Debe borrar empresas antes o Transferir sus empresas a otro usuario
			out.StatusMsg = "accountNotEmpty"
			return
		}
		// Desactiva Status
		if(r.FormValue("desactiva")=="1") {
			s.Expiration = now.AddDate(-1,0,0)
			if _, err := datastore.Put(c, datastore.NewKey(c, "Sess", s.User, 0, nil), &s); err != nil {
				out.StatusMsg = err.Error()
				return
			}
			cta.CodigoCfm = "Desactivado"
			cta.Status = false
			if err = model.PutCta(c, cta); err != nil {
				out.StatusMsg = "writeError"
				return
			}
			w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))
			w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))

			// INICIA ENVIO DE CORREO DE MOTIVOS
			// Este tramo no debe arrojar errores al usuario
			var hbody bytes.Buffer
			cta.CodigoCfm = r.FormValue("motivo")
			cancelMessageTpl.Execute(&hbody, cta)
			msg := &mail.Message{
				Sender:  "Cancelación de cuenta / Buen Fin <contacto@elbuenfin.org>",
				To:      []string{"contacto@elbuenfin.org"},
				Subject: "Aviso de motivo de cuenta cancelada / El Buen Fin en línea",
				HTMLBody: hbody.String(),
			}
			mail.Send(c, msg)
			out.Ackn = "¡Gracias por participar en El Buen Fin!"
			out.StatusMsg = "ok"
		} else {
            out.StatusMsg = "notGranted"
        }
	}
}

func setWsCta(out *WsCta, e model.Cta) {
	out.Nombre =	e.Nombre
	out.Apellidos =	e.Apellidos
	out.Puesto =	e.Puesto
	out.Email =		e.Email
	out.EmailAlt =	e.EmailAlt
	out.FechaHora =	e.FechaHora
	out.Pass =		e.Pass
	out.Tel =		e.Tel
	out.Cel =		e.Cel
}

func regFill(r *http.Request) WsCta {
	e := WsCta {
        Folio:      2013,
        Nombre:	    strings.TrimSpace(r.FormValue("Nombre")),
        Apellidos:  strings.TrimSpace(r.FormValue("Apellidos")),
        Puesto:	    strings.TrimSpace(r.FormValue("Puesto")),
        Email:	    strings.TrimSpace(r.FormValue("Email")),
        EmailAlt:	strings.TrimSpace(r.FormValue("EmailAlt")),
        Tel:		strings.TrimSpace(r.FormValue("Tel")),
        Cel:		strings.TrimSpace(r.FormValue("Cel")),
        TermCond: strings.TrimSpace(r.FormValue("TermCond")),
    }
    if r.FormValue("Pass") != "" {
        e.Pass = strings.TrimSpace(r.FormValue("Pass"))
    }
    if r.FormValue("Pass") != "" {
        e.Pass2 = strings.TrimSpace(r.FormValue("Pass2"))
    }
    return e
}

func regValidate(e WsCta, tipo string) (map[string]bool, string) {
    errmsg := "ok"
    err := make(map[string]bool)
    if e.Nombre == "" || !model.ValidName.MatchString(e.Nombre) {
		 err["Nombre"] = false
    }
	if e.Apellidos == "" || !model.ValidName.MatchString(e.Apellidos) {
		err["Apellidos"] = false
	}
	if e.Puesto != "" && !model.ValidSimpleText.MatchString(e.Puesto) {
		err["Puesto"] = false
	}
	if e.Email == "" || !model.ValidEmail.MatchString(e.Email) {
		err["Email"] = false
	}
	if e.EmailAlt != "" && !model.ValidEmail.MatchString(e.EmailAlt) {
		err["EmailAlt"] = false
	}
    if tipo != "upd" {
		if (e.Pass != e.Pass2 || e.Pass == "" || e.Pass2 == "" || !model.ValidPass.MatchString(e.Pass)) {
            err["Pass"] = false
            err["Pass2"] = false
        }
    } else {
		if ((e.Pass != e.Pass2 || !model.ValidPass.MatchString(e.Pass)) && (e.Pass != "" || e.Pass2 != "")) {
            err["Pass"] = false
            err["Pass2"] = false
        }
    }
	if e.Tel == "" || !model.ValidTel.MatchString(e.Tel) {
		err["Tel"] = false
	}
	if e.Cel != "" && !model.ValidTel.MatchString(e.Cel) {
		err["Cel"] = false
	}
	if e.TermCond != "1" {
		err["TermCond"] = false
    }
    for _, v := range err {
        if v == false {
            errmsg = "invalidInput"
        }
    }
	return err, errmsg
}
