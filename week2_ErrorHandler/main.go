package main

import (
	"JeekDemo/week2_ErrorHandler/xerr"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

func main() {

	err := queryUserInfo("lisa")
	log.Println(err)
}

func queryUserInfo(username string) error {

	db, err := initDB()
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("db init error"), " db init error : %v", err)
	}
	var address, email string
	sqlFmt := fmt.Sprintf("select address, email from account_info where username = '%s' ", username)
	err = db.QueryRowContext(context.Background(), sqlFmt).Scan(&address, &email)
	switch err {
	case sql.ErrNoRows:
		return errors.Wrapf(xerr.NewErrCode(xerr.DATA_NOT_FOUND), "not found user address")
	case nil:
		return errors.Wrapf(xerr.NewErrMsg("query user address error"), "Failed to query merkle proofs err : %v", err)
	default:
		return err
	}
}

func initDB() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/fallen_angel_swap_mainnet?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
}
