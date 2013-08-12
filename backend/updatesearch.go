package backend

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"strconv"
	"fmt"
)

func init() {
    //http.HandleFunc("/backend/updatesearch", fetchUpdateSearch)
    //http.HandleFunc("/backend/updatesearch", RedirUpdateSearch)
}

func fetchUpdateSearch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	/* 
	 * En lo que encontramos una manera digna de ejecutar el cron, se
	 * podría meter una llave en el memcaché y crearla en el administrador
	 * O bien en el datastore.
	 */
	var m int
	m, _ = strconv.Atoi(r.FormValue("m"))
	if m < 30 {
		m = 30
	}
	if r.FormValue("c") == "ZWJmbWV4LXB1YnIeCxISX0FoQWRtaW5Yc3JmVG9rZW5fIgZfWFNSRl8M" {
		url := fmt.Sprintf( "http://movil.%s.appspot.com/backend/updatesearch?minutes=%d&token=%v", appengine.AppID(c), m, r.FormValue("c"))
		ret, err := client.Get(url)
		if err != nil {
			c.Errorf("updatesearch in %v, %v %v, %v", appengine.AppID(c), url, ret.Status, err)
		} else {
			c.Infof("updatesearch in %v, %v %v, minutes=%v", appengine.AppID(c), url, ret.Status, m)
		}
	}
	return
}

func RedirUpdateSearch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	/* 
	 * En lo que encontramos una manera digna de ejecutar el cron, se
	 * podría meter una llave en el memcaché y crearla en el administrador
	 * O bien en el datastore.
	 */
	var m int
	m, _ = strconv.Atoi(r.FormValue("m"))
	if m < 30 {
		m = 30
	}
	if r.FormValue("c") == "ZWJmbWV4LXB1YnIeCxISX0FoQWRtaW5Yc3JmVG9rZW5fIgZfWFNSRl8M" {
		url := fmt.Sprintf( "http://movil.%s.appspot.com/backend/updatesearch?minutes=%d&token=%v", appengine.AppID(c), m, r.FormValue("c"))
		c.Infof("updatesearch in %v, %v, minutes=%v", appengine.AppID(c), url, m)
		http.Redirect(w, r, url, http.StatusFound)
	}
	return
}
