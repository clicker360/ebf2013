package home

import (
    "appengine"
    "appengine/datastore"
    "appengine/memcache"
	"html/template"
	"encoding/json"
	"math/rand"
    "net/http"
	"sortutil"
	"strings"
	"strconv"
    "model"
	"time"
)

type response struct {
	IdEmp	string `json:"id"`
	Name	string `json:"name"`
	Url		string `json:"url"`
	Num		int `json:"num"`
}

type Paginador struct {
	Prefix string
	Pagina int
}

func init() {
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/dirtexto", directorioTexto)
    http.HandleFunc("/wsdiremp", wsDirTexto)
    http.HandleFunc("/carr", carr)
    rand.Seed(time.Now().UnixNano())
}

/*
 * La idea es hacer 60 cachés al azar con un tiempo de vida de 30 min
 * Cada que se muere un memcache se genera otro carrousel al azar de logos
 */
func carr(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
    c := appengine.NewContext(r)
	var timetolive = 7200 //seconds
	var b []byte
	var nn int = 50 // tamaño del carrousel
	logos := make([]model.Image, 0, nn)
	hit := rand.Intn(60)
	if item, err := memcache.Get(c, "carr_"+strconv.Itoa(hit)); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("EmpLogo")
		n, _ := q.Count(c)
		offset := 0;
		if(n > nn) {
			offset = rand.Intn(n-nn)
		} else {
			nn = n
		}
		q = q.Offset(offset).Limit(nn)
		if _, err := q.GetAll(c, &logos); err != nil {
			if err == datastore.ErrNoSuchEntity {
				return
			}
		}

		b, _ = json.Marshal(logos)
		item := &memcache.Item{
			Key:   "carr_"+strconv.Itoa(hit),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Memcache.Add carr_idoft : %v", err)
		}
		//c.Infof("memcache add carr_page : %v", strconv.Itoa(hit))
	} else {
		//c.Infof("memcache retrieve carr_page : %v", strconv.Itoa(hit))
		if err := json.Unmarshal(item.Value, &logos); err != nil {
			c.Errorf("Unmarshaling EmpLogo item: %v", err)
		}
		nn = len(logos)
	}

	tpl, _ := template.New("Carr").Parse(cajaTpl)
	tn := rand.Perm(nn)
	var ti response
	for i, _ := range tn {
		ti.IdEmp = logos[tn[i]].IdEmp
		ti.Name = logos[tn[i]].Name
		ti.Url = strings.Replace(logos[tn[i]].Sp4, "s180", "s70",1)
		if ti.Url != "" {
			//b, _ := json.Marshal(ti)
			//w.Write(b)
			tpl.Execute(w, ti)
		}
		if i >= nn  {
			break
		}
	}
}

