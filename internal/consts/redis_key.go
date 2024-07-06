package consts

const RedisPrefix = "shor_link"

const (
	RedisKeyShorURL           = RedisPrefix + ":%s"
	RedisKeyShortURLBlackList = RedisPrefix + ":blacklist:%s"
)
