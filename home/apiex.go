package home

import (
	"fmt"
	"net/http"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"os"
)

func init() {
    http.HandleFunc("/faqex", FAQEx)
}

func FAQEx(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	response, err := client.Get("http://movil.ebfmex-pub.appspot.com/wsfaq")
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        fmt.Fprintf(w,"%s\n", string(contents))
    }
}
