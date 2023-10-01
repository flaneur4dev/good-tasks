package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printHelp() {
	txt := `Calendar application for events notifications.
	Usage: calendar [--config=/path/to/config/config.yaml] [help] [version]`
	fmt.Println(txt)
}

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
