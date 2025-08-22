package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/heretic1321/pokedex/internal/errorhandler"
)

type Client struct {
	Base string
	http *http.Client
}

const defaultBaseUrl = "https://pokeapi.co/api/v2/"

func New(httpClient *http.Client, base string) *Client {
  if httpClient == nil { httpClient = http.DefaultClient }
  if base == "" { base = defaultBaseUrl }
  return &Client{Base: base, http: httpClient}
}


func (c *Client) DoJSON(url string, dst any) error {
	req, err := http.NewRequest("GET", url, nil )
	if err != nil {return errorhandler.Handle(err)}
	res, err := c.http.Do(req)

	if err != nil {return errorhandler.Handle(err)}
	
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300{
		return fmt.Errorf("bad status: %s", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
  err = decoder.Decode(dst)

	if err != nil {
		return errorhandler.Handle(err)
	}

	return nil
}
