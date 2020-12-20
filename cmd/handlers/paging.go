package handlers

import (
	"fmt"
)

func prevToken(offset *int64, limit *int64) string {
	if offset != nil && *offset > 0 {
		if limit != nil {
			start := *offset - *limit
			return fmt.Sprintf("%s=%d&%s=%d", "offset", start, "limit", *limit)
		}
	}
	return ""
}

func nextToken(offset *int64, limit *int64, dataLen int) string {
	if limit != nil && *limit > 0 && *limit <= int64(dataLen) {
		var start int64
		if offset != nil && *offset > 0 {
			start = *offset + *limit
		} else {
			start = *limit
		}
		return fmt.Sprintf("%s=%d&%s=%d", "offset", start, "limit", *limit)
	}
	return ""
}
