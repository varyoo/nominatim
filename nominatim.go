package nominatim

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Email      string
	HTTPClient http.Client
}

const (
	FieldPostalCode string = "postalcode"
	FieldCountry    string = "country"
	FieldStreet     string = "street"
	FieldCity       string = "city"
)

type Street struct {
	ValidNumber bool
	HouseNumber int64
	StreetName  string
}

func (s Street) String() string {
	if s.ValidNumber {
		return fmt.Sprintf("%d %s", s.HouseNumber, s.StreetName)
	} else {
		return s.StreetName
	}
}

func (c *Client) Lookup(fields url.Values) (io.ReadCloser, error) {
	fields.Set("email", c.Email)
	fields.Set("format", "json")

	req, err := http.NewRequest("GET",
		"https://nominatim.openstreetmap.org/search?"+fields.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
