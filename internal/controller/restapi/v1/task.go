package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi/v1/request"
	"github.com/bhcoder23/gin-clean-template/internal/controller/restapi/v1/response"
	"github.com/bhcoder23/gin-clean-template/internal/entity"
	"github.com/gin-gonic/gin"
)

// @Summary     Create task
// @Description Create a new task for the current user
// @ID          create-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body     request.CreateTask true "Task data"
// @Success     201     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [post]
func (r *V1) createTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	var body request.CreateTask

	if err := ctx.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - createTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - createTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Create(ctx.Request.Context(), userID, body.Title, body.Description)
	if err != nil {
		r.l.Error(err, "restapi - v1 - createTask")

		if errors.Is(err, entity.ErrTaskTitleRequired) {
			errorResponse(ctx, http.StatusBadRequest, "task title is required")

			return
		}

		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusCreated, task)
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
// @Success     200    {object} response.TaskList
// @Failure     400    {object} response.Error
// @Failure     401    {object} response.Error
// @Failure     500    {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [get]
func (r *V1) listTasks(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	var status *entity.TaskStatus

	if rawStatus := ctx.Query("status"); rawStatus != "" {
		taskStatus := entity.TaskStatus(rawStatus)
		if !taskStatus.Valid() {
			errorResponse(ctx, http.StatusBadRequest, "invalid task status")

			return
		}

		status = &taskStatus
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

	tasks, total, err := r.tk.List(ctx.Request.Context(), userID, status, ctx.Query("q"), limit, offset)
	if err != nil {
		r.l.Error(err, "restapi - v1 - listTasks")
		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, response.TaskList{
		Tasks: tasks,
		Total: total,
	})
}

// @Summary     Get task
// @Description Get a task by ID
// @ID          get-task
// @Tags        tasks
// @Produce     json
// @Param       id  path     string true "Task ID"
// @Success     200 {object} entity.Task
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [get]
func (r *V1) getTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	task, err := r.tk.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "restapi - v1 - getTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(ctx, http.StatusNotFound, "task not found")

			return
		}

		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, task)
}

// @Summary     Update task
// @Description Update task title and description
// @ID          update-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "Task ID"
// @Param       request body     request.UpdateTask  true "Updated task data"
// @Success     200     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [put]
func (r *V1) updateTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	var body request.UpdateTask

	if err := ctx.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Update(ctx.Request.Context(), userID, ctx.Param("id"), body.Title, body.Description)
	if err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(ctx, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrTaskTitleRequired) {
			errorResponse(ctx, http.StatusBadRequest, "task title is required")

			return
		}

		if errors.Is(err, entity.ErrTaskCompleted) {
			errorResponse(ctx, http.StatusBadRequest, "completed task cannot be modified")

			return
		}

		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, task)
}

// @Summary     Transition task status
// @Description Change task status (todo -> in_progress -> done, or in_progress -> todo)
// @ID          transition-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string                true "Task ID"
// @Param       request body     request.TransitionTask  true "New status"
// @Success     200     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id}/status [patch]
func (r *V1) transitionTask(ctx *gin.Context) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	var body request.TransitionTask

	if err := ctx.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Transition(ctx.Request.Context(), userID, ctx.Param("id"), body.Status)
	if err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(ctx, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrInvalidTransition) {
			errorResponse(ctx, http.StatusBadRequest, "invalid status transition")

			return
		}

		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.JSON(http.StatusOK, task)
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
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized")

		return
	}

	err := r.tk.Delete(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "restapi - v1 - deleteTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(ctx, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrTaskCompleted) {
			errorResponse(ctx, http.StatusBadRequest, "completed task cannot be modified")

			return
		}

		errorResponse(ctx, http.StatusInternalServerError, "internal server error")

		return
	}

	ctx.Status(http.StatusNoContent)
}
