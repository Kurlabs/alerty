package handlers

import (
	"log"
	"net/http"

	"github.com/Kurlabs/alerty/shared/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"admin",
		jwt.StandardClaims{
			Id: "main_user_id",
			// ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte(env.Config.BrainToken))
	if err != nil {
		return "", err
	}
	return token, nil
}

func Login(c echo.Context) error {
	// TODO create jwt token
	token, err := createJwtToken()
	if err != nil {
		log.Println("Error creating JWT token", err)
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "You were logged in!",
		"token":   token,
	})

}
