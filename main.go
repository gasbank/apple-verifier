package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type RequestBody struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions"`
}

type AppleReceipt struct {
	Receipt struct {
		ReceiptType                string `json:"receipt_type"`
		AdamId                     int    `json:"adam_id"`
		AppItemId                  int    `json:"app_item_id"`
		BundleId                   string `json:"bundle_id"`
		ApplicationVersion         string `json:"application_version"`
		DownloadId                 int64  `json:"download_id"`
		VersionExternalIdentifier  int    `json:"version_external_identifier"`
		ReceiptCreationDate        string `json:"receipt_creation_date"`
		ReceiptCreationDateMs      string `json:"receipt_creation_date_ms"`
		ReceiptCreationDatePst     string `json:"receipt_creation_date_pst"`
		RequestDate                string `json:"request_date"`
		RequestDateMs              string `json:"request_date_ms"`
		RequestDatePst             string `json:"request_date_pst"`
		OriginalPurchaseDate       string `json:"original_purchase_date"`
		OriginalPurchaseDateMs     string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePst    string `json:"original_purchase_date_pst"`
		OriginalApplicationVersion string `json:"original_application_version"`
		InApp                      []struct {
			Quantity                string `json:"quantity"`
			ProductId               string `json:"product_id"`
			TransactionId           string `json:"transaction_id"`
			OriginalTransactionId   string `json:"original_transaction_id"`
			PurchaseDate            string `json:"purchase_date"`
			PurchaseDateMs          string `json:"purchase_date_ms"`
			PurchaseDatePst         string `json:"purchase_date_pst"`
			OriginalPurchaseDate    string `json:"original_purchase_date"`
			OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
			OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
			IsTrialPeriod           string `json:"is_trial_period"`
			InAppOwnershipType      string `json:"in_app_ownership_type"`
		} `json:"in_app"`
	} `json:"receipt"`
	Environment   string `json:"environment"`
	LatestReceipt string `json:"latest_receipt"`
	Status        int    `json:"status"`
}

func main() {
	log.Println("apple-receipt-verifier: Voided Receipt Checker")

	log.Println("apple-receipt-verifier: load .env file")
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	reqBody := RequestBody{
		ExcludeOldTransactions: true,
		Password:               os.Getenv("APPLE_RECEIPT_VERIFIER_SHARED_SECRET"),
		ReceiptData:            os.Getenv("APPLE_RECEIPT_VERIFIER_TEST_RECEIPT_DATA"),
	}

	marshal, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("https://buy.itunes.apple.com/verifyReceipt", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
		var receipt AppleReceipt
		err := json.Unmarshal(respBody, &receipt)
		if err != nil {
			log.Print(err)
		} else {
			err := ioutil.WriteFile(path.Join("results", receipt.Receipt.InApp[0].TransactionId+".json"), respBody, 0644)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
