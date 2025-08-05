CREATE TABLE accounts (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(225) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE (email, id)
) PARTITION BY HASH (id);

CREATE TABLE accounts_p0 PARTITION OF accounts FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE accounts_p1 PARTITION OF accounts FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE accounts_p2 PARTITION OF accounts FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE accounts_p3 PARTITION OF accounts FOR VALUES WITH (MODULUS 4, REMAINDER 3);

CREATE INDEX idx_accounts_p0_email ON accounts_p0 (email);
CREATE INDEX idx_accounts_p1_email ON accounts_p1 (email);
CREATE INDEX idx_accounts_p2_email ON accounts_p2 (email);
CREATE INDEX idx_accounts_p3_email ON accounts_p3 (email);

