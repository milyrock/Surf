BEGIN;

CREATE SCHEMA IF NOT EXISTS surf;

DROP TYPE IF EXISTS surf.record_category CASCADE;
CREATE TYPE surf.record_category AS ENUM (
    'порча', 'срок_годности', 'ланч', 'еда'
);

DROP TABLE IF EXISTS surf.records;
CREATE TABLE surf.records (
    id SERIAL PRIMARY KEY,
    username varchar(255) not null,
    product  varchar(255)  not null,
    category surf.record_category not null,
    amount   integer        NOT NULL CHECK (amount BETWEEN 0 AND 20) DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
); 

COMMIT;