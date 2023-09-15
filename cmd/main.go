package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/karman-dev-team/xero-transaction-bq-uploader/models"
	"golang.org/x/oauth2"
)

var (
	oauth2Config = oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
		Scopes: []string{"offline_access accounting.transactions accounting.settings"},
	}
)

var App models.App

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	oauth2Config.ClientID = os.Getenv("CLIENT_ID")
	oauth2Config.ClientSecret = os.Getenv("CLIENT_SECRET")
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/connect", handleConnect)
	http.HandleFunc("/import", handleImport)
	http.ListenAndServe(":8080", nil)
}
