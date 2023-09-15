package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/karman-dev-team/xero-transaction-bq-uploader/models"
	"golang.org/x/oauth2"
)

func getAllTransactions(token *oauth2.Token, tenantID string) ([]models.XeroTransaction, error) {
	transactions := []models.XeroTransaction{}
	page := 1
	for {
		transaction := models.TransactionBody{}
		transactionBytes, err := getTransactions(token, page, tenantID)
		if err != nil {
			fmt.Println("Error getting invoices", err)
			log.Fatal(err)
		}
		err = json.Unmarshal(transactionBytes, &transaction)
		if err != nil {
			fmt.Println("Error unmarshalling", err)
			log.Fatal(err)
		}
		transactions = append(transactions, transaction.BankTransactions...)
		if len(transaction.BankTransactions) < 100 {
			break
		}
		page++
		if page%59 == 0 {
			time.Sleep(60 * time.Second)
		}
	}
	return transactions, nil
}

func getTransactions(token *oauth2.Token, page int, tenantID string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/BankTransactions", nil)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("order", "Date DESC")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("xero-tenant-id", tenantID)
	req.Header.Add("Accept", "application/json")
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("body: %s\n", body)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return body, nil
}
