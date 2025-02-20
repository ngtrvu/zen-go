package httpserver

const (
	ItemsPerPage = 20
)

type HTTPServerConfig struct {
	Port             int `config:"HTTP_SERVER_PORT"`
	GracefulShutdown int `config:"HTTP_SERVER_GRACEFUL_SHUTDOWN"`
	ReadTimeout      int `config:"HTTP_SERVER_READ_TIMEOUT"`
	WriteTimeout     int `config:"HTTP_SERVER_WRITE_TIMEOUT"`
	IdleTimeout      int `config:"HTTP_SERVER_IDLE_TIMEOUT"`
}
