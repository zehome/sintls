module github.com/zehome/sintls

require (
	cloud.google.com/go v0.49.0 // indirect
	github.com/Azure/azure-sdk-for-go v36.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.3 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.60.282 // indirect
	github.com/aws/aws-sdk-go v1.25.47 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/cenkalti/backoff/v3 v3.1.1 // indirect
	github.com/cheynewallace/tabby v1.1.0
	github.com/cloudflare/cloudflare-go v0.10.8 // indirect
	github.com/coreos/go-systemd v0.0.0-20181012123002-c6f51f82210d
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/exoscale/egoscale v1.19.0 // indirect
	github.com/go-acme/lego/v3 v3.2.0
	github.com/go-ini/ini v1.51.0 // indirect
	github.com/go-pg/migrations v6.7.4-0.20190416172638-348a4f943ff4+incompatible
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/gophercloud/gophercloud v0.7.0 // indirect
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/kolo/xmlrpc v0.0.0-20190909154602-56d5ec7c422e // indirect
	github.com/kr/pty v1.1.3 // indirect
	github.com/labstack/echo/v4 v4.1.11
	github.com/linode/linodego v0.12.2 // indirect
	github.com/liquidweb/liquidweb-go v1.6.1 // indirect
	github.com/logrusorgru/aurora v0.0.0-20191116043053-66b7ad493a23
	github.com/miekg/dns v1.1.22 // indirect
	github.com/oracle/oci-go-sdk v13.0.0+incompatible // indirect
	github.com/ovh/go-ovh v0.0.0-20181109152953-ba5adb4cf014
	github.com/rhysd/go-github-selfupdate v1.1.0
	github.com/sacloud/libsacloud v1.32.0 // indirect
	github.com/transip/gotransip v5.8.2+incompatible // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/urfave/cli v1.22.2
	github.com/valyala/fasttemplate v1.1.0 // indirect
	github.com/vultr/govultr v0.1.7 // indirect
	go.opencensus.io v0.22.2 // indirect
	go.uber.org/ratelimit v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20191202143827-86a70503ff7e
	golang.org/x/net v0.0.0-20191204025024-5ee1b9f4859a
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20191203220235-3fa9dbf08042 // indirect
	google.golang.org/grpc v1.25.1 // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/ns1/ns1-go.v2 v2.0.0-20191126161805-25b9eac84517 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

// broken go.mod inside a package needing this
replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

// Otherwise, the build will fail in dh_golang
//replace google.golang.org/api => google.golang.org/api v0.3.0

go 1.13
