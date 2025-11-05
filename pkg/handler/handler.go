package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"test-rest-api/pkg/model"
	"test-rest-api/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	logger   *slog.Logger
}

func NewHandler(services *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())

	api := router.Group("/")
	{

		api.POST("/sign-up", h.signUp)
		api.POST("/sign-in", h.signIn)

		users := api.Group("/users", h.userIdentity)
		{
			users.GET("/:id/status", h.getStatus)
			users.GET("/leaderboard", h.getLeaders)
			users.POST("/:id/task/complete", h.completeTask)
			users.POST("/:id/referrer", h.addRefer)
		}
		tasks := api.Group("/tasks", h.userIdentity)
		{
			tasks.POST("/", h.addTask)
			tasks.DELETE("/:id", h.removeTask)
		}

	}
	return router
}

func (h *Handler) getStatus(c *gin.Context) {
	id, err := h.GetIdParam(c)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusInternalServerError, err.Error())
		return
	}
	userstatus, err := h.services.GetUserStatus(id)

	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"userStatus": userstatus})

}

func (h *Handler) getLeaders(c *gin.Context) {
	leaders, err := h.services.GetLeaders()
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"leaders": leaders})
}

func (h *Handler) completeTask(c *gin.Context) {
	id, err := h.GetIdParam(c)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.CompleteTask
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.CompleteTask(id, &input)

	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "The Task succesfully received!"})

}

func (h *Handler) addRefer(c *gin.Context) {
	id, err := h.GetIdParam(c)

	h.services.AddRefferer(id)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "The Refferorer code succesfully received!"})

}

func (h *Handler) addTask(c *gin.Context) {
	var input model.TaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.AddTask(&input)

	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "The Task succesfully received!"})

}
func (h *Handler) removeTask(c *gin.Context) {
	id, err := h.GetIdParam(c)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.RemoveTask(id)

	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "The Task succesfully deleted!"})

}

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}
	if input.Email == "" || input.Password == "" || input.Phone == "" || input.FullName == "" {
		newErrorResponse(h.logger, c, http.StatusBadRequest, "Incorrect data")
		return
	}
	user, err := h.services.GetUser(input.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}
	if user != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, "The Email "+user.Email+" is  already registered")
		return
	}
	_, err = h.services.CreateUser(input)

	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "The user succesfully created!"})
}

func (h *Handler) updateUser(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.UpdateUser(input)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"modifiedCount": id})
}

func (h *Handler) signIn(c *gin.Context) {
	var input model.SignIn

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
	}
	token, refreshToken, user, err := h.services.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = ""

	c.JSON(http.StatusOK, map[string]interface{}{"token": token, "refreshToken": refreshToken, "user": user})
}

func (h *Handler) GetIdParam(c *gin.Context) (int64, error) {
	var d struct {
		ID int64 `uri:"id" `
	}

	if err := c.ShouldBindUri(&d); err != nil {
		return 0, err
	}
	if d.ID == 0 {

		return 0, errors.New("Can't get ID")
	}
	return d.ID, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000/*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
