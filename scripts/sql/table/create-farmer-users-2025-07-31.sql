CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    full_name VARCHAR(225) NOT NULL,
    email VARCHAR(225) NOT NULL,
    phone VARCHAR(225) NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    last_login TIMESTAMP WITH TIME ZONE,
    verified BOOLEAN DEFAULT FALSE,
    profile_photo_url TEXT,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE (id, email)
) PARTITION BY HASH(id);

CREATE TABLE users_p0 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE users_p1 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE users_p2 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE users_p3 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 3);

CREATE INDEX idx_users_p0_email ON users_p0 (email);
CREATE INDEX idx_users_p1_email ON users_p1 (email);
CREATE INDEX idx_users_p2_email ON users_p2 (email);
CREATE INDEX idx_users_p3_email ON users_p3 (email);

