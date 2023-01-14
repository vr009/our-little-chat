FROM tarantool/tarantool:latest

COPY internal/message_queue/core.lua /opt/tarantool

EXPOSE 3301

CMD ["tarantool", "/opt/tarantool/core.lua"]
