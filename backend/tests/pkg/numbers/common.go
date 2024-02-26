package numbers

import "math/rand"

func RandomInt(left, right int) int {
	return rand.Intn(right-left+1) + left
}
