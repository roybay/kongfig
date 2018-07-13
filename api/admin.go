package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	contentType     string = "Content-Type"
	applicationJSON string = "application/json; charset=utf-8"
)

type Client struct {
	config  *Config
	client  *http.Client
	BaseURL string
}

func configFromPath(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := Config{}
	json.NewDecoder(file).Decode(&c)
	return &c, nil
}

func NewClient(filePath string) *Client {
	config, err := configFromPath(filePath)
	if err != nil {
		// FIXME bail
	}

	c := &Client{
		config:  config,
		client:  &http.Client{Timeout: time.Duration(5 * time.Second)},
		BaseURL: adminURL(config),
	}

	return c
}

func adminURL(c *Config) string {
	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, c.Host)
}

func (c *Client) GetRoutes(s Service) ([]Route, error) {
	r := Routes{}
	res, err := c.httpGetRequest(fmt.Sprintf("%s/services/%s/routes", c.BaseURL, s.Name), &r)
	if err != nil {
		return r.Data, err
	}
	if res.StatusCode != http.StatusOK {
		return r.Data, fmt.Errorf("bad response")
	}

	return r.Data, nil
}

func (c *Client) DeleteRoutes(s Service) error {
	routes, err := c.GetRoutes(s)
	if err != nil {
		return err
	}

	log.Printf("GET /routes - [%d]: %v", len(routes), routes)

	for _, r := range routes {
		// FIXME err check
		c.DeleteRoute(r)
	}

	return nil
}

func (c *Client) DeleteRoute(r Route) error {
	res, err := c.httpDeleteRequest(fmt.Sprintf("%s/routes/%s", c.BaseURL, r.ID))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad response")
	}

	log.Printf("route [%s] deleted (%d)", r.ID, res.StatusCode)

	return nil
}

func (c *Client) CreateRoutes(s Service) (string, error) {
	url := fmt.Sprintf("%s/services/%s/routes", c.BaseURL, s.Name)

	for _, r := range s.Routes {
		payload, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		res, err := c.httpPostRequest(url, payload)
		if err != nil {
			return "", err
		}
		if res.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("bad response")
		}
		// FIXME no route id - 204
		log.Printf("route [%s] created (%d)", s.Name, res.StatusCode)
	}

	return "", nil
}

func (c *Client) UpdateAllRecursively() {
	for _, s := range c.config.Services {
		// FIXME err check
		c.UpdateService(s)
		c.DeleteRoutes(s)
		c.CreateRoutes(s)
	}
}

func (c *Client) UpdateService(s Service) error {
	url := fmt.Sprintf("%s/services/%s", c.BaseURL, s.Name)
	s.Routes = nil
	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}
	res, err := c.httpPutRequest(url, payload)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response")
	}

	log.Printf("service [%s] created (%d)", s.Name, res.StatusCode)

	return nil
}

func (c *Client) httpGetRequest(url string, response interface{}) (*http.Response, error) {
	return c.httpRequest(http.MethodGet, url, nil, response)
}

func (c *Client) httpPostRequest(url string, payload []byte) (*http.Response, error) {
	return c.httpRequest(http.MethodPost, url, payload, nil)
}

func (c *Client) httpPutRequest(url string, payload []byte) (*http.Response, error) {
	return c.httpRequest(http.MethodPut, url, payload, nil)
}

func (c *Client) httpDeleteRequest(url string) (*http.Response, error) {
	return c.httpRequest(http.MethodDelete, url, nil, nil)
}

func (c *Client) httpRequest(method, url string, payload []byte, response interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set(contentType, applicationJSON)
	req.Header.Set("User-Agent", "kongfig")
	res, err := c.client.Do(req)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&response)

	// DEBUG
	// dump, _ := httputil.DumpResponse(res, true)
	// log.Printf("[http]: %q", dump)

	return res, err
}
