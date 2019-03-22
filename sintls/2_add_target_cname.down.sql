ALTER TABLE sintls_host
DROP CONSTRAINT sintls_host_check,
ADD CONSTRAINT sintls_host_check CHECK (dns_target_a is not null or dns_target_aaaa is not null),
DROP COLUMN sintls_target_cname;