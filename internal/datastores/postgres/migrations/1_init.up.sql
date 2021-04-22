CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE todo
(
    id          UUID    PRIMARY KEY DEFAULT uuid_generate_v4(),
    description TEXT    NOT NULL,
    done        BOOLEAN NOT NULL DEFAULT false
);