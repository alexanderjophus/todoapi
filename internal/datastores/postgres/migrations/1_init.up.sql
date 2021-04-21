CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE todo
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
);