package errors

import (
	"runtime"
	"strconv"
	"strings"
)

func callers(skip int) (stack string) {
	var sb strings.Builder
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		} else {
			sb.WriteByte('\n')
		}
		sb.WriteString(file)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(line))
	}

	stack = sb.String()
	return
}
