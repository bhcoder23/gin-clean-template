package v1

import (
	"net/http"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/gin-gonic/gin"
)

// @Summary     Register
// @Description Register a new user
// @ID          register
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     request.RegisterReq true "Registration data"
// @Success     201     {object} response.UserResp
// @Failure     400     {object} response.Error
// @Failure     409     {object} response.Error
// @Failure     500     {object} response.Error
// @Router      /auth/register [post]
func (r *V1) register(ctx *gin.Context) {
	var body request.RegisterReq

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apperror.Log(r.l, err, "restapi - v1 - register")
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		apperror.Log(r.l, err, "restapi - v1 - register")
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return
	}

	user, err := r.u.Register(ctx.Request.Context(), body.Username, body.Email, body.Password)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - register")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, response.NewUserResp(user))
}

// @Summary     Login
// @Description Authenticate user and get JWT token
// @ID          login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     request.LoginReq true "Login credentials"
// @Success     200     {object} response.TokenResp
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Router      /auth/login [post]
func (r *V1) login(ctx *gin.Context) {
	var body request.LoginReq

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apperror.Log(r.l, err, "restapi - v1 - login")
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		apperror.Log(r.l, err, "restapi - v1 - login")
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return
	}

	token, err := r.u.Login(ctx.Request.Context(), body.Email, body.Password)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - login")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, response.TokenResp{Token: token})
}

// @Summary     Get profile
// @Description Get current user profile
// @ID          profile
// @Tags        user
// @Produce     json
// @Success     200 {object} response.UserResp
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /user/profile [get]
func (r *V1) profile(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	user, err := r.u.GetUser(ctx.Request.Context(), userID)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - profile")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, response.NewUserResp(user))
}
