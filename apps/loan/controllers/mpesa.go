package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jwambugu/mpesa-golang-sdk"
)

func PaymentController(c *gin.Context) {
	consumerKey := "YOUR_CONSUMER_KEY"
	consumerSecret := "YOUR_CONSUMER_SECRET"
	passkey := "YOUR_PASSKEY"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize Mpesa App
	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesa.EnvironmentSandbox)

	// Perform STK Push
	stkResp, err := mpesaApp.STKPush(ctx, passkey, mpesa.STKPushRequest{
		BusinessShortCode: 174379,
		TransactionType:   mpesa.CustomerPayBillOnlineTransactionType,
		Amount:            1,
		PartyA:            254702817040,
		PartyB:            174379,
		PhoneNumber:       254702817040,
		CallBackURL:       "https://your-server.com/api/v1/mpesa/stk-callback", // Webhook
		AccountReference:  "TestRef",
		TransactionDesc:   "Test Payment",
	})

	if err != nil {
		log.Printf("STK Push Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "STK Push failed. Check logs."})
		return
	}

	log.Printf("STK Push Success: %+v\n", stkResp)

	// Store CheckoutRequestID in DB for tracking (to query later if needed)
	// db.SavePaymentStatus(stkResp.CheckoutRequestID, "Pending")

	// Return CheckoutRequestID to client
	c.JSON(http.StatusOK, gin.H{
		"message":             "STK Push initiated successfully",
		"CheckoutRequestID":   stkResp.CheckoutRequestID,
		"ResponseCode":        stkResp.ResponseCode,
		"ResponseDescription": stkResp.ResponseDescription,
	})
}

func PaymentCallbackController(c *gin.Context) {
	// Decode the callback body using UnmarshalSTKPushCallback
	callback, err := mpesa.UnmarshalSTKPushCallback(c.Request.Body)
	if err != nil {
		log.Printf("Error decoding STK Callback: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("STK Callback Received: %+v\n", callback)

	// Extract transaction details
	resultCode := callback.Body.STKCallback.ResultCode
	resultDesc := callback.Body.STKCallback.ResultDesc
	merchantRequestID := callback.Body.STKCallback.MerchantRequestID
	checkoutRequestID := callback.Body.STKCallback.CheckoutRequestID

	// TODO: Update DB based on ResultCode
	// Example:
	// if resultCode == 0 {
	//     db.UpdateTransactionStatus(checkoutRequestID, "Success")
	// } else {
	//     db.UpdateTransactionStatus(checkoutRequestID, "Failed")
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":           "STK Callback processed",
		"MerchantRequestID": merchantRequestID,
		"CheckoutRequestID": checkoutRequestID,
		"ResultCode":        resultCode,
		"ResultDesc":        resultDesc,
	})
}

func DisburseController(c *gin.Context) {
	consumerKey := "YOUR_CONSUMER_KEY"
	consumerSecret := "YOUR_CONSUMER_SECRET"
	initiatorPassword := "YOUR_INITIATOR_PASSWORD"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize Mpesa Client
	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesa.EnvironmentSandbox)

	// Perform B2C Payment
	b2cResp, err := mpesaApp.B2C(ctx, initiatorPassword, mpesa.B2CRequest{
		InitiatorName:   "TestInitiator",
		Amount:          1000,         // Loan amount
		PartyA:          600000,       // Your paybill or shortcode
		PartyB:          254712345678, // Borrower's phone number
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

	log.Printf("B2C Payment Success: %+v\n", b2cResp)

	c.JSON(http.StatusOK, gin.H{
		"message":             "Loan sent successfully",
		"ConversationID":      b2cResp.ConversationID,
		"ResponseCode":        b2cResp.ResponseCode,
		"ResponseDescription": b2cResp.ResponseDescription,
	})
}

func DisburseCallbackController(c *gin.Context) {
	// Decode the callback body using UnmarshalCallback
	callback, err := mpesa.UnmarshalCallback(c.Request.Body)
	if err != nil {
		log.Printf("Error decoding B2C Callback: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("B2C Callback Received: %+v\n", callback)

	// Extract transaction details
	resultCode := callback.Result.ResultCode
	resultDesc := callback.Result.ResultDesc
	originatorConversationID := callback.Result.OriginatorConversationID
	conversationID := callback.Result.ConversationID
	transactionID := callback.Result.TransactionID

	// TODO: Update DB based on ResultCode
	// Example:
	// if resultCode == 0 {
	//     db.UpdateLoanStatus(originatorConversationID, "Success", transactionID)
	// } else {
	//     db.UpdateLoanStatus(originatorConversationID, "Failed", "")
	// }

	// Respond to Safaricom confirming receipt
	c.JSON(http.StatusOK, gin.H{
		"message":                  "B2C Callback processed",
		"OriginatorConversationID": originatorConversationID,
		"ConversationID":           conversationID,
		"TransactionID":            transactionID,
		"ResultCode":               resultCode,
		"ResultDesc":               resultDesc,
	})
}

// func CreateCountryController(c *gin.Context) {
// 	consumerKey := "MYLQBVtEgjAzqV6AFAoxzjtywD1ebF8A58rbDj5N2wgGnzya"
// 	consumerSecret := "8dIKV7RVk0o9UtPO4InGqG8gQwmNclcbpSHkpXUQU2F15i5tAYspcdGjKXCsiZMt"
// 	passkey := "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout
// 	defer cancel()

// 	// Create Mpesa App
// 	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesa.EnvironmentSandbox)

// 	// Perform STK Push
// 	stkResp, err := mpesaApp.STKPush(ctx, passkey, mpesa.STKPushRequest{
// 		BusinessShortCode: 174379,
// 		TransactionType:   mpesa.CustomerPayBillOnlineTransactionType,
// 		Amount:            1,
// 		PartyA:            254702817040,
// 		PartyB:            174379,
// 		PhoneNumber:       254702817040,
// 		CallBackURL:       "https://your-server.com/stk-callback",
// 		AccountReference:  "TestRef",
// 		TransactionDesc:   "Test Payment",
// 	})

// 	if err != nil {
// 		log.Printf("STK Push Error: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "STK Push failed. Check logs."})
// 		return
// 	}

// 	log.Printf("STK Push Success: %+v\n", stkResp)

// 	// 🕒 Wait for the user to complete the transaction
// 	time.Sleep(10 * time.Second) // Wait 10 seconds before querying

// 	// Perform STK Query
// 	stkQueryRes, err := mpesaApp.STKQuery(ctx, passkey, mpesa.STKQueryRequest{
// 		BusinessShortCode: 174379,
// 		CheckoutRequestID: stkResp.CheckoutRequestID,
// 	})

// 	if err != nil {
// 		log.Printf("STK Query Error: %v\n", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "STK Query failed. Check logs."})
// 		return
// 	}

// 	log.Printf("STK Query Success: %+v\n", stkQueryRes)
// 	c.JSON(http.StatusOK, stkQueryRes)

// }
