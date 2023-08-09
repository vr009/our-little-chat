CREATE TABLE IF NOT EXISTS users
(
    user_id      uuid         NOT NULL PRIMARY KEY,
    nickname     varchar(100) NOT NULL UNIQUE,
    name         varchar(30),
    surname      varchar(30),
    password     varchar(200),
    last_auth    DATE,
    registered   DATE         NOT NULL,
    avatar       varchar
);
