package main

import "github.com/faridtriwicaksono/forgebe/internal/cli"

var version = "dev"

func main() {
	cli.SetVersion(version)
	cli.Execute()
}