func directorioTexto(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/*
		Loop para recorrer todas las empresas 
	*/
	now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
	prefixu := strings.ToLower(r.FormValue("prefix"))
	ultimos := r.FormValue("ultimos")
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	page -= 1
	const batch = 200
    q := datastore.NewQuery("EmpresaNm")
	var timetolive = 14400 //seconds
	if ultimos != "1" {
		q = q.Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd").Order("Nombre")
		/*
		 * Se pagina ordenado alfabéticamente el resutlado de la búsqueda 
		 * y se guarda en Memcache
		 */
		lot, _ := q.Count(c)
		pages := lot/batch
		if lot%batch > 0 {
			pages += 1
		}
		if pages > 1 {
			Paginas := make([]Paginador, pages)
			//c.Infof("lote: %d, paginas : %d", lot, pages)
			for i := 0; i < pages; i++ {
				Paginas[i].Prefix = prefixu
				Paginas[i].Pagina = i+1
				//c.Infof("pagina : %d", i)
			}
			tplp, _ := template.New("paginador").Parse(paginadorTpl)
			tplp.Execute(w, Paginas)
		}

		var empresas []model.EmpresaNm
		if item, err := memcache.Get(c, "dirprefix_"+prefixu+"_"+strconv.Itoa(page)); err == memcache.ErrCacheMiss {
			//c.Infof("memcached prefix: %v, pagina : %d", prefixu, page)
			offset := batch * page
			q = q.Offset(offset).Limit(batch)
			if _, err := q.GetAll(c, &empresas); err != nil {
				return
			}
			b, _ := json.Marshal(empresas)
			item := &memcache.Item {
				Key:   "dirprefix_"+prefixu+"_"+strconv.Itoa(page),
				Value: b,
				Expiration: time.Duration(timetolive)*time.Second,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("memcache.Add dirprefix : %v", err)
			}
		} else {
			if err := json.Unmarshal(item.Value, &empresas); err != nil {
				c.Errorf("Memcache Unmarshalling dirprefix item: %v", err)
			}
		}

		sortutil.CiAscByField(empresas, "Nombre")
		var ti response
		var tictac int
		var repetido string
		tictac = 1
		for k, _ := range empresas {
			tpl, _ := template.New("pagina").Parse(empresaTpl)
			ti.Num = tictac
			ti.IdEmp = empresas[k].IdEmp
			ti.Name = strings.Title(empresas[k].Nombre)
			if repetido != ti.Name {
				if tictac != 1 {
					tictac = 1
				} else {
					tictac = 2
				}
				repetido = ti.Name
				tpl.Execute(w, ti)
			}
		}
	} else {
		prefixu = ""
		var empresas []model.EmpresaNm
		if item, err := memcache.Get(c, "dirprefix_"+prefixu+"_"+strconv.Itoa(page)); err == memcache.ErrCacheMiss {
			//c.Infof("memcached prefix: %v, pagina : %d", prefixu, page)
			q = datastore.NewQuery("Empresa").Filter("FechaHora >=", now.AddDate(0,0,-2)).Limit(400)
			var empresas []model.Empresa
			if _, err := q.GetAll(c, &empresas); err != nil {
				return
			}
			b, _ := json.Marshal(empresas)
			item := &memcache.Item {
				Key:   "dirprefix_"+prefixu+"_"+strconv.Itoa(page),
				Value: b,
				Expiration: time.Duration(timetolive)*time.Second,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("memcache.Add dirprefix : %v", err)
			}
		} else {
			if err := json.Unmarshal(item.Value, &empresas); err != nil {
				c.Errorf("Memcache Unmarshalling dirprefix item: %v", err)
			}
		}

		sortutil.CiAscByField(empresas, "Nombre")
		var ti response
		var tictac int
		var repetido string
		tictac = 1
		for k, _ := range empresas {
			tpl, _ := template.New("pagina").Parse(empresaTpl)
			ti.Num = tictac
			ti.IdEmp = empresas[k].IdEmp
			ti.Name = strings.Title(strings.ToLower(empresas[k].Nombre))
			if repetido != ti.Name {
				if tictac != 1 {
					tictac = 1
				} else {
					tictac = 2
				}
				repetido = ti.Name
				tpl.Execute(w, ti)
			}
		}
	}
}

//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('/spic?IdEmp={{.IdEmp}}')"></div></div>`
const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('{{.Url}}')"></div></div>`
const empresaTpl = `<div class="gridsubRow bg-Gry{{.Num}}"><a href="http://www.elbuenfin.org/micrositio.html?id={{.IdEmp}}" target="_blank">{{.Name}}</a></div>`
const paginadorTpl = `<div class="pagination-H"><ul id="letters">{{range .}}<li><a href="#" class="letter" prfx="{{.Prefix}}" onclick="javascript:paginar({{.Pagina}});"> {{.Pagina}} </a></li>{{end}}</ul></div>`
//const paginadorTpl = `<div>{{range .}}<a href="javascript:pager({{.Prefix}}, {{.Pagina}});"> {{.Pagina}} </a>{{end}}</div>`
//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><img class="centerimg" src="/spic?IdEmp={{.IdEmp}}" /></div>`

type WsEmpresa struct{
	Id		string `json:"id"`
	Empresa	string `json:"empresa"`
	Url		string `json:"url"`
}

func wsDirTexto(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	prefixu := strings.ToUpper(r.FormValue("prefix"))
    q := datastore.NewQuery("Empresa").Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd") //.Filter("Status =", true)
	em, _ := q.Count(c)
	empresas := make([]model.Empresa, em, em)
	if _, err := q.GetAll(c, empresas); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}

	var b []byte
	wsout := make([]WsEmpresa, em, em)
	sortutil.CiAscByField(empresas, "Nombre")
	for i, _ := range empresas {
		wsout[i].Id = empresas[i].IdEmp
		wsout[i].Empresa = empresas[i].Nombre
		wsout[i].Url = empresas[i].Url
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ = json.Marshal(wsout)
	w.Write(b)
}


