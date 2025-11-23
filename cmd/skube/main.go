package main

import (
	"fmt"
	"os"

	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/executor"
	"github.com/geminal/skube/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		executor.PrintHelp()
		os.Exit(0)
	}

	args := os.Args[1:]
	ctx := parser.ParseNaturalLanguage(args)

	if err := executor.ExecuteCommand(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
		os.Exit(1)
	}
}
