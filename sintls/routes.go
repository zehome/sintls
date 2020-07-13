package sintls

import (
	"fmt"
	"github.com/go-acme/lego/v3/challenge"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/zehome/sintls/dns"
	"log"
	"net"
	"net/http"
)

// Lego httpreq RAW request
type LegoMessage struct {
	Domain      string `json:"domain"`
	Token       string `json:"token"`
	KeyAuth     string `json:"keyAuth"`
	TargetA     net.IP `json:"dnstarget_a,omitempty"`
	TargetAAAA  net.IP `json:"dnstarget_aaaa,omitempty"`
	TargetCNAME string `json:"dnstarget_cname,omitempty"`
	TargetMX    string `json:"dnstarget_mx,omitempty"`
}

func updateDNSRecords(tx *pg.Tx, req LegoMessage, user *Authorization, dnsupdater dns.DNSUpdater) (err error) {
	err = user.CreateOrUpdateHost(
		tx, req.Domain, req.TargetA, req.TargetAAAA, req.TargetCNAME,
		req.TargetMX)
	if err != nil {
		tx.Rollback()
		log.Printf("sintls: update host failed: %s", err)
		return
	}
	// Remove previous entries
	dnsupdater.RemoveRecords(req.Domain)

	// Update DNS records
	if len(req.TargetA) != 0 {
		err = dnsupdater.SetRecord(req.Domain, "A", req.TargetA.String())
		if err != nil {
			log.Printf("sintls: setrecord A failed: %s", err)
			return
		}
	}
	if len(req.TargetAAAA) != 0 {
		err = dnsupdater.SetRecord(req.Domain, "AAAA", req.TargetAAAA.String())
		if err != nil {
			log.Printf("sintls: setrecord AAAA failed: %s", err)
			return
		}
	}
	if len(req.TargetA) != 0 && len(req.TargetAAAA) != 0 && len(req.TargetCNAME) > 0 {
		err = dnsupdater.SetRecord(req.Domain, "CNAME", req.TargetCNAME)
		return
	}
	if len(req.TargetMX) > 0 {
		err = dnsupdater.SetRecord(req.Domain, "MX", req.TargetMX)
		return
	}
	dnsupdater.Refresh(req.Domain)
	return
}

func UpdateDNSRecords(c echo.Context) error {
	var req LegoMessage
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}
	user := c.Get("user").(*Authorization)
	db := c.Get("database").(*pg.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user auth is required")
	}
	if !user.CanUseHost(db, req.Domain) {
		return echo.NewHTTPError(http.StatusForbidden, "no permissions to use this domain")
	}
	// Custom DNS Updater
	dnsupdater := c.Get("dnsupdater").(dns.DNSUpdater)
	tx, err := db.Begin()
	if err != nil {
		log.Printf("sintls: db.Begin() failed", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "begin failed")
	}

	err = updateDNSRecords(tx, req, user, dnsupdater)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "update dnsrecords failed")
	}
	if err := tx.Commit(); err != nil {
		log.Printf("sintls: tx.Commit() failed:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	} else {
		return c.String(http.StatusOK, "success")
	}
}

func LegoPresent(c echo.Context) error {
	var req LegoMessage
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	if len(req.TargetA) == 0 && len(req.TargetAAAA) == 0 && len(req.TargetCNAME) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "either A, AAAA or CNAME must be defined")
	}
	user := c.Get("user").(*Authorization)
	db := c.Get("database").(*pg.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user auth is required")
	}
	if !user.CanUseHost(db, req.Domain) {
		return echo.NewHTTPError(http.StatusForbidden, "no permissions to use this domain")
	}
	// Lego DNS Provider
	provider := c.Get("dnsprovider").(challenge.Provider)
	// Custom DNS Updater
	dnsupdater := c.Get("dnsupdater").(dns.DNSUpdater)

	tx, err := db.Begin()
	if err != nil {
		log.Printf("sintls: db.Begin() failed:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "begin failed")
	}

	err = updateDNSRecords(tx, req, user, dnsupdater)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "update dnsrecords failed")
	}

	// Lego challenge
	err = provider.Present(req.Domain, req.Token, req.KeyAuth)
	if err != nil {
		tx.Rollback()
		log.Printf("sintls: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "lego present failed")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("sintls: tx.Commit() failed: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "commit failed")
	} else {
		return c.String(http.StatusOK, "success")
	}
}

func LegoCleanup(c echo.Context) error {
	var req LegoMessage
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}
	user := c.Get("user").(*Authorization)
	db := c.Get("database").(*pg.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user auth is required")
	}
	if !user.CanUseHost(db, req.Domain) {
		return echo.NewHTTPError(http.StatusForbidden, "no permissions to use this domain")
	}
	provider := c.Get("dnsprovider").(challenge.Provider)
	err := provider.CleanUp(req.Domain, req.Token, req.KeyAuth)
	if err != nil {
		log.Printf("sintls: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "lego Cleanup failed")
	}
	return c.String(http.StatusOK, "success")
}

// type CreateAuthMessage struct {
// 	Name   string `json:"name" binding:"required"`
// 	Secret string `json:"secret" binding:"required"`
// }

// func CreateAuth(c echo.Context) {
// 	var authmessage CreateAuthMessage
// 	user := c.Get("user").(*Authorization)
// 	if user.Admin.Bool != true {
// 		return c.AbortWithStatus(http.StatusUnauthorized)
// 	}
// 	db := c.Get("database").(*pg.DB)
// 	if err := c.Bind(&authmessage); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "bad request: %s", err.Error())
// 	}
// 	var count int = 0
// 	count, err := db.Model((*Authorization)(nil)).
// 		Where("name = ?", authmessage.Name).
// 		Count()
// 	if count != 0 || err != nil {
// 		return echo.NewHTTPError(http.StatusForbidden, "Authorization %s already exists", authmessage.Name)
// 	}
// 	hashpw, err := bcrypt.GenerateFromPassword(
// 		[]byte(authmessage.Secret), bcrypt.DefaultCost)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "unhashable password: %s", err)
// 	}
// 	dbauth := Authorization{
// 		Name:   authmessage.Name,
// 		Secret: string(hashpw),
// 	}
// 	if err := db.Insert(&dbauth); err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error: %s", err)
// 	}
// 	return c.String(http.StatusOK, "success")
// }
