package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/saurav-gahlawat/pseudo-proc/internal/procfs"
)

func main() {
	f, err := os.Open("/proc/984/stat")
	if err != nil {
		fmt.Println("Error Opening file: ", err)
		return
	}
	defer f.Close()

	info, err := procfs.ParsePidStat(f)

	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(info)
}
