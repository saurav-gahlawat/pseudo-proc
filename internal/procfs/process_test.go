package procfs

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestParsePidStat(t *testing.T) {
	f, err := os.Open("testdata/pidstat")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		input   io.Reader
		want    Process
		wantErr bool
	}{
		{
			name:  "Happy Path",
			input: f,
			want:  Process{PID: 2676257, Name: "bash", State: 83, PPID: 2676062, UTime: 10, STime: 5, Threads: 1, RSS: 8269824},
		},
		{
			name:    "invalid input",
			input:   strings.NewReader("Invalid Input\n"),
			wantErr: true,
		},
		{
			name:  "comm with spaces and parens",
			input: strings.NewReader("999 (my (evil) proc) S 1 999 999 0 -1 0 0 0 0 0 10 5 0 0 20 0 1 0 0 4096 1 18446744073709551615 0 0 0 0 0 0 0 0 0 0 0 0 17 0 0 0 0 0 0"),
			want:  Process{PID: 999, Name: "my (evil) proc", State: 'S', PPID: 1, UTime: 10, STime: 5, Threads: 1, RSS: int64(os.Getpagesize())},
		},
		{
			name:    "too few fields after comm",
			input:   strings.NewReader("123 (foo) S 456\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePidStat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Got error while parsing the reader: %v for test: %q, [wantErr : %v]", err, tt.name, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParsePidStat() = %v, want = %v", got, tt.want)
			}
		})
	}
}
