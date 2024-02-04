local var = redis.call('GET',KEYS[1])
if var == false then
    return redis.call('SET',KEYS[1],ARGV[1],'EX',ARGV[2])
elseif var == ARGV[1] then
    local exp = redis.call("EXPIRE",KEYS[1],ARGV[2])
    if exp == 1 then
        return 'OK'
    else
        return ''
    end
else
    return ''
end
