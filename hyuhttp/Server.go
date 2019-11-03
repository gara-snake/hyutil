package hyuhttp

import (
	"fmt"
	"log"
	"net/http"
)

const defaultPort = 8080

// Server Httpサーバ実装
type Server struct {
	Port uint
}

// Go ListenAndServe を開始する
func (srv *Server) Go() {

	if srv.Port <= 0 {
		srv.Port = defaultPort
	}

	p := ":" + fmt.Sprint(srv.Port)

	log.Println("Http server start on " + p)

	http.ListenAndServe(p, nil)

}
