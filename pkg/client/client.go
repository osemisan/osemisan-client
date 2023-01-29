package client

type Client struct {
	Id     string
	Secret string
	URIs   []string
	Scope  string
}

var C = Client{
	"osemisan-client-id-1",
	"osemisan-client-secret-1",
	[]string{"http://localhost:9000"},
	"abura kuma",
}
