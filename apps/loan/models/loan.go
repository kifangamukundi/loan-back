package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/deserializers"
	"github.com/kifangamukundi/gm/loan/services"
)

type Loan struct {
	ID           uint                            `gorm:"primaryKey"`
	Amount       float64                         `gorm:"not null"`
	Interest     float64                         `gorm:"not null"`
	Term         int                             `gorm:"not null"`
	DefaultImage deserializers.DefaultImageSlice `json:"DefaultImage" gorm:"type:jsonb"`
	Images       deserializers.DefaultImageSlice `json:"Images" gorm:"type:jsonb;serializer:json"`

	// Loan Approval & Disbursement
	Status      string     `gorm:"not null;default:'pending'"`
	OfficerID   *uint      `gorm:"index;default:null"`
	Officer     *Officer   `gorm:"foreignKey:OfficerID;constraint:onDelete:SET NULL"`
	ApprovedAt  *time.Time `gorm:"default:null"`
	RejectedAt  *time.Time `gorm:"default:null"`
	DisbursedAt *time.Time `gorm:"default:null"`

	// Loan Repayment Tracking
	RemainingBalance float64    `gorm:"not null;default:0"`
	IsFullyPaid      bool       `gorm:"default:false"`
	DueDate          *time.Time `gorm:"default:null"`
	LastPaymentDate  *time.Time `gorm:"default:null"`

	// Borrower Details
	MemberID uint   `gorm:"index"`
	Member   Member `gorm:"foreignKey:MemberID;constraint:onDelete:CASCADE"`

	// others
	GroupID uint  `gorm:"index"`
	Group   Group `gorm:"foreignKey:GroupID;constraint:onDelete:CASCADE"`

	// Iniateed by who
	AgentID uint  `gorm:"index"`
	Agent   Agent `gorm:"foreignKey:AgentID;constraint:onDelete:CASCADE"`

	// Other Loan Information
	LoanPurpose *string   `gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type LoanModel struct {
	Service services.Service
}

func NewLoanModel(service services.Service) *LoanModel {
	return &LoanModel{Service: service}
}

func (m *LoanModel) CreateLoan(loan *Loan) error {
	if err := m.Service.CreateEntity(loan); err != nil {
		return fmt.Errorf("failed to create loan: %v", err)
	}

	return nil
}

func (m *LoanModel) GetAgentMemberLoans(agentId, groupId, memberId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Loan, int64, int64, error) {
	searchColumns := []string{"user.first_name", "user.last_name"}

	preloads := []string{}

	loansResult, totalCount, filteredCount, err := m.Service.GetEntitiesFilteredAgentMemberLoans(&Loan{}, agentId, groupId, memberId, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get loans: %v", err)
	}

	var loans []Loan
	for _, loan := range loansResult {
		if c, ok := loan.(*Loan); ok {
			loans = append(loans, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", loan)
		}
	}

	return loans, totalCount, filteredCount, nil
}

func (m *LoanModel) GetLoans(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Loan, int64, int64, error) {
	searchColumns := []string{"loan_purpose"}

	preloads := []string{
		"Member.User",
		"Agent.User",
		"Group",
	}

	loansResult, totalCount, filteredCount, err := m.Service.GetEntitiesFilteredLoans(&Loan{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get loans: %v", err)
	}

	var loans []Loan
	for _, country := range loansResult {
		if c, ok := country.(*Loan); ok {
			loans = append(loans, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", country)
		}
	}

	return loans, totalCount, filteredCount, nil
}

func (m *LoanModel) GetLoanByFieldPreloaded(field, value string) (*Loan, error) {
	var loan Loan

	preloads := []string{"Agent", "Group", "Member"}

	result, err := m.Service.GetEntityByFieldWithPreload(&loan, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching loan with user info by %s: %v", field, err)
		return nil, err
	}

	loanPtr, ok := result.(*Loan)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return loanPtr, nil
}

func (m *LoanModel) ApproveLoan(id uint, status string) (Loan, error) {
	loan := &Loan{ID: uint(id)}

	_, err := m.Service.GetEntityByID(loan, uint(id))
	if err != nil {
		return Loan{}, fmt.Errorf("loan not found: %v", err)
	}

	now := time.Now()

	loan.Status = status
	loan.ApprovedAt = &now

	if err := m.Service.UpdateEntity(loan); err != nil {
		return Loan{}, fmt.Errorf("failed to update loan: %v", err)
	}

	return *loan, nil
}

func (m *LoanModel) RejectLoan(id, officerId uint, status string) (Loan, error) {
	loan := &Loan{ID: uint(id)}

	_, err := m.Service.GetEntityByID(loan, uint(id))
	if err != nil {
		return Loan{}, fmt.Errorf("loan not found: %v", err)
	}

	now := time.Now()

	loan.Status = status
	loan.RejectedAt = &now
	loan.OfficerID = &officerId

	if err := m.Service.UpdateEntity(loan); err != nil {
		return Loan{}, fmt.Errorf("failed to update loan: %v", err)
	}

	return *loan, nil
}
