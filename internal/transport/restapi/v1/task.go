package v1

import (
	"net/http"
	"strconv"

	"github.com/bhcoder23/gin-clean-template/internal/apperror"
	"github.com/bhcoder23/gin-clean-template/internal/domain"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/transport/restapi/v1/response"
	"github.com/gin-gonic/gin"
)

// @Summary     Create task
// @Description Create a new task for the current user
// @ID          create-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body     request.CreateTaskReq true "Task data"
// @Success     201     {object} response.TaskResp
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [post]
func (r *V1) createTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	var body request.CreateTaskReq

	if ok := r.bindJSON(ctx, &body, "restapi - v1 - createTask"); !ok {
		return
	}

	task, err := r.tk.Create(ctx.Request.Context(), userID, body.Title, body.Description)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - createTask")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, response.NewTaskResp(&task))
}

// @Summary     List tasks
// @Description List tasks for the current user with optional filtering
// @ID          list-tasks
// @Tags        tasks
// @Produce     json
// @Param       status query    string false "Filter by status" Enums(todo, in_progress, done)
// @Param       q      query    string false "Search in task title"
// @Param       limit  query    int    false "Limit"  default(10)
// @Param       offset query    int    false "Offset" default(0)
// @Success     200    {object} response.ListTasksResp
// @Failure     400    {object} response.Error
// @Failure     401    {object} response.Error
// @Failure     500    {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [get]
func (r *V1) listTasks(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	var status *domain.TaskStatus

	if rawStatus := ctx.Query("status"); rawStatus != "" {
		taskStatus := domain.TaskStatus(rawStatus)
		if !taskStatus.Valid() {
			errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid task status")

			return
		}

		status = &taskStatus
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid limit")

		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, apperror.CodeInvalidRequest, "invalid offset")

		return
	}

	tasks, total, err := r.tk.List(ctx.Request.Context(), userID, status, ctx.Query("q"), limit, offset)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - listTasks")
		errorResponse(ctx, http.StatusInternalServerError, apperror.CodeInternalServer, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, response.NewListTasksResp(tasks, total))
}

// @Summary     Get task
// @Description Get a task by ID
// @ID          get-task
// @Tags        tasks
// @Produce     json
// @Param       id  path     string true "Task ID"
// @Success     200 {object} response.TaskResp
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [get]
func (r *V1) getTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	task, err := r.tk.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - getTask")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, response.NewTaskResp(&task))
}

// @Summary     Update task
// @Description Update task title and description
// @ID          update-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "Task ID"
// @Param       request body     request.UpdateTaskReq true "Updated task data"
// @Success     200     {object} response.TaskResp
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [put]
func (r *V1) updateTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	var body request.UpdateTaskReq

	if ok := r.bindJSON(ctx, &body, "restapi - v1 - updateTask"); !ok {
		return
	}

	task, err := r.tk.Update(ctx.Request.Context(), userID, ctx.Param("id"), body.Title, body.Description)
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - updateTask")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, response.NewTaskResp(&task))
}

// @Summary     Transition task status
// @Description Change task status (todo -> in_progress -> done, or in_progress -> todo)
// @ID          transition-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string                true "Task ID"
// @Param       request body     request.TransitionTaskReq true "New status"
// @Success     200     {object} response.TaskResp
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id}/status [patch]
func (r *V1) transitionTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	var body request.TransitionTaskReq

	if ok := r.bindJSON(ctx, &body, "restapi - v1 - transitionTask"); !ok {
		return
	}

	task, err := r.tk.Transition(ctx.Request.Context(), userID, ctx.Param("id"), domain.TaskStatus(body.Status))
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - transitionTask")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, response.NewTaskResp(&task))
}

// @Summary     Delete task
// @Description Delete a task by ID
// @ID          delete-task
// @Tags        tasks
// @Param       id  path     string true "Task ID"
// @Success     204 "No Content"
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [delete]
func (r *V1) deleteTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, apperror.CodeUnauthorized, "unauthorized")

		return
	}

	err := r.tk.Delete(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		apperror.Log(r.l, err, "restapi - v1 - deleteTask")
		mappedErrorResponse(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
