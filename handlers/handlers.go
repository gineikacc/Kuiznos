package handlers

import (
	"fmt"
	"log"
	"myapp/auth"
	"myapp/database"
	"myapp/models"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

var DB database.MysqlStore
var AUTH auth.Auth
var SM *models.SessionManager

func Home_GET(c echo.Context) error {
	data := map[string]interface{}{}
	session_token, err := c.Cookie("session_token")
	if err == nil {
		data["session_token"] = session_token.Value
	}
	return c.Render(http.StatusOK, "index", data)
}

func Login_GET(c echo.Context) error {
	if c.Get("authorized") != nil {
		if c.Get("authorized").(bool) {
			return c.Redirect(http.StatusMovedPermanently, "/")
		}
	}
	return c.Render(http.StatusOK, "form", 0)
}

func Login_POST(c echo.Context) error {
	fmt.Println("REE: LOGIN POST")
	user, _ := read_user_form(c)
	authorized, _ := AUTH.Auth_with_credentials(*user)

	if authorized {
		session_token := AUTH.Create_session_token(*user)
		c.SetCookie(session_token)
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func Register_POST(c echo.Context) error {
	fmt.Println("REE: REGISTER POST")

	user, _ := read_user_form(c)
	exists := DB.User_exists(strings.ToLower(user.Name))
	if exists {
		log.Fatal(fmt.Errorf("> %v is already registered", user.Name))
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	err := DB.Create_user(*user)
	if err == nil {
		session_token := AUTH.Create_session_token(*user)
		c.SetCookie(session_token)
	} else {
		log.Fatal(err)
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func Logout_POST(c echo.Context) error {

	session, ok := c.Get("session").(*models.Session)
	if ok {
		session.Expiry = time.Now().Add(-1 * time.Hour)
		c.SetCookie(&http.Cookie{
			Name:    "session_token",
			Value:   "_",
			Expires: time.Now().Add(-1 * time.Minute),
		})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func read_user_form(c echo.Context) (*models.User, error) {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return user, err
	}
	user.Password = AUTH.Hash(user.Password)
	return user, nil
}
