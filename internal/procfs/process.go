package procfs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	PID     int64
	Name    string
	State   byte
	PPID    int64
	UTime   int64
	STime   int64
	Threads int64
	RSS     int64
}

func ParsePidStat(r io.Reader) (Process, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		output := scanner.Text()
		startName := strings.IndexByte(output, '(')
		endName := strings.LastIndexByte(output, ')')
		if startName == -1 || endName == -1 {
			return Process{}, fmt.Errorf("Invalid proc/<pid>/stat file")
		}
		pidStr := strings.TrimSpace(output[:startName])
		parsedPid, err := strconv.ParseUint(pidStr, 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing pid: %w", err)
		}

		name := output[startName+1 : endName]

		rest := output[endName+2:]
		fields := strings.Fields(rest)

		if len(fields) < 22 {
			return Process{}, fmt.Errorf("parsing /proc/<pid>/stat: expected at least 22 fields after comm, got %d", len(fields))
		}

		ppid, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing ppid: %w", err)
		}
		utime, err := strconv.ParseUint(fields[11], 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing utime: %w", err)
		}
		stime, err := strconv.ParseUint(fields[12], 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing stime: %w", err)
		}
		threads, err := strconv.ParseUint(fields[17], 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing threads: %w", err)
		}
		rss, err := strconv.ParseInt(fields[21], 10, 64)
		if err != nil {
			return Process{}, fmt.Errorf("parsing rss: %w", err)
		}

		return Process{
			PID:     int64(parsedPid),
			Name:    name,
			State:   fields[0][0],
			PPID:    int64(ppid),
			UTime:   int64(utime),
			STime:   int64(stime),
			Threads: int64(threads),
			RSS:     rss * int64(os.Getpagesize()),
		}, nil
	}
	return Process{}, scanner.Err()
}
