package bootstrap

import "os"

type Config struct {
	Addr string
}

func LoadConfig() Config {
	addr := os.Getenv("EINO_WORKBENCH_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	return Config{Addr: addr}
}
