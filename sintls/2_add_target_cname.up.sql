ALTER TABLE sintls_host ADD COLUMN dns_target_cname text,
DROP CONSTRAINT sintls_host_check,
ADD CONSTRAINT sintls_host_check CHECK (dns_target_a is not null or dns_target_aaaa is not null or dns_target_cname is not null)
