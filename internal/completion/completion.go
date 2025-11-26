package completion

import (
	_ "embed"
	"fmt"

	"github.com/geminal/skube/internal/parser"
)

//go:embed skube.zsh
var zshCompletion string

//go:embed skube.bash
var bashCompletion string

func HandleCompletion(ctx *parser.Context) error {
	shell := ctx.ResourceType
	if shell == "" {
		return fmt.Errorf("please specify shell type\nUsage: skube completion <zsh|bash>")
	}

	switch shell {
	case "zsh":
		fmt.Print(zshCompletion)
		return nil
	case "bash":
		fmt.Print(bashCompletion)
		return nil
	default:
		return fmt.Errorf("unsupported shell: %s\nSupported shells: zsh, bash", shell)
	}
}
