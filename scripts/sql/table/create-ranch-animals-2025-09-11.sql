CREATE TABLE animal (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farm(id),
    name VARCHAR(100) NOT NULL,
    species VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    origin_type TEXT CHECK(origin_type IN ('born', 'hatched', 'purchased')),
    origin_from UUID NULL REFERENCES animal(id),
    gender CHAR(1) CHECK (gender IN ('M', 'F')),
    weight_kg DECIMAL(8,2) NULL,
    health_status TEXT CHECK (health_status IN ('healthy','sick','recovered','dead','sold')) DEFAULT 'healthy',
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT
);
