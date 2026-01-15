package main

import (
	"github.com/joho/godotenv"
	"github.com/markdave123-py/crypto-portfolio-tracker/cmd"
	_ "github.com/markdave123-py/crypto-portfolio-tracker/docs"
)

// @title Crypto Portfolio Tracker API
// @version 1.0
// @description API for tracking crypto portfolios, transactions, and prices

// @contact.name David Nwaekwu
// @contact.email nwaekwudavid@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
