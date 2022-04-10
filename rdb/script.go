package rdb

const (
	popscore          = `local key=KEYS[1] local start=ARGV[1] local stop=ARGV[2] local result=redis.call('ZRANGE',key,'0','10','WITHSCORES') if result==false then return {} end for i, member in ipairs(result) do if i%2~=0  then redis.call('ZREM',key,member) end end return result`
	pop               = `local key=KEYS[1] local start=ARGV[1] local stop=ARGV[2] local result=redis.call('ZRANGE',key,start,stop) if result==false then return {} end for i, member in ipairs(result) do redis.call('ZREM',key,member) end return result`
	zrangeByScoreWith = `local key=KEYS[1] local min=ARGV[1] local max=ARGV[2] local offset=ARGV[3] local limit=ARGV[4] local result=redis.call('zrangebyscore',key,min,max,'withscores',"limit",offset,limit) if result==false then return {} end for i, member in ipairs(result) do if i%2~=0 then redis.call('ZREM',key,member) end end return result`
	zrangeByScore     = `local key=KEYS[1] local min=ARGV[1] local max=ARGV[2] local offset=ARGV[3] local limit=ARGV[4] local result=redis.call('zrangebyscore',key,min,max,"limit",offset,limit) if result==false then return {} end for i, member in ipairs(result) do redis.call('ZREM',key,member) end return result`
)

var (
	zpopScoreSha    string
	zpopSha         string
	zByScoreWithSha string
	zByScoreSha     string
)
