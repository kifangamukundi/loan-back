package models

import (
	"fmt"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type Disbursement struct {
	ID                       uint       `gorm:"primaryKey"`
	LoanID                   uint       `gorm:"index"` // Foreign key to Loan
	Loan                     Loan       `gorm:"foreignKey:LoanID;constraint:onDelete:CASCADE"`
	OriginatorConversationID string     `gorm:"uniqueIndex;not null"`       // Unique ID for request
	ConversationID           string     `gorm:"index"`                      // M-Pesa conversation ID
	TransactionID            string     `gorm:"index"`                      // M-Pesa transaction ID
	ResponseCode             string     `gorm:"not null"`                   // Response code from M-Pesa
	ResponseDesc             string     `gorm:"not null"`                   // Response description
	Status                   string     `gorm:"not null;default:'pending'"` // pending, processing, completed, failed
	DisbursedAt              *time.Time `gorm:""`                           // Nullable until processed
	OfficerID                *uint      `gorm:"index;default:NULL"`         // Optional (for manual processing)
	Officer                  *Officer   `gorm:"foreignKey:OfficerID;constraint:onDelete:SET NULL"`
	CreatedAt                time.Time  `gorm:"not null"`
	UpdatedAt                time.Time  `gorm:"not null"`
}

type DisburseModel struct {
	Service services.Service
}

func NewDisburseModel(service services.Service) *DisburseModel {
	return &DisburseModel{Service: service}
}

func (m *DisburseModel) CreateDisbursement(disbursement *Disbursement) error {
	if err := m.Service.CreateEntity(disbursement); err != nil {
		return fmt.Errorf("failed to create disbursement: %v", err)
	}

	return nil
}
