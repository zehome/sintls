package sintls

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// "user" will be set on the gin context
const AuthUserKey = "user"

// the key is the user name and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
func BasicAuth(db *pg.DB, isadmin bool) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	var dbauth Authorization
	return func(c *gin.Context) {
		var username, password string
		auth := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(auth) == 2 && auth[0] == "Basic" && len(auth[1]) != 0 {
			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pairs := strings.SplitN(string(payload), ":", 2)
			if len(pairs) != 2 || len(pairs[0]) == 0 || len(pairs[1]) == 0 {
				goto Exit401
			}
			username, password = pairs[0], pairs[1]
			err := db.Model(&dbauth).
				Where("name = ?", username).
				Where("admin is ?", isadmin).
				Limit(1).
				Select()
			if err != nil {
				log.Println("User search failed", err)
				goto Exit401
			} else {
				if bcrypt.CompareHashAndPassword([]byte(dbauth.Secret), []byte(password)) != nil {
					log.Println("password does not match")
					goto Exit401
				} else {
					c.Set(AuthUserKey, &dbauth)
					return
				}
			}
		}
	Exit401:
		// Credentials doesn't match, we return 401 and abort handlers chain.
		c.Header("WWW-Authenticate", realm)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
