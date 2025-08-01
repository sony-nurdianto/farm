package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func hashConfig(cfg kafka.ConfigMap) string {
	var keys []string
	for k := range cfg {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%s=%v", k, cfg[k]))
	}

	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}
