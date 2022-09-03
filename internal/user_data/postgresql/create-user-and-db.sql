-- file: 10-create-user-and-db.sql
CREATE DATABASE persons;
CREATE ROLE program WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE persons TO program;
GRANT ALL PRIVILEGES ON DATABASE postgres TO program;
ALTER ROLE program WITH LOGIN;

CREATE TABLE persons (
                         person_id bigserial not null primary key,
                         name text not null,
                         age integer not null,
                         work text,
                         address text not null
);