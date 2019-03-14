package sintls

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

/* Lego httpreq RAW request */
type LegoMessage struct {
	Domain  string `json:"domain"`
	Token   string `json:"token"`
	KeyAuth string `json:"keyAuth"`
}

func LegoPresent(c *gin.Context) {
	var req LegoMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request:", err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}

func LegoCleanup(c *gin.Context) {
	var req LegoMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request:", err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}

type CreateAuthMessage struct {
	Name   string `json:"name" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

func CreateAuth(c *gin.Context) {
	var authmessage CreateAuthMessage
	user := c.MustGet("user").(Authorization)
	if user.Admin.Bool != true {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	db := c.MustGet("database").(*pg.DB)
	if err := c.ShouldBindJSON(&authmessage); err != nil {
		c.String(http.StatusBadRequest, "bad request: %s", err.Error())
		return
	}
	var count int = 0
	count, err := db.Model((*Authorization)(nil)).
		Where("name = ?", authmessage.Name).
		Count()
	if count != 0 || err != nil {
		c.String(http.StatusForbidden, "Authorization %s already exists", authmessage.Name)
		return
	}
	hashpw, err := bcrypt.GenerateFromPassword(
		[]byte(authmessage.Secret), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusBadRequest, "unhashable password: %s", err)
		return
	}
	dbauth := Authorization{
		Name:   authmessage.Name,
		Secret: string(hashpw),
	}
	if err := db.Insert(&dbauth); err != nil {
		c.String(http.StatusInternalServerError, "Unknown error: %s", err)
		return
	}
	c.String(http.StatusOK, "success")
}

type DeleteAuthMessage struct {
	Name string `json:"name" binding:"required"`
}

func DeleteAuth(c *gin.Context) {
	var req DeleteAuthMessage
	//var auth Authorization
	user := c.MustGet("user").(Authorization)
	if user.Admin.Bool != true {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	db := c.MustGet("database").(*pg.DB)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "bad request: %s", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("db.Begin() failed: ", err)
		return
	}
	// Rollback tx on error.
	defer tx.Rollback()

	// tx := db.Begin()
	// if tx.Error != nil {
	// 	log.Fatal("db.Begin() failed", tx.Error)
	// 	c.String(http.StatusInternalServerError, "begin failed")
	// 	return
	// }
	// if err := db.Where("name = ?", req.Name).First(&auth).Error; err != nil {
	// 	c.String(http.StatusNotFound, "user %s not found", req.Name)
	// 	return
	// }
	// if err := db.Model(&auth).Related(&subdomains).Error; err == nil {
	// 	for _, subdomain := range subdomains {
	// 		if err := db.Model(&subdomain).Related(&hosts).Error; err == nil {
	// 			for _, host := range hosts {
	// 				db.Delete(&host)
	// 			}
	// 		}
	// 	}
	// }

	// if err := tx.Commit().Error; err != nil {
	// 	log.Fatal("tx.Commit() failed: ", err)
	// 	c.String(http.StatusInternalServerError, "commit failed")
	// } else {
	c.String(http.StatusOK, "success")
	//}
}
