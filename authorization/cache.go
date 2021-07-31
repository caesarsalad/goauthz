package authorization

import (
	"log"
	"regexp"
	"time"

	"github.com/caesarsalad/goauthz/database"
)

var Rule_regex_compiled = make(map[uint]*regexp.Regexp)
var Cached_user_rules = CachedUserRules{CacheManager: CacheManager{LastModifiedTimeKey: "user_rules_ts_",
	LastModifiedTime: time.Now().Unix()}, Cache: make(map[uint][]userRules)}

type CachedUserRules struct {
	CacheManager
	Cache map[uint][]userRules
}

/* TODO: Cache validation using Redis.
func (c *CacheManager) IsValid(suffix_key ...string) (bool, error) {

	local_cache_ts := c.LastModifiedTime
	suffix_cache_key := ""
	if len(suffix_key) > 0 {
		suffix_cache_key = suffix_key[0]
	}
	cache_key := c.LastModifiedTimeKey + suffix_cache_key
	return true, nil
}

func (c *CacheManager) UpdateTime(suffix_key ...string) error {
	now_ts := time.Now().Unix()
	c.LastModifiedTime = now_ts
	suffix_cache_key := ""
	if len(suffix_key) > 0 {
		suffix_cache_key = suffix_key[0]
	}
	cache_key := c.LastModifiedTimeKey + suffix_cache_key
	return nil
}
*/

func CompileAllRegexRules() {
	var user_rules []database.Rule
	database.DB.Where("meta_location_id = ?", MetaLocationUrl).Find(&user_rules)
	for _, rule := range user_rules {
		r, err := regexp.Compile(rule.Path)
		if err != nil {
			log.Println("error while compile regex ", rule.Path, err)
			continue
		}
		Rule_regex_compiled[rule.ID] = r
	}
}

func getUserRules(user_id uint) []userRules {
	if user_rules, ok := Cached_user_rules.Cache[user_id]; ok {
		return user_rules
	}
	log.Println("dont hit")
	var user_rules []userRules
	database.DB.Table("assigned_rules").
		Select("rules.id", "rules.path", "assigned_rules.meta_value", "rules.meta_key",
			"rules.meta_location_id", "rules.http_method_id", "rules.path_prefix").
		Joins("INNER JOIN rules ON rules.id = assigned_rules.rule_id").
		Where("assigned_rules.user_id = ?", user_id).Scan(&user_rules)
	Cached_user_rules.Cache[user_id] = user_rules
	return user_rules
}

func ReCacheUserRules() {
	var user_ids []uint
	for user_id := range Cached_user_rules.Cache {
		user_ids = append(user_ids, user_id)
	}
	for _, user_id := range user_ids {
		delete(Cached_user_rules.Cache, user_id)
		getUserRules(user_id)
	}
}
