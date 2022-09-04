-- file: 10-create-user-and-db.sql
CREATE DATABASE persons;
CREATE ROLE program WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE persons TO program;
GRANT ALL PRIVILEGES ON DATABASE postgres TO program;
ALTER ROLE program WITH LOGIN;

CREATE TABLE "persons"
(
    UserID   uuid default gen_random_uuid() NOT NULL
        PRIMARY KEY ,
    Nickname  varchar(100)                   NOT NULL ,
    LastAuth   DATE,
    Registered   DATE  NOT NULL,
    Avatar   TEXT[],
    ContactList  TEXT[]
);

ALTER TABLE "user"
    OWNER TO postgres;

