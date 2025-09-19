CREATE TABLE animal_product (
    id UUID PRIMARY KEY,
    animal_id UUID REFERENCES animal(id),
    category_id UUID NOT NULL REFERENCES product_category(id),
    quality_grade_id UUID REFERENCES quality_grade(id),
    quantity NUMERIC NOT NULL,
    unit VARCHAR(20) NOT NULL,
    price_per_unit DECIMAL(10,2) NOT NULL,
    total_value DECIMAL(12,2) GENERATED ALWAYS AS (quantity * price_per_unit) STORED,
    produced_at TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

