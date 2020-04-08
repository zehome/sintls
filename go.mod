module github.com/zehome/sintls

require (
	github.com/Azure/azure-sdk-for-go v36.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.60.282 // indirect
	github.com/aws/aws-sdk-go v1.25.47 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/cheynewallace/tabby v1.1.0
	github.com/cloudflare/cloudflare-go v0.10.8 // indirect
	github.com/coreos/go-systemd v0.0.0-20181012123002-c6f51f82210d
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/exoscale/egoscale v1.19.0 // indirect
	github.com/go-acme/lego/v3 v3.5.0
	github.com/go-pg/migrations/v7 v7.1.9
	github.com/go-pg/pg/v9 v9.1.5
	github.com/gophercloud/gophercloud v0.7.0 // indirect
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/kolo/xmlrpc v0.0.0-20190909154602-56d5ec7c422e // indirect
	github.com/kr/pty v1.1.3 // indirect
	github.com/labstack/echo/v4 v4.1.16
	github.com/linode/linodego v0.12.2 // indirect
	github.com/liquidweb/liquidweb-go v1.6.1 // indirect
	github.com/logrusorgru/aurora v0.0.0-20191116043053-66b7ad493a23
	github.com/miekg/dns v1.1.27
	github.com/oracle/oci-go-sdk v13.0.0+incompatible // indirect
	github.com/ovh/go-ovh v0.0.0-20181109152953-ba5adb4cf014
	github.com/rhysd/go-github-selfupdate v1.2.1
	github.com/sacloud/libsacloud v1.32.0 // indirect
	github.com/transip/gotransip v5.8.2+incompatible // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/urfave/cli v1.22.2
	github.com/vultr/govultr v0.1.7 // indirect
	go.uber.org/ratelimit v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/ns1/ns1-go.v2 v2.0.0-20191126161805-25b9eac84517 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
)

// broken go.mod inside a package needing this
replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

// Otherwise, the build will fail in dh_golang
//replace google.golang.org/api => google.golang.org/api v0.3.0

go 1.13
