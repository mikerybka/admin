package main

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	http.HandleFunc("GET /secrets/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("/home/mike/data", r.URL.Path)
		b, _ := os.ReadFile(path)
		w.Write(b)
	})
	http.HandleFunc("GET /tv", func(w http.ResponseWriter, r *http.Request) {
		w.Write(tvHTML)
	})
	http.ListenAndServe(":2222", nil)
}

//go:embed tv.html
var tvHTML []byte
