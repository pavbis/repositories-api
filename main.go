package main

import "github.com/pavbis/repositories-api/api"

func main() {
	s := api.Server{}
	s.Initialize()
	s.Run(":7000")
}
