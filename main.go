package main

import "fmt"

func main() {
	SystemStartup()
	Setup()
}

func SystemStartup() {
	fmt.Printf("Running on http://%v%s/...\n", BindIP, Port)
}
