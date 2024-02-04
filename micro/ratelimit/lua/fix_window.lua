local val = redis.call('GET',KEYS[1])
local limiter = tonumber(ARGV[1])
if val == false then
    if limiter < 1 then
        return "true"
    else
        redis.call('SET',KEYS[1],1,'PX',ARGV[2])
        return "false"
    end
elseif tonumber(val)<limiter then
    redis.call('incr',KEYS[1])
    return "false"
else
    return "true"
end