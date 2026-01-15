package main

import (
	"github.com/joho/godotenv"
	"github.com/markdave123-py/crypto-portfolio-tracker/cmd"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
