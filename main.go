package main

import (
	"github.com/rmarasigan/rakuten_travel/common"
)

func main() {
	SystemStartup()
	Setup()
}

func SystemStartup() {
	common.Print(common.OK, "Running on http://%v%s/...\n", BindIP, Port)
}
