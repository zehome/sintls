Simple INternal TLS
===================

This projects has a goal to provide let's encrypt x509 certificates to
hosts not reachable from the internet, with limited internet access.

SinTLS has two parts, both using Lego:
  - Server (Handle DNS updates and ACME-DNS)
  - Client (Asks and manage certificates)

Quickstart
==========



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

