INSERT INTO chat_participants (chat_id, participant_id) VALUES
      ('ddd37de5-6799-4158-9703-e536fb8071ac', '43b88a7a-f3ef-4631-be98-a2ed31137e8e'),
      ('ddd37de5-6799-4158-9703-e536fb8071ac', '7836098a-a55f-40d5-b9a7-1db75d1f68db');


INSERT INTO messages (msg_id, chat_id, sender_id, payload, created_at) VALUES (
    'b6e878eb-385c-4a30-a042-f7cc94d5c8de',
    'ddd37de5-6799-4158-9703-e536fb8071ac',
    '43b88a7a-f3ef-4631-be98-a2ed31137e8e',
    'Hello world!', 0.00);

INSERT INTO chats (chat_id, last_msg_id) VALUES
      ('ddd37de5-6799-4158-9703-e536fb8071ac','b6e878eb-385c-4a30-a042-f7cc94d5c8de');
