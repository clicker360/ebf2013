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
	IdOft       string `json:"idoft"`
	IdEmp       string `json:"idemp"`
	IdCat       int `json:"idcat"`
	Empresa		string `json:"empresa"`
	Oferta		string `json:"oferta"`
	Descripcion	string `json:"desc"`
	Codigo      string `json:"codigo,omitempty"`
	Precio      string `json:"precio,omitempty"`
	Descuento   string `json:"descuento,omitempty"`
	Promocion	string `json:"promocion,omitempty"`
	Enlinea     bool `json:"enlinea"`
	Url         string `json:"url,omitempty"`
	Meses       string `json:"meses,omitempty"`
	FechaHoraPub    time.Time `json:"fechapub"`
	StatusPub   bool `json:"publicada"`
	FechaHora   time.Time `json:"timestamp"`
	BlobKey	appengine.BlobKey `json:"blobkey"`
    ImageSmall  string `json:"imagesmall,omitempty"`
    ImageBig    string `json:"imagebig,omitempty"`
}

type OfertaSucursal struct {
	IdOft       string `json:"idoft"`
	IdEmp       string `json:"idemp"`
	IdSuc       string `json:"idsuc"`
	Sucursal    string `json:"sucursal"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Empresa     string `json:"empresa"`
	Oferta      string `json:"oferta"`
	Descripcion		string `json:"desc"`
	Promocion	string `json:"promocion,omitempty"`
	Precio      string `json:"precio,omitempty"`
	Descuento   string `json:"descuento,omitempty"`
	Enlinea     bool `json:"enlinea"`
	Url         string `json:"url,omitempty"`
	StatusPub   bool `json:"publicada"`
	FechaHora	time.Time `json:"timestamp"`
	IdCat       int `json:"idcat"`
	Categoria   string `json:"categoria"`
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
	IdCat       int `json:"idcat"`
	Categoria   string `json:"categoria"`
	Selected	string `datastore:"-" json:"selected,omitempty"`
}

type OfertaPalabra struct {
	IdOft      string `json:"idoft"`
	IdEmp      string `json:"idemp"`
	Palabra    string `json:"palabra"`
	FechaHora	time.Time `json:"timestamp"`
}

type OfertaEstado struct {
	IdOft      string `json:"idoft"`
	IdEnt      string `json:"ident"`
}

func (e *Empresa) OfertaKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Oferta", id, 0, e.Key(c))
}

func (e *Empresa) GetOferta(c appengine.Context, id string) (*Oferta, error) {
	var o Oferta
	if err := datastore.Get(c, e.OfertaKey(c, id), &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func GetOferta(c appengine.Context, id string) *Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdOft =", id).Limit(1)
	for i := q.Run(c); ; {
		var o Oferta
		_, err := i.Next(&o)
		if err == datastore.Done {
			break
		}
		return &o
	}
	return nil
}

func (e *Empresa) GetOfertaSucursales(c appengine.Context, id string) (*[]OfertaSucursal, error) {
	q := datastore.NewQuery("OfertaSucursal").Ancestor(e.OfertaKey(c, id))
	n, _ := q.Count(c)
	ofersuc := make([]OfertaSucursal, 0, n)
	if _, err := q.GetAll(c, &ofersuc); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &ofersuc, nil
}

func GetOfertaSucursales(c appengine.Context, id string) (*[]OfertaSucursal, error) {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", id)
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

func (e *Empresa) PutOferta(c appengine.Context, o *Oferta) (*datastore.Key, error) {
    if(o.IdOft == "") {
		o.IdOft = RandId(20)
		_ = PutChangeControl(c, o.IdOft, "Oferta", "A")
	} else {
		_ = PutChangeControl(c, o.IdOft, "Oferta", "M")
	}
    key, err := datastore.Put(c, e.OfertaKey(c, o.IdOft), o)
	if err != nil {
		return nil, err
	}
    /* 
		Recorre las relaciones oferta sucursal para copiar los datos y tenerlos listos
        para lecturas posteriores
	*/
	if ofsucs, err := e.GetOfertaSucursales(c, o.IdOft); err != nil {
		c.Infof("GetOfertaSucursales warning IdOft: %v", o.IdOft)
        return key, err
    } else {
        for _,os:= range *ofsucs {
            var ofsuc OfertaSucursal
            ofsuc.IdOft = os.IdOft
            ofsuc.IdSuc = os.IdSuc
            ofsuc.IdEmp = os.IdEmp
            ofsuc.Sucursal = os.Sucursal
            ofsuc.Lat = os.Lat
            ofsuc.Lng = os.Lng
            ofsuc.Empresa = o.Empresa
            ofsuc.Oferta = o.Oferta
            ofsuc.Descripcion = o.Descripcion
            ofsuc.Promocion = o.Promocion
            ofsuc.Precio = o.Precio
            ofsuc.Descuento = o.Descuento
            ofsuc.Enlinea = o.Enlinea
            ofsuc.Url = o.Url
            ofsuc.StatusPub = o.StatusPub
            ofsuc.FechaHora = time.Now().Add(time.Duration(GMTADJ)*time.Second)

            e.PutOfertaSucursal(c, &ofsuc)
            TouchSuc(c, os.IdSuc)
        }
	}
	return key, nil
}

/*
	Llenar primero struct de OfertaSucursal y luego guardar
*/
func (e *Empresa) PutOfertaSucursal(c appengine.Context, ofsuc *OfertaSucursal) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaSucursal", ofsuc.IdOft+ofsuc.IdSuc, 0, e.OfertaKey(c, ofsuc.IdOft)), ofsuc)
	if err != nil {
		return err
	}
	_ = TouchSuc(c, ofsuc.IdSuc)
	return nil
}

/*
	Llenar primero struct de OfertaPalabra y luego guardar
*/
func (e *Empresa) PutOfertaPalabra(c appengine.Context, idoft string, op *OfertaPalabra) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaPalabra", e.IdEmp+op.Palabra, 0, e.OfertaKey(c, idoft)), op)
	if err != nil {
		return err
	}
	return nil
}

func (e *Empresa) DelOferta(c appengine.Context, id string) error {
    o, err := e.GetOferta(c, id)
    if err != nil {
       return err
    }
    if err := blobstore.Delete(c, o.BlobKey); err != nil {
        return err
    }
    if err:= DelOfertaSucursales(c, id); err != nil {
        return err
    }
    if err := DelOfertaPalabras(c, id); err != nil {
        return err
    }
    if err := DelOfertaSearchData(c, e.OfertaKey(c, id)); err != nil {
        return err
    }
    if err := datastore.Delete(c, e.OfertaKey(c, id)); err != nil {
        return err
    }
    _ = PutChangeControl(c, o.IdOft, "Oferta", "B")
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
func (e *Empresa) PutOfertaEstado(c appengine.Context, idoft string, edomap map[string]string) error {
	for k, v := range edomap {
		var oe OfertaEstado
		oe.IdOft = v
		oe.IdEnt = k
		_, err := datastore.Put(c, datastore.NewKey(c, "OfertaEstado", v+k, 0, e.OfertaKey(c, idoft)), &oe)
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

func (e *Empresa) ListOf(c appengine.Context) *[]Oferta {
	q := datastore.NewQuery("Oferta").Ancestor(e.Key(c)).Limit(500)
	n, _ := q.Count(c)
	ofertas := make([]Oferta, 0, n)
	if _, err := q.GetAll(c, &ofertas); err != nil {
		return nil
	}
	sortutil.AscByField(ofertas, "Oferta")
	return &ofertas
}

func ListOf(c appengine.Context, id string) *[]Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdEmp =", id).Limit(500)
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
		n, _ := q.Count(c)
		categorias = make([]Categoria, 0, n)
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
