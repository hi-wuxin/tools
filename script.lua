-- 移除带分值的集合参数
---- 移除不带分值的集合参数
--function del(result)
--    for i, member in ipairs(result) do
--        redis.call('ZREM',key,member)
--    end
--end
local key=KEYS[1] local val=ARGV[1] local num=ARGV[2] local n=redis.call("GET",key) if n==false then redis.call("IncrBy",key,val) return 1 end if tonumber(n)>=num then return 1 end redis.call("IncrBy",key,val) return 1


--local key=KEYS[1] local start=ARGV[1] local stop=ARGV[2] local result=redis.call('ZRANGE',key,'0','10','WITHSCORES') if result==false then return {} end for i, member in ipairs(result) do if i%2~=0  then redis.call('ZREM',key,member) end end return result
--
-- zrange 22 10 29
--local key=KEYS[1] local start=ARGV[1] local stop=ARGV[2] local result=redis.call('ZRANGE',key,start,stop) if result==false then return {} end for i, member in ipairs(result) do redis.call('ZREM',key,member) end return result
--
---- zrangebyscore 22 1 100 withscores limit 10 100
--local key=KEYS[1] local min=ARGV[1] local max=ARGV[2] local offset=ARGV[3] local limit=ARGV[4] local result=redis.call('zrangebyscore',key,min,max,'withscores',"limit",offset,limit) if result==false then return {} end for i, member in ipairs(result) do if i%2~=0 then redis.call('ZREM',key,member) end end

--zrangebyscore 22 1 100 limit 10 100
local key=KEYS[1] local min=ARGV[1] local max=ARGV[2] local offset=ARGV[3] local limit=ARGV[4] local result=redis.call('zrangebyscore',key,min,max,"limit",offset,limit) if result==false then return {} end for i, member in ipairs(result) do redis.call('ZREM',key,member) end return result

