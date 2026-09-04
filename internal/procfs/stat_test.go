package procfs

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseStat(t *testing.T) {
	f, err := os.Open("testdata/stat")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tests := []struct {
		name    string
		input   io.Reader
		want    CPUTimes
		wantErr bool
	}{
		{
			name:  "valid proc/stat",
			input: f,
			want:  CPUTimes{2423625, 884, 765110, 92449180, 103508, 269090, 98508, 0},
		},
		{
			name:    "malformed input",
			input:   strings.NewReader("cpu garbage values\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseStat() = %v, want %v", got, tt.want)
			}
		})
	}
}
