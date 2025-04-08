package main

import "fmt"

func main() {
	dst := []int{0, 1, 2, 3, 4, 5, 8, 9, 10, 11, 12, 13, 14, 15}

	// Remove value at index 3
	copy(dst[3:], dst[4:]) // shift elements left
	fmt.Printf("After copy: %v, len: %d\n", dst, len(dst))

	// Truncate the last element (now duplicated)
	dst = dst[:len(dst)-1]
	fmt.Printf("Truncated dst:%v, len:%d\n", dst, len(dst))
}
