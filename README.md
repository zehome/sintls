Simple INternal TLS
===================

This projects has a goal to provide let's encrypt x509 certificates to
hosts not reachable from the internet, with limited internet access.

SinTLS has two parts, both using Lego:
  - Server (Handle DNS updates and ACME-DNS)
  - Client (Asks and manage certificates)

Quickstart
==========

Server
------

```shell
# Install PostgreSQL database
apt install postgresql-11

# Create postgresql credentials
sudo -Hu postgres createuser sintls
sudo -Hu postgres createdb sintls -O sintls

# Install sintls-server
dpkg -i sintls-server*.deb

# Setup OVH API access credentials
# This setup is for "clarilab.fr" domain
curl -XPOST https://eu.api.ovh.com/1.0/auth/credential \
  -H"X-Ovh-Application: ***********" \
  -H "Content-type: application/json" \
  -d @- << EOF
{"accessRules": [
  {"method": "POST", "path": "/domain/zone/clarilab.fr/record/*"},
  {"method": "GET", "path": "/domain/zone/clarilab.fr/record/*" },
  {"method": "DELETE", "path": "/domain/zone/clarilab.fr/record/*"},
  {"method": "POST", "path": "/domain/zone/clarilab.fr/record"},
  {"method": "GET", "path": "/domain/zone/clarilab.fr/record"}
]}
EOF
# Validate the credentials on the web interface as this command tells you

# Copy your API informations to /etc/sintls/sintls-server.conf

# Start sintls socket listener using SystemD
systemctl enable sintls-server.socket

# Create a new user
sudo -Hu sintls sintls-server adduser
```

Sintls-server is configured to listen on port 443, using systemd socket activation.

If SINTLS_ENABLETLS is set in /etc/sintls/sintls-server.conf, then golang's
acme/autocert is used to retrieve a valid TLS certificate on first request.


Client
------

```shell
dpkg -i sintls-client*.deb

# Use your credentials to get a new certificate
SINTLS_USERNAME=**** SINTLS_PASSWORD=**** sintls-client --email xxx@xxx.fr -d ***.subdomain.fr \
  --ca-server https://acme-staging-v02.api.letsencrypt.org/directory \
  --target-a 10.31.254.20 run

# use the renew command to renew instead of create

# Your certificates are stored in ~/.config/sintls/certificates/
```


Overview
========

```text

                                                                  +-----------------------+-------------+
                                                                  | subdomain             | credentials |
                                                                  +-------------------------------------+
                                              +-----------+       | mca.abo.c.clarilab.fr | abocred     |
                                              | PostgreSQL| ----> |                       |             |
                                              +-----+-----+       +-----------------------+-------------+
                                                    ^
                                                    |
                                                    |
                                                    |
+----------------+ <--- Lego / ACME-httpreq ---> +--+-------------+ <--- Lego / ACME auth ---> +-------+
|internal server |                               |auth.clarilab.fr|                            |OVH DNS|
+----------------+ +---   Ask CERTIFICAT    ---> +----------------+ +--- Create A records ---> +-------+

```

