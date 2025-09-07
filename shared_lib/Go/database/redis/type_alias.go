package redis

import "github.com/redis/go-redis/v9"

const (
	RedisNil = redis.Nil

	SearchFieldTypeInvalid  = redis.SearchFieldTypeInvalid
	SearchFieldTypeText     = redis.SearchFieldTypeText
	SearchFieldTypeNumeric  = redis.SearchFieldTypeNumeric
	SearchFieldTypeTag      = redis.SearchFieldTypeTag
	SearchFieldTypeGeo      = redis.SearchFieldTypeGeo
	SearchFieldTypeVector   = redis.SearchFieldTypeVector
	SearchFieldTypeGeoShape = redis.SearchFieldTypeGeoShape
)

type (
	FailoverOptions    = redis.FailoverOptions
	Client             = redis.Client
	IntCmd             = redis.IntCmd
	StringCmd          = redis.StringCmd
	MapStringStringCmd = redis.MapStringStringCmd
	BoolCmd            = redis.BoolCmd
	Pipeliner          = redis.Pipeliner
	RedisCmd           = redis.Cmd
	FTCreateOptions    = redis.FTCreateOptions
	FieldSchema        = redis.FieldSchema
	FTSearchCmd        = redis.FTSearchCmd
	FTSearchOptions    = redis.FTSearchOptions
	SearchFieldType    = redis.SearchFieldType
)
