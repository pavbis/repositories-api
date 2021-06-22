package main

import "github.com/pavbis/zal-case-study/api"

func main() {
	s := api.Server{}
	s.Initialize()
	s.Run(":7000")
}
