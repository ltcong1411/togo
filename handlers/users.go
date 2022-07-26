package handlers

import (
	b64 "encoding/base64"
	"net/http"
	"time"
	"togo/config"
	"togo/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func registerUser(c echo.Context) error {
	request := new(models.User)
	if err := c.Bind(request); err != nil {
		return err
	}

	log.Debugf("request register: %+v", request)

	// validate request
	if err := c.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get user by username
	user, err := mongoClient.GetUserByUserName(request.Username)
	if err != nil {
		log.Errorf("GetUserByUserName failed - username: %+v - err: %v", request.Username, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Debugf("user: %+v", user)

	if user != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "this user already exists!")
	}

	user = request

	user.Password = b64.StdEncoding.EncodeToString([]byte(user.Password))
	user.CreatedAt = time.Now().Unix()

	// insert user
	err = mongoClient.InsertUser(user)
	if err != nil {
		log.Errorf("InsertUser failed - user: %+v - err: %v", user, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

func login(c echo.Context) error {
	request := new(models.User)
	if err := c.Bind(request); err != nil {
		return err
	}

	log.Debugf("request login: %+v", request)

	// validate request
	if err := c.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get user by username
	user, err := mongoClient.GetUserByUserName(request.Username)
	if err != nil {
		log.Errorf("GetUserByUserName failed - username: %+v - err: %v", request.Username, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Debugf("user: %+v", user)

	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "this user not exist!")
	}

	// compare password hash
	passwordBase64 := b64.StdEncoding.EncodeToString([]byte(request.Password))

	if passwordBase64 != user.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password!")
	}

	// generate token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["admin"] = true
	claims["iss"] = "perfume-shop"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Values.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password!")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
