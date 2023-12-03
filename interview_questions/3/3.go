package main

import "fmt"

// Flood fill detect
// each pixel is 2 bits (only 4 colors)
// flood fill == you fill everywhere where the color is the same as the place you're filling
// The question is : given a pattern (2bit) and a byte, figure out if the pattern is included in the byte

func main() {

	assert(containsColorNoLoop(0b10000011, 0b11), true)
	assert(containsColorNoLoop(0b10000001, 0b11), false)
	assert(containsColorNoLoop(0b00000011, 0b10), false)
	assert(containsColorNoLoop(0b00010011, 0b10), false)
	assert(containsColorNoLoop(0b00100011, 0b10), true)
	assert(containsColorNoLoop(0b11000000, 0b11), true)
	assert(containsColorNoLoop(0b01100000, 0b11), false)
	assert(containsColorNoLoop(0b01100000, 0b00), true)
	assert(containsColorNoLoop(0b00111111, 0b00), true)
	assert(containsColorNoLoop(0b11100111, 0b00), false)

	fmt.Println("all tests passed")
}

// This solution is not really optimized..
func containsColor(pixel, color uint8) bool {
	mask := uint8(0b11)
	for i := 0; i < 4; i++ {
		if pixel&mask == color {
			return true
		}
		pixel >>= 2
	}
	return false
}

// Pixels contains 4 pixel that are 2 bits, color contains only 2 bits
func containsColorNoLoop(pixels, color uint8) bool {
	pattern := color | color<<2 | color<<4 | color<<6

	identicalBits := (pixels ^ pattern) ^ (0b11111111)
	highBitMatch := identicalBits & 0b10101010
	bothBitMatch := (highBitMatch >> 1) & identicalBits
	return bothBitMatch > 0
}

func assert(have, want bool) {
	if have != want {
		panic("")
	}
}
