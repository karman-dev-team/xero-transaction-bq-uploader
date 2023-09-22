package main

import (
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/karman-dev-team/xero-transaction-bq-uploader/models"
)

func mergeTransactionsAndJournals(transactions []models.XeroTransaction, journals []models.Journal) ([]models.AccountTransaction, error) {
	journalAcccountTransactions, err := convertJournalsToAccountTransactions(journals)
	if err != nil {
		return nil, err
	}
	transactionAccountTransactions, err := convertTransactionsToAccountTransactions(transactions)
	if err != nil {
		return nil, err
	}
	filteredTransactionAccountTransactions, err := filterBankAccountTransactions(transactionAccountTransactions)
	if err != nil {
		return nil, err
	}
	accountTransactions := append(journalAcccountTransactions, filteredTransactionAccountTransactions...)
	return accountTransactions, nil
}

func convertJournalsToAccountTransactions(journals []models.Journal) ([]models.AccountTransaction, error) {
	accountTransactions := []models.AccountTransaction{}
	for _, journal := range journals {
		for _, journalLine := range journal.JournalLines {
			if journalLine.AccountType == "REVENUE" || journalLine.AccountType == "EXPENSE" || journalLine.AccountType == "OVERHEADS" || journalLine.AccountType == "OTHERINCOME" || journalLine.AccountType == "DIRECTCOSTS" {
				dateStr := journal.JournalDate
				re := regexp.MustCompile(`/Date\((\d+)\+\d+\)/`)
				dateSplit := re.FindStringSubmatch(dateStr)
				var date time.Time
				if len(dateSplit) > 1 {
					dateUnix, err := strconv.Atoi(dateSplit[1])
					if err != nil {
						return nil, err
					}
					date = time.Unix(int64(dateUnix/1000), 0)
				} else {
					date = time.Now()
				}
				accountTransaction := models.AccountTransaction{
					TransactionID: journalLine.JournalLineID,
					AccountCode:   journalLine.AccountCode,
					Date:          date,
					Amount:        math.Abs(journalLine.GrossAmount),
					Reference:     journal.Reference,
					Description:   journalLine.Description,
				}
				accountTransactions = append(accountTransactions, accountTransaction)
			}
		}
	}
	return accountTransactions, nil
}

func convertTransactionsToAccountTransactions(transactions []models.XeroTransaction) ([]models.AccountTransaction, error) {
	accountTransactions := []models.AccountTransaction{}
	for _, transaction := range transactions {
		date, err := time.Parse("2006-01-02T15:04:05", transaction.DateString)
		if err != nil {
			return nil, err
		}
		accountTransaction := models.AccountTransaction{
			TransactionID: transaction.BankTransactionID,
			AccountCode:   transaction.LineItems[0].AccountCode,
			Date:          date,
			Amount:        math.Abs(transaction.Total),
			Reference:     transaction.Reference,
			Description:   transaction.LineItems[0].Description,
		}
		accountTransactions = append(accountTransactions, accountTransaction)
	}
	return accountTransactions, nil
}

func filterBankAccountTransactions(transactions []models.AccountTransaction) ([]models.AccountTransaction, error) {
	filteredTransactions := []models.AccountTransaction{}
	for _, transaction := range transactions {
		if transaction.AccountCode != "7003" {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}
	return filteredTransactions, nil
}
