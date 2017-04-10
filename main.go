package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/datake914/actibeat/beater"
)

func main() {
	err := beat.Run("actibeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
