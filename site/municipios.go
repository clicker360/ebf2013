package site

import (
    "appengine"
	"appengine/memcache"
	"html/template"
	"encoding/json"
    "net/http"
	"fmt"
	"model"
	"sess"
)

func init() {
    http.HandleFunc("/r/msp", municipios)
}

func municipios(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	_, ok := sess.IsSess(w, r, c)
	if ok {
		if item, err := memcache.Get(c, "msp_"+r.FormValue("CveEnt")); err == memcache.ErrCacheMiss {
			if entidad, err := model.GetEntidad(c, r.FormValue("CveEnt")); err == nil {
				if municipios, _ := entidad.GetMunicipios(c); err == nil {
					b, _ := json.Marshal(municipios)
					item := &memcache.Item{
						Key:   "msp_"+r.FormValue("CveEnt"),
						Value: b,
					}
					if err := memcache.Add(c, item); err == memcache.ErrNotStored {
						c.Infof("item with key %q already exists", item.Key)
					} else if err != nil {
						c.Errorf("Error memcache.Add Municipio : %v", err)
					}
					c.Infof("CveMun generado: %v", item.Key)
					htmlComboMuns(w, municipios, r.FormValue("CveMun"))
				}
			}
		} else {
			//c.Infof("Memcache activo: %v", item.Key)
			var municipios []model.Municipio
			if err := json.Unmarshal(item.Value, &municipios); err != nil {
				c.Errorf("error adding item: %v", err)
			}
			htmlComboMuns(w, &municipios, r.FormValue("CveMun"))
		}
	}
	return
}

func htmlComboMuns(w http.ResponseWriter, municipios *[]model.Municipio, cvemun string) {
	tpl, _ := template.New("Mun").Parse(OptionTpl)
	fmt.Fprintf(w, `<select name="DirMun" id="MunSelector" onchange="locateAddress();">`)
	for _, m := range *municipios {
		if (m.CveMun == cvemun) {
			m.Selected = "selected"
		}
		tpl.Execute(w, m)
	}
	fmt.Fprintf(w, `</select>`)
}

const OptionTpl = `<option value="{{.CveMun}}" {{if .Selected}}selected="{{.Selected}}"{{end}}>{{.Municipio}}</option>
`
