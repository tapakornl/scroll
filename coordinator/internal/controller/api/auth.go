package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"scroll-tech/coordinator/internal/types"
)

// AuthController is auth API
type AuthController struct {
}

// NewAuthController returns an AuthController instance
func NewAuthController() *AuthController {
	return &AuthController{}
}

// HealthCheck the api controller for coordinator health check
func (a *AuthController) HealthCheck(c *gin.Context) {
	types.RenderJSON(c, types.Success, nil, nil)
}

// Login the api controller for login
func (a *AuthController) Login(c *gin.Context) (interface{}, error) {
	var login types.LoginParameter
	if err := c.ShouldBind(&login); err != nil {
		return "", fmt.Errorf("missing the public_key, err:%w", err)
	}

	if !a.checkValidPublicKey(&login) {
		return nil, errors.New("incorrect public_key")
	}

	return types.LoginParameter{ProverName: login.ProverName}, nil
}

func (a *AuthController) checkValidPublicKey(param *types.LoginParameter) bool {
	return strings.TrimSpace(param.ProverName) != ""
}

// LoginResponse response login api
func (a *AuthController) LoginResponse(c *gin.Context, code int, message string, time time.Time) {
	resp := types.LoginSchema{
		Time:  time,
		Token: message,
	}
	types.RenderJSON(c, code, nil, resp)
}

// Authorizator validate the token
func (a *AuthController) Authorizator(data interface{}, c *gin.Context) bool {
	tokenCliams, ok := data.(*types.LoginParameter)
	if !ok {
		return false
	}

	if tokenCliams.ProverName == "" {
		return false
	}

	c.Set(types.ProverNameCtxKey, tokenCliams.ProverName)
	return true
}
