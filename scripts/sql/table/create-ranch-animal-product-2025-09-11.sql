CREATE TABLE animal_product (
    id UUID PRIMARY KEY,
    animal_id UUID,
    category_id UUID NOT NULL,
    quality_grade_id UUID,
    quantity NUMERIC NOT NULL,
    unit VARCHAR(20) NOT NULL,
    price_per_unit DECIMAL(10,2) NOT NULL,
    total_value DECIMAL(12,2) GENERATED ALWAYS AS (quantity * price_per_unit) STORED,
    produced_at TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY HASH (id);

CREATE TABLE animal_product_p0 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 0);

CREATE TABLE animal_product_p1 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 1);

CREATE TABLE animal_product_p2 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 2);

CREATE TABLE animal_product_p3 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 3);

CREATE TABLE animal_product_p4 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 4);

CREATE TABLE animal_product_p5 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 5);

CREATE TABLE animal_product_p6 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 6);

CREATE TABLE animal_product_p7 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 7);

CREATE TABLE animal_product_p8 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 8);

CREATE TABLE animal_product_p9 PARTITION OF animal_product
    FOR VALUES WITH (modulus 10, remainder 9);


CREATE INDEX idx_animal_product_animal_id ON animal_product (animal_id);
CREATE INDEX idx_animal_product_category_id ON animal_product (category_id);
CREATE INDEX idx_animal_product_produced_at ON animal_product (produced_at);
CREATE INDEX idx_animal_product_created_at ON animal_product (created_at);


