local fiber = require('fiber')
local uuid = require('uuid')

box.cfg{
    listen = 3301,
}

box.schema.space.create('msg_queue',
        {
            if_not_exists = true,
            format = { { 'chat_id', type = 'uuid' },
                       { 'msg_id', type = 'uuid' , unique = true },
                       { 'sender_id', type = 'uuid' },
                       { 'payload', type = 'string' },
                       { 'created_at', type = 'number' } }
        })

box.space.msg_queue:create_index('chat_index', {if_not_exists = true, parts = { { 1, 'uuid' }, {5, 'number'} }})

box.space.msg_queue:create_index('msg_index', {if_not_exists = true, parts = { { 2, 'uuid' }}})

box.space.msg_queue:create_index('sender_index', {if_not_exists = true, unique=false, parts = { { 3, 'uuid' }, {1, 'uuid'}}})

-- This space is needed to keep the info about what for chats have any user
-- and what users take part in chat.
box.schema.space.create('user_chat_list',
        {
            if_not_exists = true,
            format = { { 'participant', type = 'uuid' },
                       { 'chat_id' , type = 'uuid'},
                       { 'last_read', type = 'number' }} -- last message stamp? last_read field?
        })

box.space.user_chat_list:create_index('participant_index', { if_not_exists = true, unique =true, parts = { { 1, 'uuid' }, { 2, 'uuid' }} })
box.space.user_chat_list:create_index('chat_id_index', { if_not_exists = true, unique =true, parts = { { 2, 'uuid' }, { 1, 'uuid' }} })

box.schema.space.create('chats_upd',
        {
            if_not_exists = true,
            format = { { 'chat_id', type = 'uuid', unique = true  },
                       { 'sender_id', type = 'uuid' },
                       { 'payload', type = 'string' },
                       { 'created_at', type = 'number' } }
        })
box.space.chats_upd:create_index('chat_id_index', { if_not_exists = true, unique =true, parts = { { 1, 'uuid' }} })

local queue = {}

--chat_id, sender_id, receiver_id, payload
--c5f0ae14-d06b-4ffd-9bcd-6df01243a9c5, 62391bd9-157c-4513-8e7c-c082e00d2b7e, 61f98c94-de3c-491a-b9c7-1ea214b3ec13
function queue.put(chat_id, sender_id, receiver_id, payload)
    local msg_id = uuid()
    local created_at = fiber.time()

    print('put start')
    print(chat_id, sender_id, receiver_id, payload)
    print('put end')

    -- we put the id of the last message for chat list update
    box.space.chats_upd:replace{
        uuid.fromstr(chat_id),
        uuid.fromstr(sender_id),
        payload,
        created_at,
    }

    local res = box.space.msg_queue:insert{
        uuid.fromstr(chat_id),
        msg_id,
        uuid.fromstr(sender_id),
        payload, created_at}

    if res ~= nil then
        return{chat_id, msg_id:str(), sender_id, payload, created_at}
    end

    return res
end


function queue.take_new_messages_from_space(chat_id, receiver_id)
    local chat_info = box.space.chats_upd.index.chat_id_index:get({ uuid.fromstr(chat_id) })
    local last_updated = nil
    if chat_info == nil then
        box.space.user_chat_list:replace({uuid.fromstr(chat_id), uuid.fromstr(receiver_id), 0})
    else
        last_updated = chat_info[4]
    end

    local since_not_read = 0
    local user_info = box.space.user_chat_list.index.participant_index:get({uuid.fromstr(receiver_id), uuid.fromstr(chat_id)})
    if user_info == nil then
        box.space.user_chat_list:replace({uuid.fromstr(receiver_id), uuid.fromstr(chat_id), 0})
    else
        since_not_read = user_info[3]
    end

    local batch = {}
    for _, tuple in box.space.msg_queue.index.chat_index:pairs({ uuid.fromstr(chat_id) }) do
        if (since_not_read <= tuple[5]) then
            table.insert(batch, { tuple[1]:str(), tuple[2]:str(), tuple[3]:str(), tuple[4], tuple[5] })
            --print(tuple[1]:str(), tuple[2]:str(), tuple[3]:str(), tuple[4], tuple[5])
        end
    end

    box.space.user_chat_list:replace({
        uuid.fromstr(receiver_id),
        uuid.fromstr(chat_id),
        fiber.time()
    })
    return batch
end


function queue.fetch_chat_list_update(user_id, chat_list)
    local batch = {}
    for _, chat_id in ipairs(chat_list) do
        if chat_id ~= nil and user_id ~= nil then
            print('passed 1')
            print(chat_id)
            print(user_id)
            local chat = box.space.chats_upd.index.chat_id_index:get({chat_id[2]})
            print('passed 2')
            local user_info = box.space.user_chat_list.index.participant_index:get({uuid.fromstr(user_id), chat_id[2]})
            print('passed 3')
            if chat[4] > user_info[3] then
                table.insert(batch, { chat[1]:str(), chat[2]:str(), chat[3], chat[4] })
            end
            print('passed 4')
        end
    end
    return batch
end

function queue.fetch_chat_list_update_for_single_user(user_id)
    local chat_list = box.space.user_chat_list:select(uuid.fromstr(user_id))
    return queue.fetch_chat_list_update(user_id, chat_list)
end

function queue.flush_all()
    local batch = {}
    for k, tuple in box.space.msg_queue:pairs() do
        table.insert(batch, { tuple[1]:str(), tuple[2]:str(), tuple[3]:str(), tuple[4]:str(), tuple[5], tuple[6] })
    end
    box.space.msg_queue:truncate()
    return batch
end

function queue.fetch_all(id)
    local chat_id = uuid.fromstr(id)
    local batch = {}
    for k, tuple in box.space.msg_queue:pairs(chat_id) do
        table.insert(batch, { tuple[1]:str(), tuple[2]:str(), tuple[3]:str(), tuple[4]:str(), tuple[5], tuple[6] })
    end
    return batch
end

function queue.fetch_chats()
    local batch = {}
    for k, chat in box.space.user_chat_list:pairs() do
        local sibling_chat = box.space.user_chat_list.index.chat_list_index_inverted:select(chat[2])
        table.insert(batch, {chat[2]:str(), chat[1]:str(), sibling_chat[1][1]:str(), chat[3]})
    end
    box.space.user_chat_list:truncate()
    return batch
end

rawset(_G, 'put', queue.put)
rawset(_G, 'take_msgs', queue.take_new_messages_from_space)
rawset(_G, 'fetch_chats_upd', queue.fetch_chat_list_update_for_single_user)
rawset(_G, 'flush', queue.flush_all)
rawset(_G, 'fetch', queue.fetch_all)
rawset(_G, 'fetch_chats', queue.fetch_chats)

box.once('debug', function() box.schema.user.grant('guest', 'super') end)

box.schema.user.create('test', {password = 'test'})
box.schema.user.grant('test', 'execute', 'universe')
box.schema.user.grant('test', 'read,write', 'space', 'msg_queue')
box.schema.user.grant('test', 'read,write', 'space', 'user_chat_list')

return queue

-- uuid.fromstr('2743f114-524d-4e3b-8ed0-20666e976d39')