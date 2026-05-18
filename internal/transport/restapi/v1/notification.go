package v1

import (
	"net/http"
	"strconv"

	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errlog"
	"github.com/bhcoder23/gin-clean-template/internal/transport/errmap"
	_ "github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response" // for swaggo
	"github.com/gin-gonic/gin"
)

// @Summary     List notifications
// @Description List notifications for the current user
// @ID          list-notifications
// @Tags        notifications
// @Produce     json
// @Param       unread_only query    bool false "Only unread notifications"
// @Param       limit       query    int  false "Limit"  default(10)
// @Param       offset      query    int  false "Offset" default(0)
// @Success     200 {object} domain.NotificationList
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /notifications [get]
func (r *V1) listNotifications(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	var unreadOnly *bool

	if rawUnread := ctx.Query("unread_only"); rawUnread != "" {
		parsed, err := strconv.ParseBool(rawUnread)
		if err != nil {
			errorResponse(ctx, http.StatusBadRequest, "invalid unread_only")

			return
		}

		unreadOnly = &parsed
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid limit")

		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid offset")

		return
	}

	notifications, total, err := r.n.List(ctx.Request.Context(), userID, unreadOnly, limit, offset)
	if err != nil {
		errlog.Log(r.l, err, "restapi - v1 - listNotifications")
		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, domain.NotificationList{
		Notifications: notifications,
		Total:         total,
	})
}

// @Summary     Mark notification as read
// @Description Mark a notification as read for the current user
// @ID          mark-notification-read
// @Tags        notifications
// @Produce     json
// @Param       id path string true "Notification ID"
// @Success     200 {object} domain.Notification
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /notifications/{id}/read [patch]
func (r *V1) markNotificationRead(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	notification, err := r.n.MarkRead(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		errlog.Log(r.l, err, "restapi - v1 - markNotificationRead")
		statusCode, message := errmap.HTTP(err)
		errorResponse(ctx, statusCode, message)

		return
	}

	ctx.JSON(http.StatusOK, notification)
}
