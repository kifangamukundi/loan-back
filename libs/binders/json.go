package binders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Constants for repeated strings
const (
	SuccessMessage  = "Action successfully"
	ItemKey         = "item"
	ImagesKey       = "images"
	UserKey         = "user"
	AccessTokenKey  = "accessToken"
	RefreshTokenKey = "refreshToken"
	TotalCountKey   = "totalCount"
	CountKey        = "count"
	PageKey         = "page"
	LimitKey        = "limit"
	ItemsKey        = "items"
	OtherKey        = "other"
)

// Generic JSON response handler
func ReturnJSONResponse(c *gin.Context, status int, success bool, data interface{}) {
	c.JSON(status, gin.H{
		"success": success,
		"data":    data,
	})
}

// Helper function for common responses
func returnSuccessResponse(c *gin.Context, status int, data interface{}) {
	ReturnJSONResponse(c, status, true, data)
}

// Paginated response
type PaginatedResponse struct {
	TotalCount int         `json:"totalCount"`
	Count      int         `json:"count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Items      interface{} `json:"items"`
}

func ReturnJSONPaginateResponse(c *gin.Context, page, limit, totalCount, count int, data interface{}) {
	paginatedResponse := PaginatedResponse{
		TotalCount: totalCount,
		Count:      count,
		Page:       page,
		Limit:      limit,
		Items:      data,
	}
	returnSuccessResponse(c, http.StatusOK, paginatedResponse)
}

// General response
func ReturnJSONGeneralResponse(c *gin.Context, data interface{}) {
	returnSuccessResponse(c, http.StatusOK, gin.H{ItemKey: data})
}

// Cache response
func ReturnJSONCacheResponse(c *gin.Context, data interface{}) {
	returnSuccessResponse(c, http.StatusOK, data)
}

// Blog post response
func ReturnJSONBlogPostResponse(c *gin.Context, data, other interface{}) {
	returnSuccessResponse(c, http.StatusOK, gin.H{ItemKey: data, OtherKey: other})
}

// Permissions response
func ReturnJSONPermissionsResponse(c *gin.Context, data interface{}) {
	returnSuccessResponse(c, http.StatusOK, data)
}

// Created generic response
func ReturnJSONCreatedGenericResponse(c *gin.Context) {
	returnSuccessResponse(c, http.StatusCreated, SuccessMessage)
}

// OK generic response
func ReturnJSONOkayGenericResponse(c *gin.Context) {
	returnSuccessResponse(c, http.StatusOK, SuccessMessage)
}

// Token response
type TokenResponse struct {
	User struct {
		AccessToken  interface{} `json:"accessToken"`
		RefreshToken interface{} `json:"refreshToken"`
	} `json:"user"`
}

func ReturnJSONTokenResponse(c *gin.Context, access, refresh interface{}) {
	tokenResponse := TokenResponse{}
	tokenResponse.User.AccessToken = access
	tokenResponse.User.RefreshToken = refresh
	returnSuccessResponse(c, http.StatusCreated, tokenResponse)
}

// Media upload response
func ReturnJSONMediaUploadResponse(c *gin.Context, images interface{}) {
	returnSuccessResponse(c, http.StatusCreated, gin.H{ImagesKey: images})
}
