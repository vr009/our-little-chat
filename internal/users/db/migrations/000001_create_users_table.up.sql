CREATE TABLE IF NOT EXISTS users
(
    user_id      uuid         NOT NULL PRIMARY KEY,
    nickname     text         NOT NULL UNIQUE,
    user_name         text,
    surname      text,
    password     bytea        NOT NULL,
    registered   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    avatar       text,
    activated    bool         NOT NULL
);
