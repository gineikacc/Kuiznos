package sessions

import (
	"myapp/models"

	"github.com/labstack/echo/v4"
)

func SessionMiddleware(sm *models.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionToken, err := c.Cookie("session_token")
			if err == nil {
				// Session token found, retrieve session from the manager
				session, exists := sm.Get_session(sessionToken.Value)
				if exists {
					c.Set("session", session)
					c.Set("authorized", !session.Is_expired())
				} else {
					c.Set("authorized", false)
				}
			} else {
				c.Set("authorized", false)
			}
			return next(c)
		}
	}
}
func TrailingSlash() echo.MiddlewareFunc {
	// Defaults

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			req := c.Request()
			url := req.URL
			path := url.Path
			qs := c.QueryString()
			if path != "/" && path[len(path)-1] != '/' {
				path += "/"
				uri := path
				if qs != "" {
					uri += "?" + qs
				}
			}
			return next(c)
		}
	}
}
