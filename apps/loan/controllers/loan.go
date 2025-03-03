package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/transformations"
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/models"

	"github.com/gin-gonic/gin"
	"github.com/jwambugu/mpesa-golang-sdk"
)

type LoanController struct {
	LoanModel     *models.LoanModel
	DisburseModel *models.DisburseModel
	UserModel     *models.UserModel
	OfficerModel  *models.OfficerModel
	AgentModel    *models.AgentModel
	GroupModel    *models.GroupModel
	MemberModel   *models.MemberModel
}

func NewLoanController(loanModel *models.LoanModel, disburseModel *models.DisburseModel, userModel *models.UserModel, officerModel *models.OfficerModel, agentModel *models.AgentModel, groupModel *models.GroupModel, memberModel *models.MemberModel) *LoanController {
	return &LoanController{
		LoanModel:     loanModel,
		DisburseModel: disburseModel,
		UserModel:     userModel,
		OfficerModel:  officerModel,
		AgentModel:    agentModel,
		GroupModel:    groupModel,
		MemberModel:   memberModel,
	}
}

func (ctrl *LoanController) CreateLoanController(c *gin.Context) {
	var req bindings.CreateLoanRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Agent not found"})
		return
	}

	u := user.(models.User)

	agent, err := ctrl.AgentModel.GetAgentByField("user_id", strconv.Itoa(int(u.ID)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	group, err := ctrl.GroupModel.GetGroupByField("id", strconv.Itoa(req.GroupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	member, err := ctrl.MemberModel.GetMemberByField("id", strconv.Itoa(req.MemberID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	description := parameters.SanitizeText(*req.LoanPurpose, false)

	const InterestRate float64 = 10.0

	newLoan := models.Loan{
		AgentID:      agent.ID,
		Amount:       req.Amount,
		Interest:     InterestRate,
		Term:         req.Term,
		LoanPurpose:  &description,
		DefaultImage: req.DefaultImage,
		Images:       req.Images,
		GroupID:      group.ID,
		MemberID:     member.ID,
	}

	if err := ctrl.LoanModel.CreateLoan(&newLoan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *LoanController) GetAgentGroupMemberLoansController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decodedUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	u := decodedUser.(models.User)

	user, err := ctrl.UserModel.GetUserByFieldPreloaded("id", strconv.Itoa(int(u.ID)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	groupId, valid := parameters.ConvertParamToValidID(c, "groupId")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := ctrl.GroupModel.GetGroupByField("id", string(groupId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	memberId, valid := parameters.ConvertParamToValidID(c, "memberId")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	member, err := ctrl.MemberModel.GetMemberByField("id", string(memberId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	loans, totalCount, count, err := ctrl.LoanModel.GetAgentMemberLoans(int(user.Agent.ID), int(group.ID), int(member.ID), skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching loans: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedLoans := transformations.Transform(loans, fieldNames,
		func(loan models.Loan) interface{} { return loan.ID },
		func(loan models.Loan) interface{} { return *loan.LoanPurpose },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedLoans)
}

func (ctrl *LoanController) GetLoansController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loans, totalCount, count, err := ctrl.LoanModel.GetLoans(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching loans: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"amount",
		"term",
		"created",
		"updated",
		"agent_first_name",
		"agent_last_name",
		"member_first_name",
		"member_last_name",
	}

	transformedLoans := transformations.Transform(loans, fieldNames,
		func(loan models.Loan) interface{} { return loan.ID },
		func(loan models.Loan) interface{} { return loan.Amount },
		func(loan models.Loan) interface{} { return loan.Term },
		func(loan models.Loan) interface{} { return loan.CreatedAt },
		func(loan models.Loan) interface{} { return loan.UpdatedAt },
		func(loan models.Loan) interface{} { return loan.Agent.User.FirstName },
		func(loan models.Loan) interface{} { return loan.Agent.User.LastName },
		func(loan models.Loan) interface{} { return loan.Member.User.FirstName },
		func(loan models.Loan) interface{} { return loan.Member.User.LastName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedLoans)
}

func (ctrl *LoanController) GetLoanByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	loan, err := ctrl.LoanModel.GetLoanByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	agent, err := ctrl.AgentModel.GetAgentByFieldPreloaded("id", fmt.Sprintf("%d", loan.AgentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	group, err := ctrl.GroupModel.GetGroupByField("id", fmt.Sprintf("%d", loan.GroupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	member, err := ctrl.MemberModel.GetMemberByFieldPreloaded("id", fmt.Sprintf("%d", loan.MemberID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	response := bindings.LoanResponse{
		ID:               loan.ID,
		Amount:           loan.Amount,
		Interest:         loan.Interest,
		Term:             loan.Term,
		DefaultImage:     loan.DefaultImage,
		Images:           loan.Images,
		LoanPurpose:      loan.LoanPurpose,
		Status:           loan.Status,
		DueDate:          loan.DueDate,
		LastPaymentDate:  loan.LastPaymentDate,
		RemainingBalance: loan.RemainingBalance,
		AgentFirstName:   agent.User.FirstName,
		AgentLastName:    agent.User.LastName,
		GroupName:        group.GroupName,
		MemberFirstName:  member.User.FirstName,
		MemberLastName:   member.User.LastName,
		CreatedAt:        loan.CreatedAt,
		UpdatedAt:        loan.UpdatedAt,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}

func (ctrl *LoanController) ApproveLoanController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	loan, err := ctrl.LoanModel.GetLoanByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	decodedUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	u := decodedUser.(models.User)

	officer, err := ctrl.OfficerModel.GetOfficerByFieldPreloaded("user_id", fmt.Sprintf("%d", u.ID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Officer not found"})
		return
	}

	member, err := ctrl.UserModel.GetUserByFieldPreloaded("id", fmt.Sprintf("%d", loan.Member.UserID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}
	mobileNumber, err := strconv.ParseUint(member.MobileNumber, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	if loan.Status == "approved" {
		c.JSON(http.StatusConflict, gin.H{"error": "Loan has already been approved"})
		return
	} else if loan.Status == "rejected" {
		c.JSON(http.StatusConflict, gin.H{"error": "Loan has already been rejected"})
		return
	}

	// dueDate := now.AddDate(0, 0, loan.Term) do this after disbursement
	// loan.DueDate = &dueDate

	consumerKey := os.Getenv("MPESA_CONSUMER_KEY")
	consumerSecret := os.Getenv("MPESA_CONSUMER_SECRET")
	initiatorPassword := os.Getenv("MPESA_INITIATOR_PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesa.EnvironmentSandbox)

	b2cResp, err := mpesaApp.B2C(ctx, initiatorPassword, mpesa.B2CRequest{
		InitiatorName: "testapi",
		// SalaryPaymentCommandID, BusinessPaymentCommandID, PromotionPaymentCommandID
		CommandID:       mpesa.BusinessPaymentCommandID,
		Amount:          uint(loan.Amount),
		PartyA:          600999,
		PartyB:          mobileNumber,
		QueueTimeOutURL: "https://your-server.com/api/v1/mpesa/b2c-timeout",
		ResultURL:       "https://your-server.com/api/v1/mpesa/b2c-result",
		Remarks:         "Loan disbursement",
		Occasion:        "Loan",
	})

	if err != nil {
		log.Printf("B2C Payment Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Loan disbursement failed. Check logs."})
		return
	}

	if b2cResp.ResponseCode != "0" {
		log.Printf("M-Pesa Error: %s - %s\n", b2cResp.ResponseCode, b2cResp.ResponseDescription)
		c.JSON(http.StatusBadGateway, gin.H{"error": "M-Pesa transaction failed", "message": b2cResp.ResponseDescription})
		return
	}

	disbursement := models.Disbursement{
		LoanID:                   loan.ID,
		OfficerID:                &officer.ID,
		OriginatorConversationID: b2cResp.OriginatorConversationID,
		ConversationID:           b2cResp.ConversationID,
		ResponseCode:             b2cResp.ResponseCode,
		ResponseDesc:             b2cResp.ResponseDescription,
		Status:                   "pending",
		CreatedAt:                time.Now(),
	}

	if err := ctrl.DisburseModel.CreateDisbursement(&disbursement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = ctrl.LoanModel.ApproveLoan(loan.ID, "approved")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating loan: " + err.Error()})
		return
	}

	response := struct {
		ConversationID      string `json:"ConversationID"`
		ResponseCode        string `json:"ResponseCode"`
		ResponseDescription string `json:"ResponseDescription"`
	}{
		ConversationID:      b2cResp.ConversationID,
		ResponseCode:        b2cResp.ResponseCode,
		ResponseDescription: b2cResp.ResponseDescription,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}

func (ctrl *LoanController) RejectLoanController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	loan, err := ctrl.LoanModel.GetLoanByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	decodedUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	u := decodedUser.(models.User)

	officer, err := ctrl.OfficerModel.GetOfficerByFieldPreloaded("user_id", fmt.Sprintf("%d", u.ID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Officer not found"})
		return
	}

	if loan.Status == "approved" {
		c.JSON(http.StatusConflict, gin.H{"error": "Loan has already been approved"})
		return
	} else if loan.Status == "rejected" {
		c.JSON(http.StatusConflict, gin.H{"error": "Loan has already been rejected"})
		return
	}

	_, err = ctrl.LoanModel.RejectLoan(loan.ID, officer.ID, "rejected")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating loan: " + err.Error()})
		return
	}

	response := struct {
		ID     uint    `json:"ID"`
		Amount float64 `json:"Amount"`
	}{
		ID:     loan.ID,
		Amount: loan.Amount,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}
