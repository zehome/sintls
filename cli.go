package main

import (
	"database/sql"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"strconv"
	"strings"
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
		Name:      "subdomain",
		Prompt:    &survey.Input{Message: "Subdomain:"},
		Validate:  survey.Required,
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

func CheckArg(s string, m string) bool {
	return strings.HasPrefix(strings.ToLower(s), m)
}

func RunCLI(db *pg.DB, args []string) {
	if len(args) == 0 {
		log.Println("RunCLI without arguments")
		return
	}
	if CheckArg(args[0], "h") {
		fmt.Println("commands: help, adduser, list [authorization,subdomain,host]")
	} else if CheckArg(args[0], "a") {
		adduser_anwsers := struct {
			Name     	string
			Password 	string
			Subdomain 	string
			Admin    	bool
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
			Name:				adduser_anwsers.Subdomain,
			AuthorizationId:	authorization.AuthorizationId,
		})
		if err != nil {
			log.Println("INSERT SubDomain failed: ", err)
			return
		}
	} else if CheckArg(args[0], "l") {

		//var hosts []sintls.Host
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
			t.AddHeader("Username", "SubDomain", "Name", "TargetA", "TargetAAAA", "UpdatedAt")
			for _, row := range hosts {
				t.AddLine(
					row.SubDomain.Authorization.Name,
					row.SubDomain.Name,
					row.Name,
					row.UpdatedAt.Format("2006-01-02 15:04:05"),
					row.DnsTargetA,
					row.DnsTargetAAAA,
				)
			}
			fmt.Println("Hosts:")
			t.Print()
		}
	} else {
		log.Println("unknown command")
	}
}
