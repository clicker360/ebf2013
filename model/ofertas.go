package model

import (
    "appengine"
    "appengine/datastore"
	"appengine/blobstore"
	"appengine/memcache"
	"encoding/json"
	"sortutil"
	"time"
)

type Oferta struct {
	IdOft       string
	IdEmp       string
	IdCat       int
	Empresa		string
	Oferta		string
	Descripcion		string
	Codigo      string
	Precio      string
	Descuento   string
	Promocion	string
	Enlinea     bool
	Url         string
	Meses       string
	FechaHoraPub    time.Time
	StatusPub   bool
	FechaHora   time.Time
	BlobKey	appengine.BlobKey
}

type OfertaSucursal struct {
	IdOft       string
	IdEmp       string
	IdSuc       string
	Sucursal    string
	Lat         float64
	Lng         float64
	Empresa     string
	Oferta      string
	Descripcion		string
	Promocion	string
	Precio      string
	Descuento   string
	Enlinea     bool
	Url         string
	StatusPub   bool
	FechaHora	time.Time
	IdCat       int
	Categoria   string
}

type SearchData struct {
	Sid			string
	Kind		string
	Field		string
	Value		string
	IdCat		int
	Enlinea		bool
	FechaHora	time.Time
}

type Categoria struct {
	IdCat       int
	Categoria   string
	Selected	string `datastore:"-"`
}

type OfertaPalabra struct {
	IdOft      string
	IdEmp      string
	Palabra    string
	FechaHora	time.Time
}

type OfertaEstado struct {
	IdOft      string
	IdEnt      string
}

func (r *Oferta) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Oferta", r.IdOft, 0, nil)
}

func (r *Oferta) DelOferta(c appengine.Context) error {
	if err := blobstore.Delete(c, r.BlobKey); err != nil {
		return err
	}
    if err := datastore.Delete(c, r.Key(c)); err != nil {
		return err
	}
	return nil
}

func GetOferta(c appengine.Context, id string) (*Oferta, *datastore.Key) {
	q := datastore.NewQuery("Oferta").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e Oferta
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		// Regresa la oferta
		return &e, key
	}
	// Regresa un cascarón
	var e Oferta
	e.IdEmp = "none";
	e.IdOft = "none";
	e.IdCat = 0;
	e.BlobKey = "none";
	return &e, nil
}

func PutOferta(c appengine.Context, oferta *Oferta) error {
	if oferta.BlobKey == "" {
		oferta.BlobKey = "none"
	}
	_ = PutChangeControl(c, oferta.IdOft, "Oferta", "M")
	_, err := datastore.Put(c, oferta.Key(c), oferta)
	if err != nil {
		return err
	}
	/* 
		relación oferta sucursal 
	*/
	ofsucs, _ := GetOfertaSucursales(c, oferta.IdOft)
	for _,os:= range *ofsucs {
		var ofsuc OfertaSucursal
		ofsuc.IdOft = os.IdOft
		ofsuc.IdSuc = os.IdSuc
		ofsuc.IdEmp = os.IdEmp
		ofsuc.Sucursal = os.Sucursal
		ofsuc.Lat = os.Lat
		ofsuc.Lng = os.Lng
		ofsuc.Empresa = oferta.Empresa
		ofsuc.Oferta = oferta.Oferta
		ofsuc.Descripcion = oferta.Descripcion
		ofsuc.Promocion = oferta.Promocion
		ofsuc.Precio = oferta.Precio
		ofsuc.Descuento = oferta.Descuento
		ofsuc.Enlinea = oferta.Enlinea
		ofsuc.Url = oferta.Url
		ofsuc.StatusPub = oferta.StatusPub
		ofsuc.FechaHora = time.Now().Add(time.Duration(GMTADJ)*time.Second)
		oferta.PutOfertaSucursal(c, &ofsuc)
		TouchSuc(c, os.IdSuc)
	}

	return nil
}

