package main

import (
	"github.com/jmoiron/sqlx"
)

func main() {

}

func InitDB() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/fallen_angel_swap_mainnet?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
}
