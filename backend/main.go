package main

import (
	"red-packet/router"
)

func main() {
	r := router.NewRouter()
	r.Run(":8080")
}
