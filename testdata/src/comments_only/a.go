package comments_only

import "fmt"

func A() {
	// This is is a comment with duplicate words // want `Duplicate words \(is\) found`
	fmt.Println("hello")

	// This comment has no duplicates, but the string below does
	line := "this string has has duplicate words"
	fmt.Println(line)
}
