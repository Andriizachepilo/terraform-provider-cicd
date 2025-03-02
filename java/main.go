package main

import (
	"fmt"
	"os"
)


func main() {
if os.Getenv("OLA") != "" {
	env := os.Getenv("OLA")
	fmt.Print(env)
}
}


