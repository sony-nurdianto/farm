CREATE TABLE farms (
    id UUID,
    farmer_id UUID NOT NULL,
    farm_name VARCHAR(225) NOT NULL,
    farm_type VARCHAR NOT NULL DEFAULT 'CROPLAND',
    farm_size NUMERIC(6, 2) NOT NULL,
    photo_url TEXT,
    farm_status VARCHAR NOT NULL DEFAULT 'ACTIVE',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    address_id UUID NOT NULL,
    PRIMARY KEY (id), 
    CONSTRAINT fk_address FOREIGN KEY (address_id) REFERENCES addresses(id)
) PARTITION BY HASH (id);


CREATE TABLE farms_p0 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 0);

CREATE TABLE farms_p1 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 1);

CREATE TABLE farms_p2 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 2);

CREATE TABLE farms_p3 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 3);

CREATE TABLE farms_p4 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 4);

CREATE TABLE farms_p5 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 5);

CREATE TABLE farms_p6 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 6);

CREATE TABLE farms_p7 PARTITION OF farms
    FOR VALUES WITH (modulus 8, remainder 7);


CREATE INDEX idx_farms_farmer_id ON farms (farmer_id);

-- CREATE INDEX idx_farms_status ON farms (farm_status);
--
-- CREATE INDEX idx_farms_created_at ON farms (created_at);
--
-- CREATE INDEX idx_farms_type ON farms (farm_type);
--
-- CREATE INDEX idx_farms_farmer_status ON farms (farmer_id, farm_status);

