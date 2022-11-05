CREATE SCHEMA IF NOT EXISTS person;

CREATE TABLE IF NOT EXISTS users
(
    user_id     uuid                     NOT NULL PRIMARY KEY,
    nickname    varchar(100)             NOT NULL UNIQUE,
    password    varchar,
    last_auth   DATE,
    registered  DATE                     NOT NULL,
    avatar      TEXT[],
    contact_list TEXT[]
);
