package handlers

import (
	"myapp/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Home_GET(c echo.Context) error {
	data := map[string]interface{}{}
	session_token, err := c.Cookie("session_token")
	if err == nil {
		data["session_token"] = session_token.Value
	}
	return c.Render(http.StatusOK, "index", data)
}

func Login_GET(c echo.Context) error {
	return c.Render(http.StatusOK, "form", 0)
}

func Login_POST(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	data := map[string]interface{}{}

	if user.Name == "airolen11" {
		cookie := NewCookie()
		c.SetCookie(cookie)
		c.SetCookie(&http.Cookie{Name: "name", Value: "airolen11", Expires: cookie.Expires})

		data["name"] = user.Name
		data["email"] = user.Email
		data["password"] = user.Password
		data["session_token"] = cookie.Value
	}
	return c.Redirect(http.StatusMovedPermanently, "/")

}

func Logout_POST(c echo.Context) error {

	_, ok := c.Get("session").(models.Session)
	if ok {
		c.SetCookie(&http.Cookie{
			Name:    "session_token",
			Value:   "_",
			Expires: time.Now().Add(-1 * time.Minute),
		})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")

}

func NewCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = "WEOWCOOK"
	cookie.Expires = time.Now().Add(time.Minute * 2)
	return cookie
}
