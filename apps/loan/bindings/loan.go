package bindings

import (
	"time"

	"github.com/kifangamukundi/gm/loan/deserializers"
)

type CreateLoanRequest struct {
	Amount       float64                      `json:"Amount" binding:"required"`
	Term         int                          `json:"Term" binding:"required"`
	DefaultImage []deserializers.DefaultImage `json:"DefaultImage"`
	Images       []deserializers.DefaultImage `json:"Images"`
	LoanPurpose  *string                      `json:"LoanPurpose" binding:"required,min=10"`
	GroupID      int                          `json:"GroupID" binding:"required"`
	MemberID     int                          `json:"MemberID" binding:"required"`
}

type LoanResponse struct {
	ID               uint                         `json:"ID"`
	Amount           float64                      `json:"Amount"`
	Interest         float64                      `json:"Interest"`
	Term             int                          `json:"Term"`
	DefaultImage     []deserializers.DefaultImage `json:"DefaultImage"`
	Images           []deserializers.DefaultImage `json:"Images"`
	LoanPurpose      *string                      `json:"LoanPurpose"`
	Status           string                       `json:"Status"`
	DueDate          *time.Time                   `json:"DueDate"`
	LastPaymentDate  *time.Time                   `json:"LastPaymentDate"`
	RemainingBalance float64                      `json:"RemainingBalance"`
	AgentFirstName   string                       `json:"AgentFirstName"`
	AgentLastName    string                       `json:"AgentLastName"`
	GroupName        string                       `json:"GroupName"`
	MemberFirstName  string                       `json:"MemberFirstName"`
	MemberLastName   string                       `json:"MemberLastName"`
	CreatedAt        time.Time                    `json:"CreatedAt"`
	UpdatedAt        time.Time                    `json:"UpdatedAt"`
}
