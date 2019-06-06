package main

import (
	"database/sql"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/go-pg/pg"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"strconv"
	"strings"
	"time"
)

import "github.com/zehome/sintls/sintls"

// Anwser to bcrypt hash
func TransformBcrypt(ans interface{}) interface{} {
	transformer := survey.TransformString(
		func(s string) string {
			out, _ := bcrypt.GenerateFromPassword(
				[]byte(s), bcrypt.DefaultCost)
			return string(out)
		},
	)
	return transformer(ans)
}

var adduserQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Username:",
		},
		Validate: survey.Required,
	},
	{
		Name:     "subdomain",
		Prompt:   &survey.Input{Message: "Subdomain:"},
		Validate: survey.Required,
	},
	{
		Name:      "password",
		Prompt:    &survey.Password{Message: "Password:"},
		Validate:  survey.Required,
		Transform: TransformBcrypt,
	},
	{
		Name:     "admin",
		Prompt:   &survey.Confirm{Message: "Admin user?", Default: false},
		Validate: survey.Required,
	},
}

var deluserQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Username",
		},
	},
}

func CheckArg(s string, m string) bool {
	return strings.HasPrefix(strings.ToLower(s), m)
}

func RunCLI(db *pg.DB, disable_colors bool, args []string) {
	// colorizer
	var au = aurora.NewAurora(!disable_colors)

	if len(args) == 0 {
		log.Println("RunCLI without arguments")
		return
	}
	if CheckArg(args[0], "h") {
		fmt.Println("commands: help, adduser, deluser, list [authorization,subdomain,host]")
	} else if CheckArg(args[0], "a") {
		adduser_anwsers := struct {
			Name      string
			Password  string
			Subdomain string
			Admin     bool
		}{}
		// ask the question
		err := survey.Ask(adduserQs, &adduser_anwsers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		authorization := sintls.Authorization{
			Name:   adduser_anwsers.Name,
			Secret: adduser_anwsers.Password,
			Admin:  sql.NullBool{Valid: true, Bool: adduser_anwsers.Admin},
		}
		_, err = db.Model(&authorization).Returning("*").Insert()
		if err != nil {
			log.Println("INSERT authorization failed: ", err)
			return
		}
		err = db.Insert(&sintls.SubDomain{
			Name:            adduser_anwsers.Subdomain,
			AuthorizationId: authorization.AuthorizationId,
		})
		if err != nil {
			log.Println("INSERT SubDomain failed: ", err)
			return
		}
	} else if CheckArg(args[0], "d") {
		deluser_anwsers := struct {
			Name string
		}{}
		err := survey.Ask(deluserQs, &deluser_anwsers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Authorizations to remove
		var authorizations []sintls.Authorization
		var authorization_ids []uint64
		err = db.Model(&authorizations).Where("name = ?", deluser_anwsers.Name).Select()
		if err != nil {
			log.Println("Select authorization failed: ", err)
			return
		}
		t := tabby.New()
		t.AddHeader("Username", "Admin", "CreatedAt", "UpdatedAt")
		for _, row := range authorizations {
			t.AddLine(
				row.Name,
				strconv.FormatBool(row.Admin.Bool),
				row.CreatedAt.Format("2006-01-02 15:04:05"),
				row.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
			authorization_ids = append(authorization_ids, row.AuthorizationId)
		}
		fmt.Println(au.Red("Following users will be removed").Bold())
		t.Print()
		fmt.Println()

		if len(authorization_ids) > 0 {
			// Hosts to remove
			var hosts []sintls.Host
			err := db.Model(&hosts).Column(
				"host.name",
				"host.updated_at",
				"host.dns_target_a",
				"host.dns_target_aaaa",
				"host.dns_target_cname",
				"SubDomain",
				"SubDomain.Authorization").
				Order("sub_domain__authorization.name ASC").
				Order("sub_domain.name ASC").
				Order("host.name ASC").
				WhereIn("sub_domain__authorization.authorization_id IN (?)", authorization_ids).
				Select()
			if err != nil {
				log.Println(err)
				return
			}
			t := tabby.New()
			t.AddHeader("Username", "SubDomain", "Name", "A", "AAAA", "CNAME", "UpdatedAt", "Expires")
			for _, row := range hosts {
				expires := time.Now().Add(90 * 24 * time.Hour).Sub(row.UpdatedAt)
				t.AddLine(
					row.SubDomain.Authorization.Name,
					row.SubDomain.Name,
					row.Name,
					row.DnsTargetA,
					row.DnsTargetAAAA,
					row.DnsTargetCNAME,
					row.UpdatedAt.Format("2006-01-02 15:04:05"),
					fmt.Sprintf("%d days", expires / (24 * time.Hour)),
				)
			}
			fmt.Println(au.Red("Following hosts will be removed").Bold())
			t.Print()
		}

		anwser := false
		survey.AskOne(
			&survey.Confirm{Message: "Confirm you want to delete?"},
			&anwser, nil,
		)
		if anwser {
			fmt.Printf("Removing user %s\n", deluser_anwsers.Name)
			_, err = db.Model(&authorizations).Delete()
			if err != nil {
				panic(err)
			} else {
				log.Printf("Removed user %s (CLI)\n", deluser_anwsers.Name)
			}
		}
	} else if CheckArg(args[0], "l") {
		if len(args) < 2 || CheckArg(args[1], "a") {
			var authorizations []sintls.Authorization
			err := db.Model(&authorizations).Select()
			if err != nil {
				log.Println(err)
				return
			}
			t := tabby.New()
			t.AddHeader("Username", "Admin", "CreatedAt", "UpdatedAt")
			for _, row := range authorizations {
				t.AddLine(
					row.Name,
					strconv.FormatBool(row.Admin.Bool),
					row.CreatedAt.Format("2006-01-02 15:04:05"),
					row.UpdatedAt.Format("2006-01-02 15:04:05"),
				)
			}
			fmt.Println("Authorizations:")
			t.Print()
		}
		if len(args) < 2 || CheckArg(args[1], "s") {
			var subdomains []sintls.SubDomain
			err := db.Model(&subdomains).Column(
				"sub_domain.name",
				"sub_domain.updated_at",
				"Authorization").Order("authorization.name ASC").Order("sub_domain.name ASC").Select()
			if err != nil {
				log.Println(err)
				return
			}
			t := tabby.New()
			t.AddHeader("Username", "Name", "Updated")
			for _, row := range subdomains {
				t.AddLine(
					row.Authorization.Name,
					row.Name,
					row.UpdatedAt.Format("2006-01-02 15:04:05"),
				)
			}
			fmt.Println("Subdomains:")
			t.Print()
		}
		if len(args) < 2 || CheckArg(args[1], "h") {
			var hosts []sintls.Host
			err := db.Model(&hosts).Column(
				"host.name",
				"host.updated_at",
				"host.dns_target_a",
				"host.dns_target_aaaa",
				"host.dns_target_cname",
				"SubDomain",
				"SubDomain.Authorization").
				Order("sub_domain__authorization.name ASC").
				Order("sub_domain.name ASC").
				Order("host.name ASC").Select()
			if err != nil {
				log.Println(err)
				return
			}
			t := tabby.New()
			t.AddHeader("Username", "SubDomain", "Name", "A", "AAAA", "CNAME", "UpdatedAt", "Expires")
			for _, row := range hosts {
				expires := time.Now().Add(90 * 24 * time.Hour).Sub(row.UpdatedAt)
				t.AddLine(
					row.SubDomain.Authorization.Name,
					row.SubDomain.Name,
					row.Name,
					row.DnsTargetA,
					row.DnsTargetAAAA,
					row.DnsTargetCNAME,
					row.UpdatedAt.Format("2006-01-02 15:04:05"),
					fmt.Sprintf("%d days", expires / (24 * time.Hour)),
				)
			}
			fmt.Println("Hosts:")
			t.Print()
		}
	} else {
		log.Println("unknown command")
	}
}
