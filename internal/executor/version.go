package executor

import (
	"fmt"

	"github.com/geminal/skube/internal/config"
)

const Version = "0.2.2"

func PrintVersion() {
	fmt.Printf("%sskube%s version %s%s%s\n",
		config.ColorGreen, config.ColorReset,
		config.ColorCyan, Version, config.ColorReset)
}
