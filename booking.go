package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

type Booking struct {
	IP        string
	Name      string
	Bike      string
	Enquiry   string
	Email     string
	Telephone string
	Address   string
	Message   string
	Date      time.Time
}

func (h *WebHandler) newBooking(booking Booking) int {
	h.log.WithFields(logrus.Fields{
		"ip":      booking.IP,
		"name":    booking.Name,
		"bike":    booking.Bike,
		"enq":     booking.Enquiry,
		"email":   booking.Email,
		"tel":     booking.Telephone,
		"address": booking.Address,
		"message": booking.Message,
		"date":    booking.Date,
	}).Print("NewBooking")

	id := 0

	// Store the data in the DB
	err := h.db.Update(func(tx *bolt.Tx) error {
		bookingBytes, err := json.Marshal(booking)
		if err != nil {
			return err
		}
		idBkt := tx.Bucket([]byte("ids"))
		if idBkt == nil {
			return nil
		}
		// get the bookingID
		idBytes := string(idBkt.Get([]byte("bookingID")))
		id, _ = strconv.Atoi(idBytes)
		id++
		idBytes = fmt.Sprintf("%d", id)
		idBkt.Put([]byte("bookingID"), []byte(idBytes))

		bookingBkt := tx.Bucket([]byte("bookings"))
		if bookingBkt == nil {
			return nil
		}
		return bookingBkt.Put([]byte(idBytes), bookingBytes)
	})

	// Generate an email to the person that placed the booking
	userMail := fmt.Sprintf(`
<div style="background:#90a1cf;border-bottom:20px solid #32439d">
<img src="http://cycle2u.com.au/images/logo.jpg">
</div>
<h3>Thanks for your Enquiry (%s)</h3>
<p>We will get back to you soon.

<p><strong>Bike: </strong> %s</p>
<p><strong>Telephone: </strong> %s</p>
<p><strong>Address: </strong> %s</p>
<p><strong>Message: </strong></p>
<p>%s</p>
<p><strong>Reference: #2019%04d</strong></p>
<p><img src="http://cycle2u.com.au/images/logos.jpg">
<div style="background:#90a1cf;border-top:20px solid #32439d">
<p><img src="http://cycle2u.com.au/images/experience.png">
</div>
`,
		booking.Enquiry,
		booking.Bike,
		booking.Telephone,
		booking.Address,
		booking.Message,
		id,
	)
	err = h.mailer.Send(h.cfg.Mail.From, booking.Email, "Cycle2U Booking Confirmation", userMail, nil)
	if err != nil {
		h.log.WithError(err).Error("Mailing customer")
	}

	// Generate emails to the service reps
	svcMail := fmt.Sprintf(`
<div style="background:#90a1cf;border-bottom:20px solid #32439d">
<img src="http://cycle2u.com.au/images/logo.jpg">
</div>
<h3>Cycle2U Website Enquiry (%s)</h3>
<p><strong>Name: </strong> %s </p> 
<p><strong>Bike: </strong> %s </p> 
<p><strong>Email Address: </strong> %s </p> 
<p><strong>Home Address: </strong> %s
<p><a href="http://maps.google.com/?q=%s">Click for MAP</a>
<p><strong>Telephone: </strong> %s </p> 
<hr><p><strong>Message: </strong></p><p> %s </p> 
<p><strong>Reference: #2019%04d</strong></p>
<hr>
<p><img src="http://cycle2u.com.au/images/logos.jpg">
<p><small>This message was sent from the IP Address: %s at %s on %s</small></p>
<div style="background:#90a1cf;border-top:20px solid #32439d">
<p><img src="http://cycle2u.com.au/images/experience.png">
`,
		booking.Enquiry,
		booking.Name,
		booking.Bike,
		booking.Email,
		booking.Address,
		url.QueryEscape(booking.Address),
		booking.Telephone,
		booking.Message,
		id,
		booking.IP,
		time.Now().Format(time.RFC850),
		h.cfg.Name,
	)
	err = h.mailer.Send(h.cfg.Mail.From,
		h.cfg.Mail.Email,
		fmt.Sprintf("Cycle2U Booking - %s @%s <%s>", booking.Bike, booking.Name, booking.Email),
		svcMail,
		h.cfg.Mail.BCC,
	)
	if err != nil {
		h.log.WithError(err).Error("Mailing service rep")
	}

	// Generate an SMS to the prime service rep
	smsText := fmt.Sprintf("#2019%04d\n%s\nPh: %s\n\n%s\n\n%s\n---\n%s",
		id,
		booking.Name,
		booking.Telephone,
		booking.Bike,
		booking.Address,
		booking.Message)
	if len(smsText) > 160 {
		smsText = smsText[:160]
	}

	rsp, err := h.sms.Send(smsText)
	if err != nil {
		h.log.WithError(err).Error("Sending SMS")
	}
	h.log.WithField("smsref", rsp).Info("SMS Response")
	return id
}
