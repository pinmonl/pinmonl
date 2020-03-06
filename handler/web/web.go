package web

// ServerOpts defines the options of Server initiation.
type ServerOpts struct {
	DevServer string
}

// Server handles webui mainly.
type Server struct {
	devSvr string
}

// NewServer creates Server.
func NewServer(opts ServerOpts) *Server {
	return &Server{
		devSvr: opts.DevServer,
	}
}
