package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/karman-dev-team/xero-transaction-bq-uploader/models"
)

func uploadInvoices(invoices []models.XeroTransaction, company string) error {
	bqInvoices, err := convertToBQInvoice(invoices, company)
	if err != nil {
		return err
	}
	batchSize := 1000
	batches := splitIntoBatches(bqInvoices, batchSize)
	err = uploadToBQ(batches)
	if err != nil {
		return err
	}
	return nil
}

func uploadToBQ(batches [][]models.BQTransaction) error {
	maxRetries := 10
	retryInterval := 5 * time.Second
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "reporting-393509")
	if err != nil {
		log.Printf("Failed to create BigQuery client: %v", err)
		return err
	}
	defer client.Close()

	dataset := client.Dataset("internal_reporting")
	table := dataset.Table("xero_transactions")
	uploader := table.Uploader()

	for i, batch := range batches {
		var retryCount int
		for retryCount < maxRetries {
			err := uploader.Put(ctx, batch)
			if err == nil {
				fmt.Printf("Uploaded Batch %d...\n", i+1)
				break
			}
			retryCount++
			log.Printf("Failed to insert data for Batch %d: %v. Retrying", i+1, err)
			if retryCount < maxRetries {
				time.Sleep(retryInterval)
			}
		}
		if retryCount == maxRetries {
			fmt.Printf("Exceeded maximum retries for Batch %d, giving up.\n", i+1)
		}
	}
	return nil
}

func splitIntoBatches(slice []models.BQTransaction, batchSize int) [][]models.BQTransaction {
	var batches [][]models.BQTransaction

	for batchSize < len(slice) {
		slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
	}
	batches = append(batches, slice)

	return batches
}

func convertToBQInvoice(transactions []models.XeroTransaction, company string) ([]models.BQTransaction, error) {
	bqTransactions := []models.BQTransaction{}
	for _, transaction := range transactions {
		date, err := time.Parse("2006-01-02TT15:04:05", transaction.DateString)
		if err != nil {
			return nil, err
		}
		bqTransaction := models.BQTransaction{
			TransactionID: transaction.BankTransactionID,
			Company:       company,
			Date:          date,
			Amount:        transaction.Total,
			Reference:     transaction.Reference,
			Description:   transaction.LineItems[0].Description,
		}
		bqTransactions = append(bqTransactions, bqTransaction)
	}
	return bqTransactions, nil
}
