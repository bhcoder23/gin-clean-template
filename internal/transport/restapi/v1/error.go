package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/bhcoder23/gin-clean-template/pkg/requestid"
	"github.com/gin-gonic/gin"
)

func errorResponse(ctx *gin.Context, status int, code, msg string) {
	id, _ := requestid.FromContext(ctx.Request.Context())
	ctx.AbortWithStatusJSON(status, response.Error{
		Error: response.ErrorBody{
			Code:      code,
			Message:   msg,
			RequestID: id,
		},
	})
}

func mappedErrorResponse(ctx *gin.Context, err error) {
	mapping := apperror.From(err)
	errorResponse(ctx, mapping.HTTPStatus, mapping.Code, mapping.Message)
}
