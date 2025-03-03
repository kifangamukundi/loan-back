package parameters

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

type ValidID string

func (v ValidID) Validate() bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(string(v))
}
func ConvertParamToValidID(c *gin.Context, param string) (ValidID, bool) {
	id := c.Param(param)
	validID := ValidID(id)

	if !validID.Validate() {
		return validID, false
	}

	return validID, true
}

type ValidString string

func (v ValidString) Validate() bool {
	re := regexp.MustCompile("^[A-Za-z0-9_-]+={0,2}$") // Base64 URL-safe regex
	return re.MatchString(string(v))
}

func ConvertParamToValidString(c *gin.Context, param string) (ValidString, bool) {
	paramValue := c.Param(param)
	validParam := ValidString(paramValue)

	if !validParam.Validate() {
		return validParam, false
	}

	return validParam, true
}

type ValidSlug string

func (v ValidSlug) Validate() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$")
	return re.MatchString(string(v))
}

func ConvertParamToValidSlug(c *gin.Context, param string) (ValidSlug, bool) {
	slug := c.Param(param)
	validSlug := ValidSlug(slug)

	if !validSlug.Validate() {
		return validSlug, false
	}

	return validSlug, true
}
