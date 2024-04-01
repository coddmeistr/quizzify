package numbers

import "math/rand"

// RandomInt returns random integer in range [left, right]
func RandomInt(left, right int) int {
	return rand.Intn(right-left+1) + left
}
