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
					c.Set("authorized", !session.IsExpired())
				}
			}
			return next(c)
		}
	}
}
