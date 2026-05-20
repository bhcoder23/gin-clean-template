package v1

import (
	"net/http"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/gin-gonic/gin"
)

func (r *V1) bindJSON(ctx *gin.Context, target any, op string) bool {
	if err := ctx.ShouldBindJSON(target); err != nil {
		apperror.Log(r.l, err, op)
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return false
	}

	if err := r.v.Struct(target); err != nil {
		apperror.Log(r.l, err, op)
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid request body")

		return false
	}

	return true
}
