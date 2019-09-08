package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/steveoc64/gomail"
	"github.com/steveoc64/smsbroadcast"

	rice "github.com/GeertJohan/go.rice"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/memdebug"
)

type WebHandler struct {
	cfg    *configData
	log    *logrus.Logger
	assets *rice.Box
	mailer *gomail.Mailer
	sms    *smsbroadcast.SMS
	db     *bolt.DB
}

func NewWebHandler(cfg *configData, log *logrus.Logger) *WebHandler {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("bookings"))
		tx.CreateBucketIfNotExists([]byte("ids"))
		return nil
	})

	return &WebHandler{
		cfg:    cfg,
		log:    log,
		assets: rice.MustFindBox("assets"),
		mailer: gomail.New(cfg.Mail.Server, cfg.Mail.Username, cfg.Mail.Password),
		sms:    smsbroadcast.New(cfg.SMS.API, cfg.SMS.Username, cfg.SMS.Password, cfg.SMS.Destination, cfg.SMS.Source),
		db:     db,
	}

}

func (h *WebHandler) Run() {
	http.HandleFunc("/booking", h.bookings)
	http.Handle("/", http.FileServer(rice.MustFindBox("assets").HTTPBox()))

	addr := fmt.Sprintf("%s:%d", h.cfg.Address, h.cfg.Port)
	h.log.WithField("port", h.cfg.Port).Info("Starting up")
	h.log.Fatal(http.ListenAndServe(addr, nil))
}

func getIPAdress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

func (h *WebHandler) bookings(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	switch r.Method {
	case "GET":
		// so we never actually getting into here
		/*
			filename := r.RequestURI
			if filename == "/" {
				filename = "/index.html"
			}
			b, _ := h.assets.Bytes(filename[1:])
			w.Write(b)
			h.log.WithField("uri", filename).Print("GET")
		*/
		return
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// immediate feedback on the booking form
		b, _ := h.assets.Bytes("thanks.html")
		w.Write(b)

		ip := getIPAdress(r)
		id := h.newBooking(Booking{
			IP:        ip,
			Name:      r.FormValue("name"),
			Bike:      r.FormValue("bike"),
			Enquiry:   r.FormValue("enquiry"),
			Email:     r.FormValue("email"),
			Telephone: r.FormValue("telephone"),
			Address:   r.FormValue("address"),
			Message:   r.FormValue("message"),
			Date:      time.Now(),
		})
	}
	memdebug.Print(t1, r.Method, r.RequestURI)
}
