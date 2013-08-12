package oferta

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"strings"
	"model"
	"time"
	"fmt"
)

const MailServer = "gmail"
func init() {
}

func putSearchData(c appengine.Context, value string, key *datastore.Key, idoft string, idcat int, enlinea bool) {
	r := strings.NewReplacer("."," ",","," ",";"," ",":"," ","!"," ","~"," ","¿"," ","?"," ","#"," ","_"," ","-"," ","+"," ","'"," ","\""," ","*"," ","$"," ","("," ",")"," ","="," ","%"," ","&"," ","<"," ",">"," ","|"," ","@"," ","·"," ","["," ","]"," ","{"," ","}"," ","¡"," ","!"," ", "\n", " ", "\r", " ", "\t", " ")
	blacklist := []string{
		"imbec",
		"imbéc",
		"pende",
		"pinch",
		"perron",
		"perrón",
		"mamona",
		"mamad",
		"mamo ",
		"mamon",
		"mamón",
		"mierd",
		"mrda",
		"puta",
		"puto",
		"puti",
		"pute",
		"verg",
		"pene",
		"concha",
		"maam",
		"caca",
		"caga",
		"cago",
		"boludo",
		"ching",
		"jodi",
		"cogon",
		"chiga",
		"fecal",
		"pucha",
		"coger",
		"cogido",
		"mame",
		"cule",
		"culo",
		"cula",
		"narco",
		"coño",
		"maric",
		"mamd",
		"mamar",
		"caga",
		"weba",
		"webo",
		"hueva",
		"huevos",
		"guevo",
		"güevo",
		"droga",
		"estupi",
		"estúpi",
		"fuck",
		"f*ck",
		"zeta",
		"joto",
		"jota",
		"marih",
		"coca",
		"crac",
		"nomam",
		"peda ",
		"pedo ",
		"nalg",
		"teta",
		"teto",
		"puch",
		"cabro",
		"chaqueta ",
		"puñeta ",
		"fecal",
		"f ecal",
		"fe cal",
		"fec al",
		"feca l",
		"f.e.c.a.l",
		"f cal",
		"felipe cal",
		"p nieto",
		"peña nieto",
		"pena nieto",
		"secuestr",
	}
	if err := model.DelOfertaSearchData(c, key); err != nil {
		c.Errorf("Datastore Delete Kind:SearchData, key:%s", key)
	} else {
		/*
		 * blacklist lookahead loop
		 */
		for _, v := range strings.Split(r.Replace(strings.ToLower(value)), " ") {
			i := 0
			for _, vv := range blacklist {
				if strings.HasPrefix(v, vv) {
					i = i+1
				}
			}
			if i > 0 {
				c.Errorf("BLACKLIST PutSearchData, IdOft:%s, word:%s", idoft, value)
				if (MailServer == "gmail") {
					msg := &mail.Message{
						Sender:  "contacto@elbuenfin.org",
						To:      []string{"ahuezo@clicker360.com", "daniela@iniciativamexico.com", "adan@clicker360.com"},
						Subject: "Alerta de actividad anormal de blacklist / El Buen Fin",
						HTMLBody: fmt.Sprintf("Se evitó integrar a búsquedas la siguiente Oferta ID: %v, \n Con el siguiente Texto: %v \n Favor de revisar el resto de la actividad relacionada a la oferta.", idoft, value),
					}
					//To:      []string{"ahuezo@clicker360.com", "thomas@clicker360.com", "adan@clicker360.com"},
					if err := mail.Send(c, msg); err != nil {
						/* Problemas para enviar el correo NOK */
						c.Errorf("BLACKLIST MAIL NOT SEND, IdOft:%s, word:%s", idoft, value)
						panic(err)
					}
				} else {
					c.Errorf("NO MAILSERVER BLACKLIST, IdOft:%s, word:%s", idoft, value)
					panic(err)
				}
				//return
			}
		}
		/*
		 * All clear
		 */
		for _, v := range strings.Split(r.Replace(value), " ") {
			if(len(v)>3) {
				w := strings.ToLower(v)
				if(model.ValidSearchData.MatchString(w)) {
					sd := &model.SearchData {
						Sid: key.Encode(),
						Kind: "Oferta",
						Field: "Descripcion",
						Value: w,
						IdCat: idcat,
						Enlinea: enlinea,
						FechaHora: time.Now(),
					}
					// Pa llave es el idoft + palabra clave
					_, err := datastore.Put(c, datastore.NewKey(c, "SearchData", idoft+w, 0, nil), sd)
					if err != nil {
						c.Errorf("Datastore Put Kind:SearchData, IdOft:%s, word:%s", idoft, w)
					}
				} else {
					c.Infof("Intento de palabra inválida al diccionario, IdOft:%s, word:%s", idoft, w)
				}
			}
		}
	}
}

