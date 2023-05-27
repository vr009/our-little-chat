local fiber = require('fiber')
local uuid = require('uuid')
local log = require('log')
log.cfg{ level='debug', log='tarantool.log'}

box.cfg{
    listen = 3301,
}

box.schema.space.create('messages',
        {
            if_not_exists = true,
            format = { { 'chat_id', type = 'uuid' },
                       { 'msg_id', type = 'uuid' , unique = true },
                       { 'sender_id', type = 'uuid' },
                       { 'payload', type = 'string' },
                       { 'created_at', type = 'number' } }
        })

box.space.messages:create_index('chat_index', 
        {if_not_exists = true, parts = { { 1, 'uuid' }, {5, 'number'} }})

box.space.messages:create_index('msg_index', 
        {if_not_exists = true, parts = { { 2, 'uuid' }}})

box.space.messages:create_index('sender_index', 
        {if_not_exists = true, unique=false, parts = { { 3, 'uuid' }, {1, 'uuid'}}})

-- This space is needed to keep the info about what for chats have any user
-- and what users take part in chat. Also here is info about the last time when user visited the chat.
box.schema.space.create('chat_participants',
        {
            if_not_exists = true,
            format = { { 'participant_id', type = 'uuid' },
                       { 'chat_id' , type = 'uuid'},
                       { 'last_read_msg_id', type = 'uuid' }}
        })

box.space.chat_participants:create_index('participant_index', 
        { if_not_exists = true, unique =true, parts = { { 1, 'uuid' }, { 2, 'uuid' }} })

box.space.chat_participants:create_index('chat_id_index', 
        { if_not_exists = true, unique =true, parts = { { 2, 'uuid' }, { 1, 'uuid' }} })

-- contains last messages in chat
box.schema.space.create('chat_last_msgs',
        {
            if_not_exists = true,
            format = { { 'chat_id', type = 'uuid', unique = true  },
                       { 'sender_id', type = 'uuid' },
                       { 'msg_id', type = 'uuid' },
                       { 'payload', type = 'string' },
                       { 'created_at', type = 'number' } }
        })
box.space.chat_last_msgs:create_index('chat_id_index',
        { if_not_exists = true, unique =true, parts = { { 1, 'uuid' }} })

-- contains unread messages of a chat
box.schema.space.create('unread_msgs',
        {
            if_not_exists = true,
            format = { { 'chat_id', type = 'uuid', unique = true },
                       { 'msg_id', type = 'uuid' },
                       { 'sender_id', type = 'uuid'}}
        })
box.space.unread_msgs:create_index('chat_id_index',
        { if_not_exists = true, unique =true, parts = { { 1, 'uuid' }} })

local queue = {}

function queue.create_chat(users, chat_id)
    if chat_id == nil then
        chat_id = uuid.str()
    end

    local chat = box.space.chat_participants.index.chat_id_index:select({ uuid.fromstr(chat_id) })
    if chat[1] ~= nil then
        log.info('chat exists')
        return
    end

    for _, user_id in pairs(users) do
        uuid.fromstr(user_id)
        box.space.chat_participants:insert({
            uuid.fromstr(user_id),
            uuid.fromstr(chat_id),
            nil,
        })
    end
    return {chat_id}
end

function queue.add_user_to_chat(chat_id, user_id)
    box.space.chat_participants:replace({user_id, chat_id, nil})
end

function queue.get_members(chat_id)
    local batch = {}
    local usrs = box.space.chat_participants.index.chat_id_index:select({uuid.fromstr(chat_id)})
    for _, tuple in pairs(usrs) do
        table.insert(batch, { tuple['participant_id']:str() })
    end
    return batch
end

--chat_id, sender_id, receiver_id, payload
--c5f0ae14-d06b-4ffd-9bcd-6df01243a9c5, 62391bd9-157c-4513-8e7c-c082e00d2b7e, 61f98c94-de3c-491a-b9c7-1ea214b3ec13
function queue.put(chat_id, sender_id, payload)
    local res_chat_id = box.space.chat_participants.index.participant_index:select({ uuid.fromstr(sender_id), uuid.fromstr(chat_id) })[1]
    log.info(res_chat_id)
    if res_chat_id == nil then
        log.info('doesnt exist')
        return
    end

    local msg_id = uuid.new()
    local created_at = fiber.time()

    -- we put the id of the last message for chat list update
    box.space.chat_last_msgs:replace{
        uuid.fromstr(chat_id),
        uuid.fromstr(sender_id),
        msg_id,
        payload,
        created_at,
    }

    -- replace an unread message in the chat
    box.space.unread_msgs:replace{
        uuid.fromstr(chat_id),
        msg_id,
        uuid.fromstr(sender_id)
    }

    local res = box.space.messages:insert{
        uuid.fromstr(chat_id),
        msg_id,
        uuid.fromstr(sender_id),
        payload, created_at}

    if res ~= nil then
        log.error({ 'failed to put a message', chat_id, msg_id:str(), sender_id, payload, created_at})
        return{chat_id, msg_id:str(), sender_id, payload, created_at}
    end
    return res
end


