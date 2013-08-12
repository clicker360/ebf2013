package backend

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"fmt"
)

func init() {
    //http.HandleFunc("/backend/mvblob", MvBlob)
    http.HandleFunc("/mvblob", RedirMvBlob)
}

func MvBlob(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	url := fmt.Sprintf( "http://movil.%s.appspot.com/mvblob/generate", appengine.AppID(c))
	ret, err := client.Get(url)
	if err != nil {
		c.Errorf("mvblob in %v, %v %v, %v", appengine.AppID(c), url, ret.Status, err)
	} else {
		c.Infof("mvblob in %v, %v %v", appengine.AppID(c), url, ret.Status)
	}
	return
}

func RedirMvBlob(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url := fmt.Sprintf( "http://movil.%s.appspot.com/mvblob/generate", appengine.AppID(c))
	c.Infof("redirect to mvblob/generate in %v, %v", appengine.AppID(c), url)
	http.Redirect(w, r, url, http.StatusFound)
	return
}