func NewOferta(c appengine.Context, oferta *Oferta) (*Oferta, error) {
	oferta.IdOft = RandId(20)
    _, err := datastore.Put(c, datastore.NewKey(c, "Oferta", oferta.IdOft, 0, nil), oferta)
	if err != nil {
		return nil, err
	}
	_ = PutChangeControl(c, oferta.IdOft, "Oferta", "A")
	return oferta, nil
}

func GetOfertaSucursales(c appengine.Context, idoft string) (*[]OfertaSucursal, error) {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", idoft)
	n, _ := q.Count(c)
	ofersuc := make([]OfertaSucursal, 0, n)
	if _, err := q.GetAll(c, &ofersuc); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &ofersuc, nil
}

func GetOfertaPalabras(c appengine.Context, idoft string, idemp string) *[]OfertaPalabra {
	q := datastore.NewQuery("OfertaPalabra")
	if(idemp != "") {
		q = q.Filter("IdEmp =", idemp)
	} else {
		q = q.Filter("IdOft =", idoft)
	}
	n, _ := q.Count(c)
	op := make([]OfertaPalabra, 0, n)
	if _, err := q.GetAll(c, &op); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil
		}
	}
	return &op
}

func GetOfertaSucursalesGeo(c appengine.Context, lat string, lng string, rad string) (*Sucursal, error) {
	/*
	q := datastore.NewQuery("Sucursal")
	for i := q.Run(c); ; {
		var s Sucursal
        _, err := i.Next(&s)
		if err == datastore.Done {
			break
        }
		geo1, _ := strconv.ParseFloat(s.Geo1, 64)
		geo2, _ := strconv.ParseFloat(s.Geo2, 64)
		sqdist := (lat - geo1) * (lat - geo1)  + (long - geo2) * (long - geo2);
		if ( sqdist <= rad * rad) {
			fmt.Fprintf(w, "lat, long: %s, %s\n", s.Geo1, s.Geo2);
		}
	}
	*/
	return nil,nil
}

func GetCategoria(c appengine.Context, id int) *Categoria {
	q := datastore.NewQuery("Categoria").Filter("IdCat =", id).Limit(1)
	for i := q.Run(c); ; {
		var c Categoria
		_, err := i.Next(&c)
		if err == datastore.Done {
			break
		}
		return &c
	}
	return nil
}

/*
	Llenar primero struct de OfertaSucursal y luego guardar
*/
func (r *Oferta) PutOfertaSucursal(c appengine.Context, ofsuc *OfertaSucursal) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaSucursal", r.IdOft+ofsuc.IdSuc, 0, r.Key(c)), ofsuc)
	if err != nil {
		return err
	}
	_ = TouchSuc(c, ofsuc.IdSuc)
	return nil
}

/*
	Llenar primero struct de OfertaPalabra y luego guardar
*/
func (r *Oferta) PutOfertaPalabra(c appengine.Context, op *OfertaPalabra) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaPalabra", r.IdEmp+op.Palabra, 0, r.Key(c)), op)
	if err != nil {
		return err
	}
	return nil
}

func DelOferta(c appengine.Context, id string) error {
	q := datastore.NewQuery("Oferta").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e Oferta
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := blobstore.Delete(c, e.BlobKey); err != nil {
			return err
		}
		if err:= DelOfertaSucursales(c, id); err != nil {
			return err
		}
		if err := DelOfertaPalabras(c, id); err != nil {
			return err
		}
		if err := DelOfertaSearchData(c, key); err != nil {
			return err
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
		_ = PutChangeControl(c, e.IdOft, "Oferta", "B")
	}
	return nil
}

/*
	Método para borrar todas las palabras de SearchData
*/
func DelOfertaSearchData(c appengine.Context, key *datastore.Key) error {
	q := datastore.NewQuery("SearchData").Filter("Sid =", key.Encode()).KeysOnly()
	n, _ := q.Count(c)
	sd := make([]*datastore.Key, 0, n)
	if _, err := q.GetAll(c, &sd); err != nil {
		return nil
	}
	if err := datastore.DeleteMulti(c, sd); err != nil {
		return err
	}
	return nil
}

