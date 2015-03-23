package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/ryanbillingsley/blackmagic"
)

func main() {

	mongoUrl := flag.String("mongo", "localhost", "The mongo db address.  It can be as simple as `localhost` or involved as `mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb`")
	databaseName := flag.String("db", "blackmagic", "The name of the database you are connecting to.  Defaults to blackmagic")
	port := flag.Int("p", 8080, "The port that the server will listen on.")
	flag.Parse()

	database := blackmagic.NewDatabase(*mongoUrl, *databaseName)

	server := NewServer(database)

	server.Start(*port)
}

type Server interface {
	Start() error
}

type server struct {
	Database blackmagic.Database
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

func NewServer(db blackmagic.Database) *server {
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
