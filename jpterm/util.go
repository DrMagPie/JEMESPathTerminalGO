package jpterm

import "fmt"

func print(s string, x int) {
	// \033c
	fmt.Printf("\033[%dA\033[2J\n%s\n", x, s)
}
