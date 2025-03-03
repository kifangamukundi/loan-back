package emails

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

var mailConfig = map[string]*gomail.Dialer{}

// InitializeMailConfig initializes mail configurations for different aliases
func InitializeMailConfig(config map[string]map[string]string) {
	for alias, configValues := range config {
		// Log the config values to verify what we are receiving
		log.Printf("Initializing mail config for alias %s with values: %+v", alias, configValues)

		// Convert the port to an integer
		port, err := strconv.Atoi(configValues["port"])
		if err != nil {
			log.Printf("Error: Invalid port number %s for alias %s. Skipping configuration.", configValues["port"], alias)
			continue
		}

		// Initialize gomail.Dialer with host, port, username, and password
		dialer := gomail.NewDialer(configValues["host"], port, configValues["auth_user"], configValues["auth_pass"])

		// Store the mail configuration for the given alias
		mailConfig[alias] = dialer
		log.Printf("Mail configuration for %s initialized.", alias)
	}
}

func SendEmail(options map[string]string, alias string) error {
	// Check if the mail configuration exists for the provided alias
	dialer, exists := mailConfig[alias]
	if !exists {
		return fmt.Errorf("mail configuration for alias %s not found", alias)
	}

	// Set "From" to the dialer's username if not set in options
	from := options["from"]
	if from == "" {
		from = dialer.Username
	}

	// Prepare the email message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", options["to"])
	m.SetHeader("Subject", options["subject"])
	m.SetBody("text/html", options["html"])

	// Send the email using the configured mail server
	err := dialer.DialAndSend(m)
	if err != nil {
		log.Printf("Failed to send email: %s", err.Error())
		return err
	}

	log.Printf("Email sent successfully to %s", options["to"])
	return nil
}
