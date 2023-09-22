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
		if page%20 == 0 {
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
	params.Add("where", "Status!=\"DELETED\"")
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
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("body: %s\n", body)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return body, nil
}

func getAllJournals(token *oauth2.Token, tenantID string) ([]models.Journal, error) {
	journals := []models.Journal{}
	offset := 0
	for {
		journal := models.JournalsResponse{}
		journalBytes, err := getJournals(token, offset, tenantID)
		if err != nil {
			fmt.Println("Error getting invoices", err)
			log.Fatal(err)
		}
		err = json.Unmarshal(journalBytes, &journal)
		if err != nil {
			fmt.Println("Error unmarshalling", err)
			log.Fatal(err)
		}
		journals = append(journals, journal.Journals...)
		if len(journal.Journals) < 100 {
			break
		}
		offset += 100
		if offset%2000 == 0 {
			time.Sleep(60 * time.Second)
		}
	}
	return journals, nil
}

func getJournals(token *oauth2.Token, offset int, tenantID string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/Journals", nil)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("offset", fmt.Sprintf("%d", offset))
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
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("body: %s\n", body)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return body, nil
}

func modifyAccountLookupTable(accountLookup map[string]models.AccountLookup) map[string]models.AccountLookup {
	delete(accountLookup, "9310")
	delete(accountLookup, "6005")
	replaceLookup := map[string]string{
		"4103": "4201",
		"4106": "4105",
		"4200": "4100",
		"4202": "4102",
		"4207": "4206",
		"4906": "4206",
		"4907": "4201",
		"4908": "4102",
		"4909": "4101",
		"4915": "4102",
		"4916": "4102",
		"4917": "4201",
		"4918": "4100",
		"4923": "4102",
		"4924": "4206",
		"4925": "4206",
		"4926": "4206",
		"7007": "8204",
		"7303": "8204",
		"7399": "7400",
		"7401": "7304",
		"7402": "7304",
		"8209": "8213",
		"8214": "8213",
		"9000": "7003",
		"9100": "7003",
		"9200": "7003",
		"9300": "7003",
		"9505": "4100",
		"9506": "4105",
		"9509": "4100",
		"477":  "7003",
		"478":  "7003",
	}
	for key, value := range replaceLookup {
		tempAccount := models.AccountLookup{
			Name:  accountLookup[value].Name,
			Group: accountLookup[value].Group,
		}
		accountLookup[key] = tempAccount
	}
	return accountLookup
}

func getAccountLookupTable(token *oauth2.Token, tenantID string) (map[string]models.AccountLookup, error) {
	accounts, err := getAccounts(token, tenantID)
	if err != nil {
		return nil, err
	}
	accountLookup := make(map[string]models.AccountLookup)
	for _, account := range accounts.Account {
		tempAccount := models.AccountLookup{
			Name:  account.Name,
			Group: setGroup(account.Type),
		}
		accountLookup[account.Code] = tempAccount
	}
	return accountLookup, nil
}

func getAccounts(token *oauth2.Token, tenantID string) (models.AccountBody, error) {
	accounts := models.AccountBody{}
	req, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/Accounts", nil)
	if err != nil {
		return accounts, err
	}
	params := url.Values{}
	params.Add("where", "Type==\"REVENUE\"||Type==\"EXPENSE\"||Type==\"OVERHEADS\"||Type==\"OTHERINCOME\"||Type==\"DIRECTCOSTS\"")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("xero-tenant-id", tenantID)
	req.Header.Add("Accept", "application/json")
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Do(req)
	if err != nil {
		return accounts, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return accounts, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("body: %s\n", body)
		return accounts, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return accounts, err
	}

	return accounts, nil
}

func setGroup(transactionType string) string {
	var group string
	if transactionType == "EXPENSE" || transactionType == "OVERHEADS" {
		group = "Administration Costs"
	}
	if transactionType == "REVENUE" || transactionType == "OTHERCOSTS" {
		group = "Revenue"
	}
	if transactionType == "DIRECTCOSTS" {
		group = "Cost of Sale"
	}
	return group
}
