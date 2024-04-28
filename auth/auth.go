package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"myapp/database"
	"myapp/models"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	DB database.MysqlStore
	SM *models.SessionManager
}

func New(store database.MysqlStore, sm *models.SessionManager) Auth {
	return Auth{DB: store, SM: sm}
}

func (a Auth) Auth_with_credentials(user models.User) (bool, error) {
	db_user, err := a.DB.Read_user(user.Name)

	if err == nil {
		if strings.EqualFold(user.Password, db_user.Password) {
			return true, nil
		}
		return false, errors.New("passwords dont match")

	}
	return false, errors.New("user doesn't exist")

}
func (a Auth) Auth_with_session(session_cookie http.Cookie) (bool, error) {
	session, exists := a.SM.Get_session(session_cookie.Value)
	if !exists || session.Is_expired() {
		return false, errors.New("session cookie mismatch")
	}
	return true, nil
}

func (a Auth) Create_session_token(user models.User) *http.Cookie {
	cookie := new(http.Cookie)
	expiry := time.Now().Add(24 * time.Hour)

	cookie.Name = "session_token"
	cookie.Value = generate_secure_token(32)
	cookie.Expires = expiry
	session := models.Session{Username: user.Name, Session_token: cookie.Value, Expiry: expiry}
	a.SM.Set_session(&session)
	return cookie
}

func (a Auth) Hash(key string) string {
	h := sha512.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func generate_secure_token(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
