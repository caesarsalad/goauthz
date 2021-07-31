package authorization

var (
	MetaLocationQuery uint = 1
	MetaLocationBody  uint = 2
	MetaLocationUrl   uint = 3
)

var MetaLocationIDMap = map[string]uint{
	"query": 1,
	"body":  2,
	"url":   3,
}

var (
	HTTPMethodGET     uint = 1
	HTTPMethodHEAD    uint = 2
	HTTPMethodPOST    uint = 3
	HTTPMethodPUT     uint = 4
	HTTPMethodDELETE  uint = 5
	HTTPMethodCONNECT uint = 6
	HTTPMethodOPTIONS uint = 7
	HTTPMethodTRACE   uint = 8
	HTTPMethodPATCH   uint = 9
)

var HttpMethodIDMap = map[string]uint{
	"GET":     1,
	"HEAD":    2,
	"POST":    3,
	"PUT":     4,
	"DELETE":  5,
	"CONNECT": 6,
	"OPTIONS": 7,
	"TRACE":   8,
	"PATCH":   9,
}

type CacheManager struct {
	LastModifiedTimeKey string
	LastModifiedTime    int64
}
