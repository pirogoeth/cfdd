package util

import (
	"fmt"
)

func BuildFQDN(recordName, zoneName string) string {
	return fmt.Sprintf("%s.%s", recordName, zoneName)
}
