package handlers

import (
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/models"
	"github.com/saalikmubeen/goravel/mailer"
)

type Handlers struct {
	App    *goravel.Goravel
	Models *models.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "bar"

	h.App.Session.Put(r.Context(), "foo", myData)

	myValue := h.App.Session.GetString(r.Context(), "foo")

	vars := make(jet.VarMap)
	vars.Set("foo", myValue)

	err := h.App.Render.JetPage(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// JSON is the handler to demonstrate json responses
func (h *Handlers) JSON(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID      int64    `json:"id"`
		Name    string   `json:"name"`
		Hobbies []string `json:"hobbies"`
	}

	payload.ID = 10
	payload.Name = "Jack Jones"
	payload.Hobbies = []string{"karate", "tennis", "programming"}

	err := h.App.WriteJSON(w, http.StatusOK, goravel.Response{
		"data": payload,
	})
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// XML is the handler to demonstrate XML responses
func (h *Handlers) XML(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		ID      int64    `xml:"id"`
		Name    string   `xml:"name"`
		Hobbies []string `xml:"hobbies>hobby"`
	}

	var payload Payload
	payload.ID = 10
	payload.Name = "John Smith"
	payload.Hobbies = []string{"karate", "tennis", "programming"}

	err := h.App.WriteXML(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// DownloadFile is the handler to demonstrate file download reponses
func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	h.App.DownloadFile(w, r, "./public/images", "ocean.jpg")
}

func (h *Handlers) TestCrypto(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello, world"
	fmt.Fprint(w, "Unencrypted: "+plainText+"\n")

	enc := goravel.Encryption{Key: []byte(h.App.EncryptionKey)}
	encrypted, err := enc.Encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w, r)
		return
	}

	fmt.Fprint(w, "Encrypted: "+encrypted+"\n")

	decrypted, err := enc.Decrypt(encrypted)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Error500(w, r)
		return
	}

	fmt.Fprint(w, "Decrypted: "+decrypted+"\n")
}

func (h *Handlers) TestMail(w http.ResponseWriter, r *http.Request) {
	msg := mailer.Message{
		From:        "info@goravel.com",
		To:          "saalik@gmail.com",
		Subject:     "Test Subject",
		Template:    "test",
		Attachments: nil,
		Data:        nil,
	}

	h.App.Mail.Jobs <- msg
	res := <-h.App.Mail.Results
	if res.Error != nil {
		h.App.ErrorLog.Println(res.Error)
	}

	// To send directly using SMTP or using API(such as SendGrid, Mailgun, etc)
	// err := a.App.Mail.SendSMTPMessage(msg)
	// err := h.App.Mail.SendUsingAPI(msg, "sendgrid")
	// if err != nil {
	// 	a.App.ErrorLog.Println(err)
	// 	return
	// }

	fmt.Fprint(w, "Send mail!")
}
