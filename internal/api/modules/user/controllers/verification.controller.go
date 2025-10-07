package users

import (
	"net/http"

	"github.com/deva-labs/univia/internal/api/modules/user/models"
	verificationservices "github.com/deva-labs/univia/internal/api/modules/user/services"
	"github.com/deva-labs/univia/pkg/types"
	_ "github.com/deva-labs/univia/pkg/types"
	"github.com/deva-labs/univia/pkg/utils"

	"github.com/gin-gonic/gin"
)

// VerifyCode godoc
// @Summary Verify the code
// @Description handles the verification code process
// @Tags Authentication
// @Accept       json
// @Produce      json
// @Param request body	types.VerifyCodeRequest true "Verify Code"
// @Success 200 {object} types.SuccessVerifyCodeResponse "Tokens response"
// @Failure 400 {object} types.StatusBadRequest "Invalid Input"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/auth/confirm-forgot-password [post]
// VerifyCode handles the verification code process and token generation
func VerifyCode(c *gin.Context) {
	var code users.VerificationCode

	// Parse JSON input
	if err := utils.BindJson(c, &code); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Call the service to verify the code and generate tokens
	err := verificationservices.VerifyCode(code.Code, code.Email)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: ", err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Verification successful", nil)
}

// VerifyCodeAndGenerateToken godoc
// @Summary Verify the code and generate token
// @Description handles the verification code process and token generation
// @Tags Authentication
// @Accept       json
// @Produce      json
// @Param request body	types.VerifyCodeRequest true "Verify Code"
// @Success 200 {object} types.SuccessVerifyCodeResponse "Tokens response"
// @Failure 400 {object} types.StatusBadRequest "Invalid Input"
// @Failure 500 {object} types.StatusInternalError "Internal server error"
// @Router /api/v1/auth/verification [post]
// VerifyCode handles the verification code process and token generation
// VerifyCodeAndGenerateToken handles the verification code process and token generation
func VerifyCodeAndGenerateToken(c *gin.Context) {
	var code users.VerificationCode

	if err := utils.BindJson(c, &code); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	meta := types.SessionMetadata{
		IP:        ip,
		UserAgent: userAgent,
	}

	response, status, err := verificationservices.VerifyCodeAndGenerateTokens(code, meta)
	if err != nil {
		utils.SendErrorResponse(c, status, err.Error(), nil)
		return
	}
	utils.SetHttpOnlyCookieForSession(c, response.SessionID)

	utils.SendSuccessResponse(c, http.StatusOK, "Verification successful", response)
}
