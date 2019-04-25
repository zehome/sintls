Simple INternal TLS
===================

[![Build Status](https://travis-ci.org/zehome/sintls.svg?branch=master)](https://travis-ci.org/zehome/sintls)

This projects has a goal to provide let's encrypt x509 certificates to
hosts not reachable from the internet, with limited internet access.

SinTLS has two parts, both using Lego:
  - Server (Handle DNS updates and ACME-DNS)
  - Client (Asks and manage certificates)

Quickstart
==========

OVH credentials setup
```shell
curl -XPOST -H"X-Ovh-Application: ***********" -H "Content-type: application/json" https://eu.api.ovh.com/1.0/auth/credential  -d '{"accessRules": [{"method": "POST", "path": "/domain/zone/clarilab.fr/record/*"}, { "method": "GET", "path": "/domain/zone/clarilab.fr/record/*" }, {"method": "DELETE", "path": "/domain/zone/clarilab.fr/record/*"}, {"method": "POST", "path": "/domain/zone/clarilab.fr/record"}]}'
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

