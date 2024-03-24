CREATE TABLE warehouses (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(64) NOT NULL,
    is_available    BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMP
);

---- create above / drop below ----

DROP TABLE warehouses;
