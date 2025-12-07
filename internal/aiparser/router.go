package aiparser

import (
	"github.com/geminal/skube/internal/parser"
)

func ParseWithRouter(args []string) *parser.Context {
	if !shouldUseAI(args) {
		return parser.ParseNaturalLanguage(args)
	}

	aiParser, err := NewAIParser()
	if err != nil {
		return parser.ParseNaturalLanguage(args)
	}

	ctx, err := aiParser.Parse(args)
	if err != nil {
		return parser.ParseNaturalLanguage(args)
	}

	return ctx
}

func shouldUseAI(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if args[0] == "--ai" {
		return true
	}

	return false
}

func StripAIFlag(args []string) []string {
	if len(args) > 0 && args[0] == "--ai" {
		return args[1:]
	}
	return args
}

func HasAIFlag(args []string) bool {
	if len(args) > 0 && args[0] == "--ai" {
		return true
	}
	return false
}
