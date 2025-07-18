package main

import (
	server "onfly-api/cmd/fiber_http"
	"onfly-api/cmd/fiber_http/jwt"
	"onfly-api/internal/infrasctructure/persistence"
)

func main() {

	server.StartServerHttp()
	persistence.InitDB()
	jwt.InitJWT()

}
