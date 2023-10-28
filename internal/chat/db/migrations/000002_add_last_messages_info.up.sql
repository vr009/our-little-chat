ALTER TABLE chats
    ADD COLUMN IF NOT EXISTS last_msg_id uuid DEFAULT NULL REFERENCES messages(msg_id);