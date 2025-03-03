package models

import "time"

type Payment struct {
	ID                uint      `gorm:"primaryKey"`
	LoanID            uint      `gorm:"index"` // Foreign key to Loan
	Loan              Loan      `gorm:"foreignKey:LoanID;constraint:onDelete:CASCADE"`
	Amount            float64   `gorm:"not null"`                   // Payment amount
	PhoneNumber       string    `gorm:"not null"`                   // Customer phone number
	CheckoutRequestID string    `gorm:"uniqueIndex;not null"`       // M-Pesa STK Push ID
	MerchantRequestID string    `gorm:"index"`                      // Merchant request ID from M-Pesa
	Status            string    `gorm:"not null;default:'Pending'"` // Payment status (Pending, Success, Failed)
	ResponseCode      string    `gorm:"not null"`                   // Response code from M-Pesa
	ResponseDesc      string    `gorm:"not null"`                   // Response description
	TransactionDesc   string    `gorm:"not null"`                   // Description of the transaction
	PaymentMode       string    `gorm:"not null"`                   // Payment method (e.g., 'mpesa', 'bank', 'cash')
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}