function queue.take_new_messages_from_space(chat_id, receiver_id)
    local since_not_read = 0
    -- to know the last update of the chat
    local user_info = box.space.chat_participants.index.participant_index:get({uuid.fromstr(receiver_id), uuid.fromstr(chat_id)})
    if user_info == nil then
        return {}
    else
        local msg_since_not_read_id = user_info['last_read_msg_id']
        if msg_since_not_read_id == nil then
            since_not_read = 0
        else
            since_not_read = bos.space.messages.index.msg_index:get(msg_since_not_read_id)
        end
    end

    local last_message_id
    -- collect messages since the since_not_read number
    local batch = {}
    for _, message in box.space.messages.index.chat_index:pairs({ uuid.fromstr(chat_id) }) do
        if (since_not_read <= message['created_at']) then
            table.insert(batch, {
                message['chat_id']:str(),
                message['msg_id']:str(),
                message['sender_id']:str(),
                message['payload'],
                message['created_at'],
            })
            log.info('found! : ', message[1]:str(), message[2]:str(), message[3]:str(), message[4], message[5])
            last_message_id = message['msg_id']
        end
    end

    -- delete message from unread messages
    box.space.unread_msgs:delete{
        uuid.fromstr(chat_id),
    }

    -- update the info about a user that has read the messages from chat
    box.space.chat_participants:replace({
        uuid.fromstr(receiver_id),
        uuid.fromstr(chat_id),
        last_message_id,
    })
    return batch
end

function queue.fetch_last_chat_messages_for_all_users()
    local batch = {}
    for _, tuple in box.space.chat_last_msgs:pairs() do
        table.insert(batch, {
            tuple['chat_id']:str(),
            tuple['sender_id']:str(),
            tuple['msg_id']:str(),
            tuple['payload'],
            tuple['created_at'],
        })
    end
    return batch
end

function queue.flush_all_msgs()
    local batch = {}
    for _, tuple in box.space.messages:pairs() do
        table.insert(batch, {
            tuple['chat_id']:str(),
            tuple['msg_id']:str(),
            tuple['sender_id']:str(),
            tuple['payload'],
            tuple['created_at'],
        })
    end
    box.space.messages:truncate()
    return batch
end

function queue.fetch_all_user_msgs(id)
    local chat_id = uuid.fromstr(id)
    local batch = {}
    for k, tuple in box.space.messages:pairs(chat_id) do
        table.insert(batch, {
            tuple['chat_id']:str(),
            tuple['msg_id']:str(),
            tuple['sender_id']:str(),
            tuple['payload'],
            tuple['created_at'],
        })
    end
    return batch
end

function queue.flush_chats_participants()
    local batch = {}
    for _, chat in box.space.chat_participants:pairs() do
        table.insert(batch, {
            chat['participant_id']:str(),
            chat['chat_id']:str(),
            chat['last_read_msg_id'],
        })
    end
    box.space.chat_participants:truncate()
    return batch
end

-- this function returns all unread messages for a user
function queue.fetch_unread_messages(user_id)
    local chat_ids = box.space.chat_participants.index.participant_index:select({
        uuid.fromstr(user_id)})
    local batch = {}
    for _, chat_id in pairs(chat_ids) do
        local msg_id = box.space.unread_msgs.index.chat_id_index:get{chat_id}
        local message = bos.space.messages.index.msg_index:get(msg_id)
        if message['sender_id']:str() ~= user_id then
            table.insert(batch, {
                message['chat_id']:str(),
                message['sender_id']:str(),
                message['msg_id']:str(),
                message['payload'],
                message['created_at'],
            })
        end
    end
    return batch
end

-- create_chat creates a chat for passed users, the number of users is (seems) not limited
rawset(_G, 'create_chat', queue.create_chat)
-- add_user_to_chat just adds a user to the members of chat
rawset(_G, 'add_user_to_chat', queue.add_user_to_chat)

-- returns all members of the chat
rawset(_G, 'get_members_of_chat', queue.get_members)

-- a method for posting a message to the chat
rawset(_G, 'put', queue.put)
-- get all unread messages
rawset(_G, 'take_msgs', queue.take_new_messages_from_space)

-- returns a list of last messages of all chats
rawset(_G, 'fetch_all_chat_last_msgs', queue.fetch_last_chat_messages_for_all_users)

-- returns a list of unread messages
rawset(_G, 'fetch_unread_messages', queue.fetch_unread_messages)

rawset(_G, 'flush_msgs', queue.flush_all_msgs)
rawset(_G, 'flush_chats_participants', queue.flush_chats_participants)
rawset(_G, 'fetch_msgs', queue.fetch_all_user_msgs)

box.once('debug', function() box.schema.user.grant('guest', 'super') end)

box.schema.user.create('test', {password = 'test', if_not_exists = true})
box.schema.user.grant('test', 'execute', 'universe', nil, {if_not_exists=true})
box.schema.user.grant('test', 'read,write', 'space', 'messages', {if_not_exists=true})
box.schema.user.grant('test', 'read,write', 'space', 'chat_participants', {if_not_exists=true})
box.schema.user.grant('test', 'read,write', 'space', 'chat_last_msgs', {if_not_exists=true})
box.schema.user.grant('test', 'read,write', 'space', 'unread_msgs', {if_not_exists=true})

return queue

-- uuid.fromstr('2743f114-524d-4e3b-8ed0-20666e976d39')