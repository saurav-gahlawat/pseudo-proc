package procfs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CPUTimes struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
}

func ParseStat(r io.Reader) (CPUTimes, error) {
	var result CPUTimes
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 9 {
			return result, fmt.Errorf("parsing /proc/stat: expected at least 9 fields, got %d", len(fields))
		}
		targets := []*uint64{
			&result.User, &result.Nice, &result.System, &result.Idle,
			&result.IOWait, &result.IRQ, &result.SoftIRQ, &result.Steal,
		}
		for i, target := range targets {
			val, err := strconv.ParseUint(fields[i+1], 10, 64)
			if err != nil {
				return result, fmt.Errorf("parsing /proc/stat field %d: %w", i+1, err)
			}
			*target = val
		}
	}
	return result, scanner.Err()
}
