package handlers

import (
	"net/http"
	"time"
	"togo/models"

	"github.com/labstack/echo/v4"
)

func addTask(c echo.Context) error {
	addTaskRequest := new(models.Task)
	if err := c.Bind(addTaskRequest); err != nil {
		return err
	}

	// get user_id from token
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, "missing user_id")
	}

	addTaskRequest.UserID = userID

	// validate request
	if err := c.Validate(addTaskRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	now := time.Now().Unix()

	addTaskRequest.CreatedAt = now
	addTaskRequest.UpdatedAt = now

	err := mongoClient.InsertTask(addTaskRequest)
	if err != nil {
		log.Errorf("InsertTask failed - addTaskRequest: %+v - err: %v", addTaskRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, addTaskRequest)
}
