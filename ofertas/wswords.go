package oferta

import (
    "appengine"
	"encoding/json"
    "net/http"
	"strings"
	"model"
	"time"
)

type Word struct{
	Token	string `json:"token"`
	Id		string `json:"id"`
	Status  string `json:"status"`
}

func init() {
    http.HandleFunc("/r/addword", AddWord)
    http.HandleFunc("/r/delword", DelWord)
    http.HandleFunc("/r/rmword", RmWord)
    http.HandleFunc("/r/wordsxo", ShowOfWords)
    http.HandleFunc("/r/wordsxe", ShowEmpWords)
}

func AddWord(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Word
	out.Token = r.FormValue("token")
	out.Id = r.FormValue("id")
	if model.ValidSimpleText.MatchString(out.Token) {
		tokens := strings.Fields(r.FormValue("token"))
		oferta,_ := model.GetOferta(c, out.Id)
		if(len(tokens) > 0) {
			for _,v:= range tokens {
				if oferta.IdEmp != "none" {
					var palabra model.OfertaPalabra
					palabra.IdOft = out.Id
					palabra.IdEmp = oferta.IdEmp
					palabra.Palabra = strings.ToLower(v)
					palabra.FechaHora = time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
					err := oferta.PutOfertaPalabra(c, &palabra)
					if err != nil {
						out.Status = "writeErr"
					} else {
						out.Status = "ok"
					}
				} else {
					out.Status = "notFound"
				}
			}
		} else {
			out.Status = "invalidText"
		}
	} else {
		out.Status = "invalidText"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

/*
	Se remueve una palabra de la oferta sólamente
*/
func DelWord(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Word
	out.Token = strings.ToLower(r.FormValue("token"))
	out.Id = r.FormValue("id")
	oferta,_ := model.GetOferta(c, out.Id)
	if oferta.IdEmp != "none" {
		err := model.DelOfertaPalabra(c, out.Id, out.Token)
		if err != nil {
			out.Status = "writeErr"
		} else {
			out.Status = "ok"
		}
		out.Id = "none"
	} else {
		out.Status = "notFound"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

/*
	Se remueve una palabra de la empresa
	Recibe un Id de Empresa
*/
func RmWord(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Word
	out.Token = strings.ToLower(r.FormValue("token"))
	out.Id = r.FormValue("id")
	out.Status = "ok"
	if err := model.RmOfertaPalabra(c, out.Id, out.Token); err != nil {
		out.Status = "notFound"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

func ShowOfWords(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var b []byte
	ofertas := model.GetOfertaPalabras(c, r.FormValue("id"), "")
	if *ofertas != nil {
		if(len(*ofertas) > 0) {
			words := make([]Word, len(*ofertas) ,len(*ofertas))
			for i,v:= range *ofertas {
				words[i].Id = v.IdOft
				words[i].Token = v.Palabra
				words[i].Status = ""
			}
			b, _ = json.Marshal(words)
		} else {
			var out Word
			out.Id = r.FormValue("id")
			out.Token = ""
			out.Status = "notFound"
			b, _ = json.Marshal(out)
		}
	} else {
		var out Word
		out.Id = r.FormValue("id")
		out.Token = ""
		out.Status = "notFound"
		b, _ = json.Marshal(out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func ShowEmpWords(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var b []byte
	ofertas := model.GetOfertaPalabras(c, "", r.FormValue("id"))
	if *ofertas != nil {
		if(len(*ofertas) > 0) {
			words := make([]Word, len(*ofertas) ,len(*ofertas))
			for i,v:= range *ofertas {
				words[i].Id = v.IdOft
				words[i].Token = v.Palabra
				words[i].Status = ""
			}
			b, _ = json.Marshal(words)
		} else {
			var out Word
			out.Id = r.FormValue("id")
			out.Token = ""
			out.Status = "notFound"
			b, _ = json.Marshal(out)
		}
	} else {
		var out Word
		out.Id = r.FormValue("id")
		out.Token = ""
		out.Status = "notFound"
		b, _ = json.Marshal(out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
