CREATE TABLE farms (
    id UUID PRIMARY KEY,
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
    CONSTRAINT fk_address FOREIGN KEY (address_id) REFERENCES addresses(id)
);

