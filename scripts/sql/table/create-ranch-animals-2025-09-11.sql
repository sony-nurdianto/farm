-- Create partitioned table using hash partitioning
CREATE TABLE animal (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    origin_type TEXT CHECK(origin_type IN ('born', 'hatched', 'purchased')),
    origin_from UUID NULL,
    purchased_price DECIMAL(12,2) NULL,
    gender CHAR(1) CHECK (gender IN ('M', 'F')),
    weight_kg DECIMAL(8,2) NULL,
    health_status TEXT CHECK (health_status IN ('healthy','sick','recovered','dead','sold')) DEFAULT 'healthy',
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT
) PARTITION BY HASH (id);

-- Create partitions (10 partisi)
CREATE TABLE animal_p0 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 0);

CREATE TABLE animal_p1 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 1);

CREATE TABLE animal_p2 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 2);

CREATE TABLE animal_p3 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 3);

CREATE TABLE animal_p4 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 4);

CREATE TABLE animal_p5 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 5);

CREATE TABLE animal_p6 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 6);

CREATE TABLE animal_p7 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 7);

CREATE TABLE animal_p8 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 8);

CREATE TABLE animal_p9 PARTITION OF animal
    FOR VALUES WITH (modulus 10, remainder 9);

-- Alternative: Jika ingin partition berdasarkan farm_id untuk better data locality
-- CREATE TABLE animal (
--     id UUID PRIMARY KEY,
--     farm_id UUID NOT NULL REFERENCES farm(id),
--     name VARCHAR(100) NOT NULL,
--     species VARCHAR(100) NOT NULL,
--     breed VARCHAR(100),
--     birth_date DATE,
--     origin_type TEXT CHECK(origin_type IN ('born', 'hatched', 'purchased')),
--     origin_from UUID NULL REFERENCES animal(id),
--     gender CHAR(1) CHECK (gender IN ('M', 'F')),
--     weight_kg DECIMAL(8,2) NULL,
--     health_status TEXT CHECK (health_status IN ('healthy','sick','recovered','dead','sold')) DEFAULT 'healthy',
--     registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     is_active BOOLEAN DEFAULT TRUE,
--     notes TEXT
-- ) PARTITION BY HASH (farm_id);

-- Untuk farm_id partitioning, gunakan 10 partisi juga:
-- CREATE TABLE animal_p0 PARTITION OF animal FOR VALUES WITH (modulus 10, remainder 0);
-- CREATE TABLE animal_p1 PARTITION OF animal FOR VALUES WITH (modulus 10, remainder 1);
-- ... dst sampai remainder 9

-- Create indexes on partitioned table
CREATE INDEX idx_animal_farm_id ON animal (farm_id);
CREATE INDEX idx_animal_species ON animal (species);
CREATE INDEX idx_animal_health_status ON animal (health_status);
CREATE INDEX idx_animal_registration_date ON animal (registration_date);

-- Untuk melihat partisi yang ada
-- SELECT schemaname, tablename, partitionbounds 
-- FROM pg_tables 
-- WHERE tablename LIKE 'animal_%';
