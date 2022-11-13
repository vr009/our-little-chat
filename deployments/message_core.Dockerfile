FROM tarantool/tarantool:latest

COPY core.lua /opt/tarantool

EXPOSE 3301

CMD ["tarantool", "/opt/tarantool/core.lua"]
