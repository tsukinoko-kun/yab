package main

import (
	"fmt"
	"os"
	"strings"

	hash "github.com/segmentio/fasthash/fnv1a"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf(
			"// %s\n%s = 0x%x\n",
            arg,
			strings.ReplaceAll(arg, "-", "_"),
			hash.HashString32(arg),
		)
	}
}
