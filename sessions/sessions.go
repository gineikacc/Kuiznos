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
				session, _ := sm.GetSession(sessionToken.Value)
				c.Set("session", session)
			}
			return next(c)
		}
	}
}
