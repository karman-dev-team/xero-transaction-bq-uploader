package models

import (
	"time"

	"golang.org/x/oauth2"
)

type App struct {
	Oauth2Token *oauth2.Token
}

type XeroCompany struct {
	ID      string
	Company string
}

type PageData struct {
	TokenSet bool
}

type TransactionBody struct {
	ID               string            `json:"Id"`
	Status           string            `json:"Status"`
	ProviderName     string            `json:"ProviderName"`
	DateTimeUTC      string            `json:"DateTimeUTC"`
	BankTransactions []XeroTransaction `json:"BankTransactions"`
}

type XeroTransaction struct {
	BankTransactionID string      `json:"BankTransactionID"`
	BankAccount       BankAccount `json:"BankAccount"`
	Type              string      `json:"Type"`
	Reference         string      `json:"Reference"`
	IsReconciled      bool        `json:"IsReconciled"`
	HasAttachments    bool        `json:"HasAttachments"`
	Contact           Contact     `json:"Contact"`
	DateString        string      `json:"DateString"`
	Date              string      `json:"Date"`
	Status            string      `json:"Status"`
	LineAmountTypes   string      `json:"LineAmountTypes"`
	LineItems         []LineItem  `json:"LineItems"`
	SubTotal          float64     `json:"SubTotal"`
	TotalTax          float64     `json:"TotalTax"`
	Total             float64     `json:"Total"`
	UpdatedDateUTC    string      `json:"UpdatedDateUTC"`
	CurrencyCode      string      `json:"CurrencyCode"`
}

type BankAccount struct {
	AccountID string `json:"AccountID"`
	Code      string `json:"Code"`
	Name      string `json:"Name"`
}

type Contact struct {
	ContactID           string `json:"ContactID"`
	Name                string `json:"Name"`
	Addresses           []any  `json:"Addresses"`
	Phones              []any  `json:"Phones"`
	ContactGroups       []any  `json:"ContactGroups"`
	ContactPersons      []any  `json:"ContactPersons"`
	HasValidationErrors bool   `json:"HasValidationErrors"`
}

type LineItem struct {
	Description string  `json:"Description"`
	UnitAmount  float64 `json:"UnitAmount"`
	TaxType     string  `json:"TaxType"`
	TaxAmount   float64 `json:"TaxAmount"`
	LineAmount  float64 `json:"LineAmount"`
	AccountCode string  `json:"AccountCode"`
	Tracking    []any   `json:"Tracking"`
	Quantity    float64 `json:"Quantity"`
	LineItemID  string  `json:"LineItemID"`
	AccountID   string  `json:"AccountID"`
}

type BQTransaction struct {
	TransactionID string    `bigquery:"id"`
	Company       string    `bigquery:"company"`
	Date          time.Time `bigquery:"date"`
	Amount        float64   `bigquery:"amount"`
	Reference     string    `bigquery:"reference"`
	RevenueLine   string    `bigquery:"revenue_line"`
	Description   string    `bigquery:"description"`
	Group         string    `bigquery:"transfer_group"`
	AccountCode   string    `bigquery:"account_code"`
}

type AccountBody struct {
	Account []Account `json:"Accounts"`
}

type Account struct {
	AccountID               string `json:"AccountID"`
	Code                    string `json:"Code"`
	Name                    string `json:"Name"`
	Type                    string `json:"Type"`
	TaxType                 string `json:"TaxType"`
	EnablePaymentsToAccount bool   `json:"EnablePaymentsToAccount"`
	BankAccountNumber       string `json:"BankAccountNumber"`
	BankAccountType         string `json:"BankAccountType"`
	CurrencyCode            string `json:"CurrencyCode"`
}

type AccountLookup struct {
	Name  string
	Group string
}

type Journal struct {
	JournalID      string        `json:"JournalID"`
	JournalDate    string        `json:"JournalDate"`
	JournalNumber  int           `json:"JournalNumber"`
	CreatedDateUTC string        `json:"CreatedDateUTC"`
	Reference      string        `json:"Reference"`
	JournalLines   []JournalLine `json:"JournalLines"`
}

type JournalLine struct {
	JournalLineID      string         `json:"JournalLineID"`
	AccountID          string         `json:"AccountID"`
	AccountCode        string         `json:"AccountCode"`
	AccountType        string         `json:"AccountType"`
	AccountName        string         `json:"AccountName"`
	Description        string         `json:"Description"`
	NetAmount          float64        `json:"NetAmount"`
	GrossAmount        float64        `json:"GrossAmount"`
	TaxAmount          float64        `json:"TaxAmount"`
	TaxType            string         `json:"TaxType"`
	TaxName            string         `json:"TaxName"`
	TrackingCategories []TrackingItem `json:"TrackingCategories"`
}

type TrackingItem struct {
	TrackingCategoryID string `json:"TrackingCategoryID"`
	TrackingOptionID   string `json:"TrackingOptionID"`
	Name               string `json:"Name"`
	Option             string `json:"Option"`
}

type JournalsResponse struct {
	Journals []Journal `json:"Journals"`
}

type AccountTransaction struct {
	TransactionID string
	Date          time.Time
	Amount        float64
	Reference     string
	AccountCode   string
	Description   string
}
