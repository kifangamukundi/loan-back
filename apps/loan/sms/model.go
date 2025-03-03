package sms

type SMSClient interface {
	SendSMS(to, message string) error
}