/*
	Método para borrar todas las sucursales de una oferta
*/
func DelOfertaSucursales(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaSucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

/*
	Método para borrar todas una sucursales de todas las ofertas
*/
func DelSucursalesOferta(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdSuc =", id)
	for i := q.Run(c); ; {
		var e OfertaSucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

/*
	Método para borrar todas las sucursales de una oferta
*/
func DelOfertaSucursal(c appengine.Context, idoft string, idsuc string) error {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", idoft).Filter("IdSuc =", idsuc)
	for i := q.Run(c); ; {
		var e OfertaSucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	_ = TouchSuc(c, idsuc)
	return nil
}

/*
	Las palabras clave asociadas a una oferta se ponen con idoft="none" todas juntas
*/
func DelOfertaPalabras(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		e.IdOft = "none";
		/*
			En realidad no se borra ningun entity, solo se desliga la oferta
			La palabra continua perteneciendo a la empresa para uso de las demás
			ofertas
		*/
		_, err = datastore.Put(c, key, &e)
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

/*
	Las palabras clave asociadas a una oferta se borran todas juntas
*/
func RmOfertaPalabras(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}


/*
	Las palabra clave se borra de Empresa 
*/
func RmOfertaPalabra(c appengine.Context, id string, palabra string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdEmp =", id).Filter("Palabra =", palabra)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		/*
			Aquí si se borra el entity
		*/
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

func DelOfertaPalabra(c appengine.Context, id string, palabra string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id).Filter("Palabra =", palabra)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		e.IdOft = "none";
		/*
			En realidad no se borra ningun entity, solo se desliga la oferta
			La palabra continua perteneciendo a la empresa para uso de las demás
			ofertas
		*/
		_, err = datastore.Put(c, key, &e)
	}
	return nil
}

/*
 Métodos de OfertaEstado
*/
func (r *Oferta) PutOfertaEstado(c appengine.Context, edomap map[string]string) error {
	for k, v := range edomap {
		var e OfertaEstado
		e.IdOft = v
		e.IdEnt = k
		_, err := datastore.Put(c, datastore.NewKey(c, "OfertaEstado", v+k, 0, r.Key(c)), &e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Oferta) DelOfertaEstado(c appengine.Context) error {
	q := datastore.NewQuery("OfertaEstado").Filter("IdOft =", r.IdOft)
	for i := q.Run(c); ; {
		var e OfertaEstado
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

func ListOf(c appengine.Context, IdEmp string) *[]Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdEmp =", IdEmp).Limit(500)
	n, _ := q.Count(c)
	ofertas := make([]Oferta, 0, n)
	if _, err := q.GetAll(c, &ofertas); err != nil {
		return nil
	}
	sortutil.AscByField(ofertas, "Oferta")
	return &ofertas
}

func ListCat(c appengine.Context, IdCat int) *[]Categoria {
	var categorias []Categoria
	if item, err := memcache.Get(c, "categorias"); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("Categoria")
		//n, _ := q.Count(c)
		//cats := make([]Categoria, 0, n)
		if _, err := q.GetAll(c, &categorias); err != nil {
			return nil
		}
		b, _ := json.Marshal(categorias)
		item := &memcache.Item{
			Key:   "categorias",
			Value: b,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("memcache.Add Categoria : %v", err)
		}
	} else {
		c.Infof("Memcache activo: %v", item.Key)
		if err := json.Unmarshal(item.Value, &categorias); err != nil {
			c.Errorf("Memcache Unmarchalling item: %v", err)
		}
	}
	for i, _ := range categorias {
		if(IdCat == categorias[i].IdCat) {
			categorias[i].Selected = `selected`
		}
	}
	return &categorias
}
