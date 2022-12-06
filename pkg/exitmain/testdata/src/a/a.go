package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println("Hello world")

	os.Exit(1) // want "direct call os.Exit from main function"
}
