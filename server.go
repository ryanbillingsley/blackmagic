package blackmagic

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Server interface {
	Start() error
}

type server struct {
	Database Database
	Router   *mux.Router
}

func (s *server) Start(port int) {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", s.homeHandler)

	s.Router = router

	negroni := negroni.Classic()
	negroni.UseHandler(router)

	negroni.Run(fmt.Sprintf(":%v", port))
}

func NewServer(db Database) *server {
	return &server{Database: db}
}

func (s *server) homeHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(res, "Not found", 404)
		return
	}
	if req.Method != "GET" {
		http.Error(res, "Method not allowed", 405)
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(res, "Welcome")
}
