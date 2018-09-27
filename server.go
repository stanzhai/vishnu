package vishnu

type Server struct {
	bridge    string
	localAddr string
	port      int
}

func NewServer(bridge string, port int) *Server {
	return &Server{
		bridge: bridge,
		port:   port,
	}
}

func (s *Server) Run() {

}
