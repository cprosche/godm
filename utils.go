package godm

import (
	"fmt"
	"strings"
)

const (
	CR   = "\r"
	LF   = "\n"
	CRLF = CR + LF
	LFCR = LF + CR
)

func parseLine(line string) (k string, v string, err error) {
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "COMMENT") {
		return "COMMENT", strings.TrimSpace(line[7:]), nil
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid input: %s", line)
	}

	value := strings.TrimSpace(parts[1])
	unitIndex := strings.LastIndex(value, "[")

	if unitIndex != -1 && strings.HasSuffix(value, "]") {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(value[:unitIndex]), nil
	}

	return strings.TrimSpace(parts[0]), value, nil
}

func detectLineEnding(s string) string {
	if strings.Contains(s, CRLF) {
		return CRLF
	} else if strings.Contains(s, LFCR) {
		return LFCR
	} else if strings.Contains(s, CR) {
		return CR
	} else {
		return LF
	}
}
