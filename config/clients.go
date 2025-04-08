package config

import (
	"log"

	"github.com/plutov/paypal/v4"
)

func ConnectPaypal(config *Config) (*paypal.Client, error) {
	c, err := paypal.NewClient(config.PaypalClientID, config.PaypalSecret, paypal.APIBaseSandBox)
	if err != nil {
		log.Fatalf("Failed to create PayPal client: %v", err)
	}
	return c, nil
}
