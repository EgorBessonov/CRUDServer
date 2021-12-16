package handler

import (
	"CRUDServer/internal/repository"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

/
type customClaims struct{
	email string `json:"email"`
	jwt.StandardClaims
}
// Registration is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	hashedPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Registration failed.")
	}
	err = h.rps.CreateAuthUser(repository.RegistrationForm{
		UserName: c.QueryParam("userName"),
		Email:    c.QueryParam("email"),
		Password: hashedPassword,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while saving form.")
	}
	return c.String(http.StatusOK, "Successfully.")
}

// Authentication check user password and if it correct returns access and refresh tokens 
func (h Handler) Authentication(c echo.Context) error {
	authUser, err := h.rps.GetAuthUser(c.QueryParam("password"))
	if err != nil{
		return c.String(http.StatusBadRequest, "Incorrect email")
	}
	hashedLoginPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil{
		return c.String(http.StatusInternalServerError, "Login failed, try again")
	}
	if authUser.Password != hashedLoginPassword{
		return c.String(http.StatusBadRequest, "Incorrect password")
	} 

	
	return c.String(http.StatusOK, "Successfully login")
}

func (h Handler) Authorization(c echo.Context) error{
	
	return nil
}

func (h Handler) RefreshToken(c echo.Context) error{

	return nil
} 

func (h Handler) Logout(c echo.Context) error{

	return nil
}

func hashPassword(password string) (string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "hashPassword()",
			"err":    err,
		}).Info("Authentication info.")
		return "", nil
	}
	return string(hashedPassword), nil
}
