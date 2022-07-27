package main

import (
	"fmt"
	"github.com/galenliu/chip/cmd/commission"
	"os"
)

func main() {
	if err := commission.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
