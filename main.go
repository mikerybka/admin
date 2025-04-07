package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mikerybka/twilio"
	"github.com/mikerybka/util"
)

func main() {
	twilioClient := twilio.NewClient(
		util.RequireEnvVar("TWILIO_ACCOUNT_SID"),
		util.RequireEnvVar("TWILIO_AUTH_TOKEN"),
		util.RequireEnvVar("TWILIO_PHONE_NUMBER"),
	)
	adminPhone := util.RequireEnvVar("ADMIN_PHONE_NUMBER")
	http.HandleFunc("POST /alert", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = twilioClient.SendSMS(adminPhone, strings.TrimSpace(string(b)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("GET /secrets/email", func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("/home/mike/data/secrets/email")
		w.Write(b)
	})
	http.ListenAndServe(":2222", nil)
}
