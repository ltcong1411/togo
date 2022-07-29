package handlers

import (
	b64 "encoding/base64"
	"net/http"
	"time"
	"togo/config"
	"togo/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func registerUser(c echo.Context) error {
	request := new(models.User)
	if err := c.Bind(request); err != nil {
		return err
	}

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

	if user != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "this user already exists!")
	}

	user = request

	if user.DailyTaskLimit == 0 {
		user.DailyTaskLimit = config.Values.DailyTaskLimitDefault
	}

	user.Password = b64.StdEncoding.EncodeToString([]byte(user.Password))
	now := time.Now().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now

	// insert user
	err = mongoClient.InsertUser(user)
	if err != nil {
		log.Errorf("InsertUser failed - user: %+v - err: %v", user, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func login(c echo.Context) error {
	request := new(models.User)
	if err := c.Bind(request); err != nil {
		return err
	}

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

	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "this user not exist!")
	}

	// compare password hash
	passwordBase64 := b64.StdEncoding.EncodeToString([]byte(request.Password))

	if passwordBase64 != user.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password!")
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		user.ID.Hex(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Values.JWTSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
