package main

import (
	"fmt"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/memdebug"
)

var assetBox *rice.Box

func runWeb(cfg *configData, log *logrus.Logger) {
	//r := mux.NewRouter()

	assetBox = rice.MustFindBox("assets")
	//staticFileServer := http.StripPrefix("/", http.FileServer(assetBox.HTTPBox()))
	//r.Handle("/", http.FileServer(assetBox.HTTPBox()))
	//r.Handle("/", staticFileServer).Methods("GET")
	//http.Handle("/", r)
	http.HandleFunc("/booking", mainHandler)
	http.Handle("/", http.FileServer(rice.MustFindBox("assets").HTTPBox()))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.WithField("port", cfg.Port).Info("Starting up")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	memdebug.Print(t1, r.Method, r.RequestURI)
	switch r.Method {
	case "GET":
		filename := r.RequestURI
		if filename == "/" {
			filename = "/index.html"
		}
		b, _ := assetBox.Bytes(filename[1:])
		w.Write(b)
		return
	case "POST":
		memdebug.Print(time.Now(), "Posting a booking")
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		println("name", r.FormValue("name"))
		println("bike", r.FormValue("bike"))
		println("equiry", r.FormValue("equiry"))
		println("email", r.FormValue("email"))
		println("telephone", r.FormValue("telephone"))
		println("address", r.FormValue("address"))
		println("message", r.FormValue("message"))
		b, _ := assetBox.Bytes("thanks.html")
		w.Write(b)
		sendMail("A booking", fmt.Sprintf("<h1>Booking</h1><ul><li>Name: %s<li>Email: %s<li>Bike: %s</ul>",
			r.FormValue("name"),
			r.FormValue("email"),
			r.FormValue("bike")), configDataData, logrus.New())
	}
	memdebug.Print(t1, r.Method, r.RequestURI)
}
