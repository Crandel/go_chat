package setup

// Database ...
type Database struct {
	Type string
	Path string
}

// Server ...
type Server struct {
	Host     string
	Domain   string
	Port     int
	UseHTTPS bool
	HTTPS    HTTPS
}

// HTTPS ...
type HTTPS struct {
	Port        int
	Certificate string
	Key         string
}
