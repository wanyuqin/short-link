package consts

const RedisPrefix = "shor_link"

const (
	RedisKeyShorUrl           = RedisPrefix + ":%s"
	RedisKeyShortUrlBlackList = RedisPrefix + ":blacklist:%s"
)
