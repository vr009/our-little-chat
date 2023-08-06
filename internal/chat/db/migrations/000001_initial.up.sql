CREATE SCHEMA IF NOT EXISTS chats;

CREATE TABLE IF NOT EXISTS messages
(
    msg_id     uuid         NOT NULL PRIMARY KEY,
    chat_id    uuid         NOT NULL,
    sender_id  uuid         NOT NULL,
    payload    varchar      NOT NULL,
    created_at bigint         NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_participants
(
    chat_id          uuid         NOT NULL,
    participant_id   uuid         NOT NULL,
    PRIMARY KEY (chat_id, participant_id)
);

CREATE TABLE IF NOT EXISTS chats
(
    chat_id        uuid         NOT NULL PRIMARY KEY,
    photo_url      varchar,
    created_at     bigint,
    last_msg_id    uuid         REFERENCES messages(msg_id) DEFAULT NULL
);
