package main

import (
	"fmt"

	"github.com/username/module"
)

func main() {
	// Example usage of your module
	client := module.New()
	result := client.DoSomething()
	fmt.Println(result)
}
