package client

type Client struct {
	Id     string
	Secret string
	URIs   []string
	Scope  string
}

var C = Client{
	Id:     "osemisan-client-id-1",
	Secret: "osemisan-client-secret-1",
	URIs:   []string{"http://localhost:9000/callback"},
	Scope:  "abura kuma",
}
