package load

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "net/http"
	"model"
	"fmt"
)

func init() {
	http.HandleFunc("/r/loado", loadOrg)
	http.HandleFunc("/r/loadcat", loadCat)
}

func loadOrg(w http.ResponseWriter, r *http.Request) {
//	return
	c := appengine.NewContext(r)
    if u := user.Current(c); u != nil {
		o := []model.Organismo {
			{"ABM", "Asociación de Bancos de México", ""},
			{"AMIB","Asociación Mexicana de Intermediarios Bursátiles", ""},
			{"AMIS", "Asociación Mexicana de Instituciones de Seguros", ""},
			{"AMPICI", "Asociación Mexicana de Internet", ""},
			{"ANTAD", "Asociación Nacional de Tiendas de Autoservicio y Departamentales", ""},
			{"CANACINTRA", "Cámara Nacional de la Industria de Transformación", ""},
			{"CNA", "Consejo Mexicano de Hombres de Negocios", ""},
			{"COMCE", "Consejo Empresarial Mexicano de Comercio Exterior, Inversión y Tecnología", ""},
			{"CANACOPE", "Cámara de Comercio, Servicios y Turismo en Pequeño de la Ciudad de México", ""},
			{"CONCANACO/CONCAMIN", "Confederación de Cámaras Nacionales de Comercio, Servicio y Turismo", ""},
			{"COPARMEX", "Confederación Patronal de la República Mexicana", ""},
			{"OTRO", "Otro Organismo", ""},
		}
		for _, e := range o {
			fmt.Fprintf(w, "Organismo: %d, %d", e.Siglas, e.Nombre)
			_, err := datastore.Put(c, datastore.NewKey(c, "Organismo", e.Nombre, 0, nil), &e)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
        url, _ := user.LoginURL(c, "/")
        fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
        return
    }

}

func loadCat(w http.ResponseWriter, r *http.Request) {
//	return
	c := appengine.NewContext(r)
    if u := user.Current(c); u != nil {
		o := []model.Categoria {
			{1, "Abarrotes", ""},
			{2, "Automóviles y transportes", ""},
			{3, "Bancos y seguros", ""},
			{4, "Blancos", ""},
			{5, "Calzado", ""},
			{6, "Cuidado e higiene personal", ""},
			{7, "Deportes", ""},
			{8, "Electrónica y video", ""},
			{9, "Electrodomésticos", ""},
			{10, "Farmacias", ""},
			{11, "Ferretería", ""},
			{12, "Fotografía y cómputo", ""},
			{13, "Hogar y decoración", ""},
			{14, "Jardinería", ""},
			{15, "Joyería", ""},
			{16, "Jugetería", ""},
			{17, "Línea blanca", ""},
			{18, "Muebles", ""},
			{19, "Oficina y papelería", ""},
			{20, "Oportunidades", ""},
			{21, "Óptica", ""},
			{22, "Perecederos", ""},
			{23, "Regalos", ""},
			{24, "Restaurantes", ""},
			{25, "Ropa", ""},
			{26, "Servicios", ""},
			{27, "Varios", ""},
			{28, "Viajes", ""},
		}
		for _, e := range o {
			fmt.Fprintf(w, "Categoria: %d, %d", e.IdCat, e.Categoria)
			_, err := datastore.Put(c, datastore.NewKey(c, "Categoria", "", int64(e.IdCat), nil), &e)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		fmt.Fprintf(w, "Usuario inválido")
    }
}
