CREATE TYPE size AS ENUM ('s', 'm', 'l');

CREATE TABLE goods (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(127) NOT NULL,
    size          size NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMP
);

---- create above / drop below ----

DROP TABLE goods;
DROP TYPE size;