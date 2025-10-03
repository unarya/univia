package users

import (
	"net/http"
	Users "univia/src/modules/user/models"
	verificationservices "univia/src/modules/user/services"
	"univia/src/utils"

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
// @Router /api/v1/auth/verification [post]
// VerifyCode handles the verification code process and token generation
func VerifyCode(c *gin.Context) {
	var code Users.VerificationCode

	// Parse JSON input
	if err := utils.BindJson(c, &code); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Call the service to verify the code and generate tokens
	token, err := verificationservices.VerifyCode(code.Code, code.Email)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized: ", err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Verification successful", token)
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
	var code Users.VerificationCode

	// Parse JSON input
	if err := utils.BindJson(c, &code); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Call the service to verify the code and generate tokens
	response, status, err := verificationservices.VerifyCodeAndGenerateTokens(code)
	if err != nil {
		utils.SendErrorResponse(c, status, err.Error(), nil)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Verification successful", response)
}
