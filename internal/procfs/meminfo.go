package procfs

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type MemInfo struct {
	MemTotal     uint64
	MemAvailable uint64
}

func ParseMemInfo(r io.Reader) (MemInfo, error) {
	var result MemInfo
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		val, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}
		switch key {
		case "MemTotal":
			result.MemTotal = val * 1024
		case "MemAvailable":
			result.MemAvailable = val * 1024
		}
	}
	return result, scanner.Err()
}
