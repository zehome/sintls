module github.com/zehome/sintls

require (
	cloud.google.com/go v0.37.2 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/cheynewallace/tabby v1.1.0
	github.com/coreos/go-systemd v0.0.0-20181012123002-c6f51f82210d
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/go-acme/lego/v3 v3.0.2
	github.com/go-pg/migrations v6.7.4-0.20190416172638-348a4f943ff4+incompatible
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/labstack/echo/v4 v4.0.0
	github.com/logrusorgru/aurora v0.0.0-20190428105938-cea283e61946
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ovh/go-ovh v0.0.0-20181109152953-ba5adb4cf014
	github.com/rhysd/go-github-selfupdate v1.1.0
	github.com/transip/gotransip v5.8.2+incompatible // indirect
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasttemplate v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c
	gopkg.in/AlecAivazis/survey.v1 v1.8.5
	mellium.im/sasl v0.2.1 // indirect
)

// broken go.mod inside a package needing this
replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

// Otherwise, the build will fail in dh_golang
replace google.golang.org/api => google.golang.org/api v0.3.0

// git.apache.org does not exists
replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
