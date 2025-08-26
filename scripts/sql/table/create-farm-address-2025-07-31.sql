CREATE TABLE addresses (
    id UUID PRIMARY KEY,
    street TEXT,
    village TEXT,
    sub_district TEXT,
    city TEXT,
    province TEXT,
    postal_code TEXT
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
);
