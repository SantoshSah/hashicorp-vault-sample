package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

// Global variable for hashicorp vault client
var HCVClient *api.Client

// InitHCVault Initializes Hashicorp vault with root token
func InitHCVault(rootToken string) (err error) {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	conf := &api.Config{
		Address:    os.Getenv("HCVAULT_ADDRESS"),
		HttpClient: httpClient,
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return
	}
	HCVClient = client

	HCVClient.SetToken(rootToken)
	return
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Vault
	err = InitHCVault(os.Getenv("HCVAULT_ROOT_TOKEN"))
	if err != nil {
		log.Println(err)
	}
	logicalClient := HCVClient.Logical()

	// UserID for which secrets are needed to be stored in vault
	userID := "UserID"

	// Storage path in vault to store data
	path := fmt.Sprintf("secret/data/%s/aws", userID)

	// Data to store
	data := map[string]interface{}{
		"key": "AWS_KEY",
	}

	inputData := map[string]interface{}{
		"data": data,
	}

	// Write data to vault
	secret, err := logicalClient.Write(path, inputData)
	if err != nil {
		fmt.Println("Error while saving to vault: ", err.Error())
	}
	log.Println("Input data: ", secret.Data)

	// Read data from vault
	outputData, err := logicalClient.Read(path)
	if err != nil {
		fmt.Println("Error while reading from vault: ", err.Error())
	}

	// Marshal & display output data
	b, _ := json.Marshal(outputData.Data)
	fmt.Println("Output Data: ", string(b))
}
