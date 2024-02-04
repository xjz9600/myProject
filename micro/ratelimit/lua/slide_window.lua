local now = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])
local limiter = tonumber(ARGV[3])
local min = now-interval
redis.call('ZREMRANGEBYSCORE',KEYS[1],'-inf',min)
local cnt = redis.call('ZCOUNT',KEYS[1],'-inf','+inf')
if cnt >= limiter then
    return "true"
else
    redis.call('ZADD',KEYS[1],now,now)
    redis.call('PEXPIRE',KEYS[1],interval)
    return "false"
end