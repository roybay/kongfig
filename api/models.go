package api

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
