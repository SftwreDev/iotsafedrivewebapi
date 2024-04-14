package utils

import (
	"fmt"
	"github.com/vonage/vonage-go-sdk"
	"os"
)

func SendSMS(receiver string, content string) error {
	API_KEY := os.Getenv("VONAGE_API_KEY")
	API_SECRET := os.Getenv("VONAGE_API_SECRET")

	auth := vonage.CreateAuthFromKeySecret(API_KEY, API_SECRET)
	smsClient := vonage.NewSMSClient(auth)
	response, errResp, err := smsClient.Send(
		"Vonage APIs",
		receiver,
		content,
		vonage.SMSOpts{})

	if err != nil {
		panic(err)
	}

	fmt.Println(response)

	if response.Messages[0].Status == "0" {
		fmt.Println("Account Balance: " + response.Messages[0].RemainingBalance)
		return nil
	} else {
		fmt.Println("Error code " + errResp.Messages[0].Status + ": " + errResp.Messages[0].ErrorText)
		return fmt.Errorf("Error code " + errResp.Messages[0].Status + ": " + errResp.Messages[0].ErrorText)
	}

}
