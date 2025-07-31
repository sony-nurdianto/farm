CREATE TYPE farm_type_enum AS ENUM ('CROPLAND', 'ORCHARD', 'RANCH', 'MIXED', 'OTHER');
CREATE TYPE farm_status_enum AS ENUM ('ACTIVE', 'INACTIVE', 'SOLD', 'DESERTED');

CREATE TABLE farms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    farmer_id UUID NOT NULL,
    farm_name VARCHAR(225) NOT NULL,
    farm_type farm_type_enum NOT NULL DEFAULT 'CROPLAND',
    farm_size NUMERIC(6, 2) NOT NULL,
    photo_url TEXT,
    farm_status farm_status_enum NOT NULL DEFAULT 'ACTIVE',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    address_id UUID NOT NULL,
    CONSTRAINT fk_address FOREIGN KEY (address_id) REFERENCES addresses(id),
);

