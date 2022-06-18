local tap = require('tap')
local uuid = require('uuid')
local queue = require('queue')


local console_1 = require('console')
local console_2 = require('console')

local TARANTOOL_PATH = arg[-1]
local test = tap.test('core-test')
local cmd = 'ERRINJ_STDIN_ISATTY=1 ' .. TARANTOOL_PATH .. ' -i 2>&1'

test:plan(1)

local user_a = uuid()
local user_b = uuid()
local ab_chat_id = uuid()

local message_1 = 'Hi'
local message_2 = 'Hi there'



