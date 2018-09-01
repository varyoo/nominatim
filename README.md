Go client of https://wiki.openstreetmap.org/wiki/Nominatim

```Go
type Place struct {
    Name    string  `json:"display_name"`
    Lat     string  `json:"lat"`
    Lon     string  `json:"lon"`
}

search := url.Values{}
search.Set(nominatim.FieldCountry, "France")
search.Set(nominatim.FieldCity, "Poitiers")
search.Set(nominatim.FieldStreet, nominatim.Street{
    HouseNumber: 47,
    ValidNumber: true,
    StreetName: "Avenue des Terrasses",
}.String())

jsonResp, _ := nominatim.Client{}.Lookup(search)

defer jsonResp.Close()

var places []Place

json.NewDecoder(jsonResp).Decode(&places)
```

nominatim/service:

```Go

import "github.com/varyoo/nominatim/service"

type SearchJob struct {
    City    string
    Street  nominatim.Street
}

func (j SearchJob) Search() url.Values {
   s := url.Values{}
   s.Set(nominatim.FieldCity, j.City)
   s.Set(nominatim.FieldStreet, j.Street.String())
   return s
}

// SetCoordinates is called as soon as the Nominatim search completes
func (j SearchJob) SetCoordinates(jsonResp io.ReadCloser) error {
    defer jsonResp.Close()

    var places []Place

    if err := json.NewDecoder(jsonResp).Decode(&places); err != nil {
        return err
    }

    for _, p := range places {
        log.Println(p.Name)
    }
}

func localize() {
    s := New(nominatim.Client{})

    // start the Nominatim Goroutine
    go s.Go()

    ctx, _ := context.WithTimeout(context.Backgroudn(), time.Second * 3)

    s.Localize(ctx, SearchJob{"Poitiers", nominatim.Street{
        HouseNumber: 47,
        ValidNumber: true,
        StreetName: "Avenue des Terrasses",
    })
}
```
