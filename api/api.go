package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/josephmakin/monerohub/models"
)

var Endpoint string
const contentType string = "application/json"

func MakePaymentRequest(callbackURL string) (models.Payment, error) {
	body, err := json.Marshal(models.Payment{
		CallbackURL: callbackURL,
	})
	if err != nil {
		fmt.Println("Error converting body to json...", err.Error())
		return models.Payment{}, err
	}

	response, err := http.Post(Endpoint, contentType, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return models.Payment{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response status %s\n", response.Status)
		return models.Payment{}, err
	}

	var paymentResponse models.Payment
	err = json.NewDecoder(response.Body).Decode(&paymentResponse)
	if err != nil {
		fmt.Println("Could not decode response...")
		return models.Payment{}, err
	}

	return paymentResponse, err
}
