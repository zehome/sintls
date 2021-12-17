### Simple INternal TLS
This projects has a goal to provide let's encrypt x509 certificates to
hosts not reachable from the internet, with limited internet access.

SinTLS has two parts, both using Lego:
  - Server (Handle DNS updates and ACME-DNS)
  - Client (Asks and manage certificates)


Supported DNS providers:
  - rfc2136
  - ovh

# Quickstart

## Server

```shell
# Install PostgreSQL database
apt install postgresql-14

# Create postgresql credentials
sudo -Hu postgres createuser sintls
sudo -Hu postgres createdb sintls -O sintls

# Install sintls-server
dpkg -i sintls-server*.deb
```

### RFC 2136 (bind RNDC)
On the bind server
```shell
rndc-confgen -b 256 -k sintls
# Then adapt your bind configuration and reload bind9 service
```
On the SinTLS server, adjust /etc/sintls/sintls-server.conf. All paramters must match (including the key-name)

### OVH DNS provider

# Setup OVH API access credentials
# This setup is for "zehome.com" domain
curl -XPOST https://eu.api.ovh.com/1.0/auth/credential \
  -H"X-Ovh-Application: ***********" \
  -H "Content-type: application/json" \
  -d @- << EOF
{"accessRules": [
  {"method": "POST", "path": "/domain/zone/zehome.com/record/*"},
  {"method": "GET", "path": "/domain/zone/zehome.com/record/*" },
  {"method": "DELETE", "path": "/domain/zone/zehome.com/record/*"},
  {"method": "POST", "path": "/domain/zone/zehome.com/record"},
  {"method": "POST", "path": "/domain/zone/zehome.com/refresh"},
  {"method": "GET", "path": "/domain/zone/zehome.com/record"}
]}
EOF
# Validate the credentials on the web interface as this command tells you
# Copy your API informations to /etc/sintls/sintls-server.conf


## Run SinTLS server
Start sintls socket listener using SystemD
```
systemctl enable sintls-server.socket
```

## Create a new user
```
sudo -Hu sintls sintls-server adduser
```

Sintls-server is configured to listen on port 443, using systemd socket activation.

If SINTLS_ENABLETLS is set in /etc/sintls/sintls-server.conf, then golang's
acme/autocert is used to retrieve a valid TLS certificate on first request.


# Client

```shell
dpkg -i sintls-client*.deb

# Use your credentials to get a new certificate
SINTLS_USERNAME=**** SINTLS_PASSWORD=**** sintls-client --email xxx@xxx.fr -d ***.subdomain.fr \
  --ca-server https://acme-staging-v02.api.letsencrypt.org/directory \
  --target-a 10.31.254.20 run

# use the renew command to renew instead of create

# Your certificates are stored in ~/.config/sintls/certificates/
```


# Overview

```text

                                                                  +--------------------------+-------------+
                                                                  | subdomain                | credentials |
                                                                  +-------------------------------------+
                                              +-----------+       | laurent.dev.bluemind.net | abocred     |
                                              | PostgreSQL| ----> |                          |             |
                                              +-----+-----+       +--------------------------+-------------+
                                                    ^
                                                    |
                                                    |
                                                    |
+----------------+ <--- Lego / ACME-httpreq ---> +--+-------------+ <--- Lego / ACME auth ---> +--------------+
|internal server |                               |dev.bluemind.net|                            | DNS provider |
+----------------+ +---   Ask CERTIFICAT    ---> +----------------+ +--- Create A records ---> +--------------+

```

