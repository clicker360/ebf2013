package backend

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "net/http"
    "model"
	//"strings"
	"strconv"
    "fmt"
)


func init() {
    http.HandleFunc("/backend/registros.csv", registroCsv)

}

func registroCsv(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    if u := user.Current(c); u == nil {
		return
	}

	var lote = 200
	pagina,_ := strconv.Atoi(r.FormValue("pg"));

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-type", "application/octet-stream");
	w.Header().Set("Content-Disposition", "attachment; filename=\"reportecta.csv\"");
	w.Header().Set("Accept-Charset","utf-8");

	//fmt.Fprintf(w, "'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'\n",
	//"cta.Nombre", "cta.Apellidos", "cta.Puesto", "cta.Email", "cta.EmailAlt", "cta.Pass", "cta.Tel", "cta.Cel", "cta.FechaHora", "cta.CodigoCfm", "cta.Status",
	//"IdEmp", "RFC", "Nombre Empresa", "Razon Social", "Dir.Calle", "Dir.Colonia", "Dir.Entidad", "Dir.Municipio", "Dir.Cp", "Dir.Número Suc",
	//"Organiso Emp", "Otro Organismo", "Reg Org. Empresarial", "Url", "PartLinea", "ExpComer", "Descripción", "FechaHora Alta Emp.","emp.Status")
	fmt.Fprintf(w, "'%s'|'%s'|'%s'|'%t'|'%s'|'%s'|'%s'|'%s'|'%s'|'%t'\n",
	"cta.Email", "cta.EmailAlt", "cta.Pass", "cta.Status", "IdEmp", "RFC", "Nombre Empresa", "Razon Social", "FechaHora Alta Emp","emp.Status")

    q := datastore.NewQuery("Cta").Offset(pagina*lote).Limit(lote)
	for cursor := q.Run(c); ; {
		var cta model.Cta
		_, err := cursor.Next(&cta)
		if err == datastore.Done  {
			break
		}

		q2 := datastore.NewQuery("Empresa").Ancestor(cta.Key(c))
		for cursor := q2.Run(c); ; {
			var emp model.Empresa
			_, err := cursor.Next(&emp)
			if err == datastore.Done  {
				break
			}
			/*
			var entidad string
			var municipio string
			munq := datastore.NewQuery("Municipio").Filter("CveEnt =", emp.DirEnt).Filter("CveMun =", emp.DirMun).Limit(1)
			for muncur := munq.Run(c); ; {
				var mun model.Municipio
				_, err := muncur.Next(&mun)
				if err == datastore.Done  {
					break
				}
				municipio = mun.Municipio
				entidad = mun.Entidad
			}
			desc := strings.Replace(emp.Desc, "\n", " ", -1)
			desc = strings.Replace(desc, "\r", " ", -1)
			*/

			fmt.Fprintf(w, "'%s'|'%s'|'%s'|'%t'|'%s'|'%s'|'%s'|'%s'|'%s'|'%t'\n",
			cta.Email, cta.EmailAlt, cta.Pass, cta.Status, emp.IdEmp, emp.RFC, emp.Nombre, emp.RazonSoc, emp.FechaHora, emp.Status)

			//fmt.Fprintf(w, "'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%s'|'%d'|'%d'|'%s'|'%s'|'%t'\n",
			//cta.Nombre, cta.Apellidos, cta.Puesto, cta.Email, cta.EmailAlt, cta.Pass, cta.Tel, cta.Cel, cta.FechaHora, cta.CodigoCfm, cta.Status,
			//emp.IdEmp, emp.RFC, emp.Nombre, emp.RazonSoc, emp.DirCalle, emp.DirCol, entidad, municipio, emp.DirCp, emp.NumSuc,
			//emp.OrgEmp, emp.OrgEmpOtro, emp.OrgEmpReg, emp.Url, emp.PartLinea, emp.ExpComer, desc, emp.FechaHora, emp.Status)
			/*
			cta.Nombre, cta.Apellidos, cta.Puesto, cta.Email, cta.EmailAlt, cta.Pass, cta.Tel, cta.Cel, cta.FechaHora, cta.CodigoCfm, cta.Status,
			emp.IdEmp, emp.RFC, emp.Nombre, emp.RazonSoc, emp.DirCalle, emp.DirCol, entidad, municipio, emp.DirCp, emp.NumSuc,
			emp.OrgEmp, emp.OrgEmpOtro, emp.OrgEmpReg, emp.Url, emp.PartLinea, emp.ExpComer, emp.FechaHora, emp.Status)
			*/
		}

	}
}

