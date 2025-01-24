package godm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	CR   = "\r"
	LF   = "\n"
	CRLF = CR + LF
	LFCR = LF + CR
)

// TODO: write test
func detectLineEnding(raw string) string {
	if strings.Contains(raw, CRLF) {
		return CRLF
	}
	if strings.Contains(raw, LFCR) {
		return LFCR
	}
	if strings.Contains(raw, CR) {
		return CR
	}
	return LF
}

type KV struct {
	Key   string
	Value string
}

func parseIntoKVs(raw string) ([]KV, error) {
	var (
		kvs        = []KV{}
		lineEnding = detectLineEnding(raw)
	)

	lines := strings.Split(raw, lineEnding)
	for _, line := range lines {
		// skip empty lines
		if line == "" {
			continue
		}
		k, v, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		kvs = append(kvs, KV{k, v})
	}

	return kvs, nil
}

func parseLine(line string) (k string, v string, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", fmt.Errorf("can't parse empty line")
	}

	// parse comment
	if strings.HasPrefix(line, "COMMENT") {
		return "COMMENT", strings.TrimSpace(line[7:]), nil
	}

	// parse KVN
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid input: %s", line)
	}

	value := strings.TrimSpace(parts[1])

	// remove units if present
	unitIndex := strings.LastIndex(value, "[")
	if unitIndex != -1 && strings.HasSuffix(value, "]") {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(value[:unitIndex]), nil
	}

	return strings.TrimSpace(parts[0]), value, nil
}

// TODO: write tests
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// TODO: write tests
func parseTime(s string) (time.Time, error) {
	if !strings.Contains(s, "Z") {
		s += "Z"
	}
	return time.Parse(time.RFC3339, s)
}
