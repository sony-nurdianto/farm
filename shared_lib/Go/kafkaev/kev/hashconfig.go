package kev

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// hashConfig membuat hash dari configuration untuk digunakan sebagai key
func hashConfig(cfg kafka.ConfigMap) string {
	var keys []string
	for key := range cfg {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var configStr strings.Builder
	for _, key := range keys {
		val, _ := cfg.Get(key, "")
		configStr.WriteString(fmt.Sprintf("%s=%v;", key, val))
	}

	hash := sha256.Sum256([]byte(configStr.String()))
	return fmt.Sprintf("%x", hash)
}
