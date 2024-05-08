package server

type RedisServer struct {
	port int
}

func New(port int) *RedisServer {
	return &RedisServer{
		port: port,
	}
}

func (s *RedisServer) Listen() {

	//net.Listen()

}
