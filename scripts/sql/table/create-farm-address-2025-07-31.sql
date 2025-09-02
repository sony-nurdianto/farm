CREATE TABLE addresses (
    id UUID,
    street TEXT,
    village TEXT,
    sub_district TEXT,
    city TEXT,
    province TEXT,
    postal_code TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
) PARTITION BY HASH (id);

CREATE TABLE addresses_p0 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 0);

CREATE TABLE addresses_p1 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 1);

CREATE TABLE addresses_p2 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 2);

CREATE TABLE addresses_p3 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 3);

CREATE TABLE addresses_p4 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 4);

CREATE TABLE addresses_p5 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 5);

CREATE TABLE addresses_p6 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 6);

CREATE TABLE addresses_p7 PARTITION OF addresses
    FOR VALUES WITH (modulus 8, remainder 7);

-- CREATE INDEX idx_addresses_province ON addresses (province);
--
-- -- Index untuk city (sering diquery)
-- CREATE INDEX idx_addresses_city ON addresses (city);
--
-- -- Index untuk postal_code
-- CREATE INDEX idx_addresses_postal_code ON addresses (postal_code);
--
-- -- Composite index untuk geographic queries
-- CREATE INDEX idx_addresses_province_city ON addresses (province, city);
--
-- -- Index untuk sub_district
-- CREATE INDEX idx_addresses_sub_district ON addresses (sub_district);
--
-- -- Index untuk created_at (untuk sorting/filtering)
-- CREATE INDEX idx_addresses_created_at ON addresses (created_at);
