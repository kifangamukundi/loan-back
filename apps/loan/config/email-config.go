package config

import (
	"log"
)

// MailConfig stores the mail configuration for different environments
var MailConfig = map[string]map[string]string{
	"default": {
		"host":      "localhost",
		"port":      "1025",
		"auth_user": "admin@localhost.dev",
		"auth_pass": "",
	},
	"support": {
		"host":      "smtp.zoho.com",
		"port":      "465",
		"auth_user": "support@glamonwheels.com",
		"auth_pass": "Stje4E08eHjv",
	},
	"accounts": {
		"host":      "smtp.zoho.com",
		"port":      "465",
		"auth_user": "accounts@glamonwheels.com",
		"auth_pass": "veH5AkSN1DJ1",
	},
}

// LoadMailConfig initializes the email configuration
func LoadMailConfig() map[string]map[string]string {
	// Here you can add any additional logic, like validating the loaded values.
	// For now, we just return the pre-loaded configuration.
	if len(MailConfig["support"]["auth_user"]) == 0 || len(MailConfig["accounts"]["auth_user"]) == 0 {
		log.Println("Warning: Missing environment variables for email configurations")
	}

	return MailConfig
}
