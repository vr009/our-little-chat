CREATE SCHEMA IF NOT EXISTS chats;

CREATE TABLE IF NOT EXISTS chats
(
    chat_id        uuid         NOT NULL PRIMARY KEY,
    photo_url      varchar,
    created_at     bigint
);

CREATE TABLE IF NOT EXISTS messages
(
    msg_id     uuid         NOT NULL PRIMARY KEY,
    chat_id    uuid         NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    sender_id  uuid         NOT NULL,
    payload    varchar      NOT NULL,
    created_at bigint       NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_participants
(
    chat_id          uuid         NOT NULL,
    participant_id   uuid         NOT NULL,
    chat_name        varchar      NOT NULL,
    PRIMARY KEY (chat_id, participant_id)
);

