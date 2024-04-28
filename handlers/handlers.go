package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/auth"
	"myapp/database"
	"myapp/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

var DB database.MysqlStore
var AUTH auth.Auth
var SM *models.SessionManager

func Home_GET(c echo.Context) error {
	data := map[string]interface{}{}
	if c.Get("authorized") != nil {
		authorized := c.Get("authorized").(bool)
		data["authorized"] = authorized
		if c.Get("session") != nil {
			data["username"] = c.Get("session").(*models.Session).Username
		}
		data["quizzes"] = DB.Query_Quizzes("")
		c.SetCookie(&http.Cookie{Name: "auth", Value: strconv.FormatBool(authorized), Expires: time.Now().Add(time.Hour)})
	}
	return c.Render(http.StatusOK, "index", data)
}

func Login_GET(c echo.Context) error {
	if c.Get("authorized").(bool) {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	return c.Render(http.StatusOK, "form", 0)
}

func Login_POST(c echo.Context) error {
	user, _ := read_user_form(c)
	authorized, _ := AUTH.Auth_with_credentials(*user)
	if authorized {
		session_token := AUTH.Create_session_token(*user)
		c.SetCookie(session_token)
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func Register_POST(c echo.Context) error {
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
	st, err := c.Cookie("session_token")
	if err == nil {
		session, ok := SM.Get_session(st.Value)
		if ok {
			session.Expiry = time.Now().Add(-1 * time.Minute)
			c.SetCookie(&http.Cookie{
				Name:    "session_token",
				Value:   "_",
				Expires: time.Now().Add(-1 * time.Minute),
			})
		}
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func Search_GET(c echo.Context) error {
	var quizzes []models.Quiz
	if c.Get("authorized") != nil {
		quizzes = DB.Query_Quizzes(c.QueryParams().Get("query"))
	}
	return c.Render(http.StatusOK, "searchResults", quizzes)
}

func CreateGame_POST(c echo.Context) error {
	if c.Get("authorized") != nil && c.Get("session") != nil {
		quiz, err := read_quiz_form(c)
		if err != nil {
			log.Fatalf("%+v\n", err)
		} else {
			DB.Create_quiz(*quiz, c.Get("session").(*models.Session).Username)
		}

	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func PlayGET(c echo.Context) error {
	return c.Render(http.StatusOK, "play", nil)
}

func GameGET(c echo.Context) error {
	quiz := &models.Quiz{}
	if c.Get("authorized") != nil {
		*quiz, _ = DB.Read_quiz(c.QueryParams().Get("title"))
	}
	return c.JSON(http.StatusOK, quiz)
}
func GameEndGET(c echo.Context) error {
	if c.Get("authorized") != nil && c.Get("session") != nil {
		DB.Increment_user_quiz_done_count(c.Get("session").(*models.Session).Username)
	}
	return c.String(http.StatusOK, ":)")
}

func Hx_createGame(c echo.Context) error {
	if c.Get("authorized") != nil {
		return c.Render(http.StatusOK, "gameCreationForm", nil)
	}
	return c.Render(http.StatusOK, "", nil)
}

func read_user_form(c echo.Context) (*models.User, error) {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return user, err
	}
	user.Password = AUTH.Hash(user.Password)
	return user, nil
}

func read_quiz_form(c echo.Context) (*models.Quiz, error) {
	quiz := new(models.Quiz)
	str := c.FormValue("quiz")
	err := json.Unmarshal([]byte(str), &quiz)
	if err != nil {
		return quiz, err
	}
	return quiz, nil
}

func Hx_questionEntry(c echo.Context) error {
	data := map[string]interface{}{}
	return c.Render(http.StatusOK, "questionEntry", data)
}
