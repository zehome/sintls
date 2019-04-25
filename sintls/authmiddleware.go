package sintls

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// "user" will be set on the gin context
const AuthUserKey = "user"

func BasicAuth(db *pg.DB, isadmin bool) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		var dbauth Authorization
		err := db.Model(&dbauth).
			Where("name = ?", username).
			Where("admin is ?", isadmin).
			Limit(1).
			Select()
		if err != nil {
			return false, echo.NewHTTPError(http.StatusUnauthorized, "incorrect user or password")
		} else {
			if bcrypt.CompareHashAndPassword([]byte(dbauth.Secret), []byte(password)) != nil {
				log.Println("password does not match")
			} else {
				c.Set(AuthUserKey, &dbauth)
				return true, nil
			}
		}
		return false, echo.NewHTTPError(http.StatusUnauthorized, "incorrect user or password")
	}
}
