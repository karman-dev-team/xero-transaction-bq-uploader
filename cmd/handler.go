package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/karman-dev-team/xero-transaction-bq-uploader/models"
)

func handleConnect(w http.ResponseWriter, r *http.Request) {
	authURL := oauth2Config.AuthCodeURL("")
	http.Redirect(w, r, authURL, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	App.Oauth2Token = token
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	pageData := models.PageData{
		TokenSet: false,
	}
	tmpl, err := template.ParseFiles("templates/uploader.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if App.Oauth2Token != nil {
		pageData.TokenSet = true
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleImport(w http.ResponseWriter, r *http.Request) {
	// Simulate an import process
	msg, err := importXeroData()
	response := map[string]string{
		"message": msg,
	}
	if err != nil {
		response["error"] = err.Error()
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func importXeroData() (string, error) {
	tenantID := []models.XeroCompany{
		{ID: os.Getenv("CF_TENANT_ID"), Company: "CF"},
		{ID: os.Getenv("KD_TENANT_ID"), Company: "KD"}}

	for _, tenant := range tenantID {
		transactions, err := getAllTransactions(App.Oauth2Token, tenant.ID)
		if err != nil {
			return "", err
		}
		fmt.Println(transactions)
	}
	return "Success", nil
}
