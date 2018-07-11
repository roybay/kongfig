package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	contentType     string = "Content-Type"
	applicationJSON string = "application/json; charset=utf-8"
)

type Service struct {
	ConnectTimeout int     `json:"connect_timeout,omitempty"`
	Name           string  `json:"name,omitempty"`
	Path           string  `json:"path,omitempty"`
	ReadTimeout    int     `json:"read_timeout,omitempty"`
	Retries        int     `json:"retries,omitempty"`
	URL            string  `json:"url,omitempty"`
	WriteTimeout   int     `json:"write_timeout,omitempty"`
	Routes         []Route `json:"routes,omitempty"`
}

type Routes struct {
	Next string  `json:"next"`
	Data []Route `json:"data"`
}

type Route struct {
	ID           string   `json:"id,omitempty"`
	Protocols    []string `json:"protocols,omitempty"`
	Methods      []string `json:"methods,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	Paths        []string `json:"paths,omitempty"`
	StripPath    bool     `json:"strip_path,omitempty"`
	PreserveHost bool     `json:"preserve_host,omitempty"`
	Service      string   `json:"-"`
}

type Config struct {
	Host     string    `json:"host"`
	HTTPS    bool      `json:"https"`
	version  string    `json:"version"`
	Services []Service `json:"services,omitempty"`
}

func (c *Config) GetRoutes(s Service) (string, error) {
	url := fmt.Sprintf("%s/services/%s/routes", adminURL(c), s.Name)
	r := Routes{}
	res, err := httpGetRequest(url, &r)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response")
	}

	routes := r.Data

	log.Printf("GET /routes - [%d]: %v", len(routes), routes)

	for _, rr := range routes {
		c.deleteRoute(rr)
	}

	return "", nil
}

func (c *Config) deleteRoute(r Route) (string, error) {
	url := fmt.Sprintf("%s/routes/%s", adminURL(c), r.ID)
	res, err := httpDeleteRequest(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("bad response")
	}

	log.Printf("DELETE /routes/%s - OK", r.ID)

	return "", nil
}

func (c *Config) CreateRoutes(s Service) (string, error) {
	url := fmt.Sprintf("%s/services/%s/routes", adminURL(c), s.Name)

	for _, r := range s.Routes {
		payload, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		res, err := httpPostRequest(url, payload)
		if err != nil {
			return "", err
		}
		if res.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("bad response")
		}
		log.Printf("[routes] POST - %d", res.StatusCode)
	}

	return "", nil
}

func httpGetRequest(url string, response interface{}) (*http.Response, error) {
	return httpRequest(http.MethodGet, url, nil, response)
}

func httpPostRequest(url string, payload []byte) (*http.Response, error) {
	return httpRequest(http.MethodPost, url, payload, nil)
}

func httpPutRequest(url string, payload []byte) (*http.Response, error) {
	return httpRequest(http.MethodPut, url, payload, nil)
}

func httpDeleteRequest(url string) (*http.Response, error) {
	return httpRequest(http.MethodDelete, url, nil, nil)
}

func httpRequest(method, url string, payload []byte, response interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set(contentType, applicationJSON)
	req.Header.Set("User-Agent", "kongfig")
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	res, err := client.Do(req)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&response)

	// DEBUG
	// dump, _ := httputil.DumpResponse(res, true)
	// log.Printf("[http]: %q", dump)

	return res, err
}

func adminURL(c *Config) string {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}
	host := "localhost:8001"
	return fmt.Sprintf("%s://%s", protocol, host)
}

func (c *Config) UpdateService(s Service) (string, error) {
	url := fmt.Sprintf("%s/services/%s", adminURL(c), s.Name)
	s.Routes = nil
	payload, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	res, err := httpPutRequest(url, payload)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response")
	}

	log.Printf("[services] PUT - %d", res.StatusCode)

	return "", nil
}
