package load

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"model"
	"strings"
)

func init() {
    //http.HandleFunc("/r/normalize-mail", nmail)
}

func nmail(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Cta").Order("FechaHora")
	regdata := make([]model.Cta,0,500)

    if _, err := q.GetAll(c, &regdata); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	for _, cta := range regdata {
		emailL := strings.ToUpper(cta.Email)
		cta.Email = emailL
		_, err := datastore.Put(c, cta.Key(c), &cta)

		if err != nil {
	}

		fmt.Fprintf(w, "Email: %s - %s\n", cta.Email, emailL)
	}
}
