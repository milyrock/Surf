BEGIN;

DROP TABLE IF EXISTS surf.record_category;
CREATE TYPE surf.record_category AS ENUM (
    'damage', 'expiration', 'lunch', 'additional'
);

DROP TABLE IF EXISTS surf.records;
create table surf.records (
    id SERIAL PRIMARY KEY,
    username varchar(255) not null,
    product  varchar(255)  not null,
    category surf.record_category not null,
    amount   integer        NOT NULL CHECK (amount BETWEEN 0 AND 20) DEFAULT 1
); 

COMMIT;