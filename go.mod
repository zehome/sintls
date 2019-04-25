module github.com/zehome/sintls

require (
	cloud.google.com/go v0.37.2 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/Azure/azure-sdk-for-go v26.3.1+incompatible // indirect
	github.com/Azure/go-autorest v11.2.8+incompatible // indirect
	github.com/JamesClonk/vultr v2.0.0+incompatible // indirect
	github.com/OpenDNS/vegadns2client v0.0.0-20180418235048-a3fa4a771d87 // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.7.4 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190423085641-6c7e87673c6e // indirect
	github.com/aws/aws-sdk-go v1.19.15 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/cheynewallace/tabby v1.1.0
	github.com/cloudflare/cloudflare-go v0.9.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20181012123002-c6f51f82210d
	github.com/cpu/goacmedns v0.0.1 // indirect
	github.com/decker502/dnspod-go v0.0.0-20180416134550-83a3ba562b04 // indirect
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/dnsimple/dnsimple-go v0.23.0 // indirect
	github.com/exoscale/egoscale v0.17.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-acme/lego v2.5.0+incompatible
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/go-pg/migrations v6.7.4-0.20190416172638-348a4f943ff4+incompatible
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/google/uuid v1.1.1 // indirect
	github.com/gophercloud/gophercloud v0.0.0-20190424031112-b9b92a825806 // indirect
	github.com/h2non/gock v0.0.0-00010101000000-000000000000 // indirect
	github.com/iij/doapi v0.0.0-20180911005243-8803795a9b7b // indirect
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/juju/ratelimit v1.0.1 // indirect
	github.com/kolo/xmlrpc v0.0.0-20190417161013-de6d879202d7 // indirect
	github.com/labstack/echo/v4 v4.0.0
	github.com/linode/linodego v0.7.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/miekg/dns v1.1.8 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/namedotcom/go v0.0.0-20180403034216-08470befbe04 // indirect
	github.com/nrdcg/auroradns v1.0.0 // indirect
	github.com/nrdcg/goinwx v0.6.0 // indirect
	github.com/oracle/oci-go-sdk v5.4.0+incompatible // indirect
	github.com/ovh/go-ovh v0.0.0-20181109152953-ba5adb4cf014
	github.com/rhysd/go-github-selfupdate v1.1.0
	github.com/sacloud/libsacloud v1.21.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/timewasted/linode v0.0.0-20160829202747-37e84520dcf7 // indirect
	github.com/transip/gotransip v5.8.2+incompatible // indirect
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasttemplate v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	gopkg.in/AlecAivazis/survey.v1 v1.8.3
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/ns1/ns1-go.v2 v2.0.0-20190322154155-0dafb5275fd1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

// broken go.mod inside a package needing this
replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

// Otherwise, the build will fail in dh_golang
replace google.golang.org/api => google.golang.org/api v0.3.0
