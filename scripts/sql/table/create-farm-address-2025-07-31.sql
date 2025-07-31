CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    street TEXT,
    village TEXT,
    sub_district TEXT,
    city TEXT,
    province TEXT,
    postal_code TEXT
);
