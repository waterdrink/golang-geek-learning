package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var db *sql.DB

func main() {
	fmt.Println(getUsernameWithId(123))

}

func getUsernameWithId(id int) (username string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 老师，我对于是否有必要调用cancel()一直有疑问，一般代码中都会多写这么一句，
	// 但也没有看到哪里有文档对什么时候需要cancel有详细的解释
	defer cancel()

	err = db.QueryRowContext(ctx, "SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	// ErrNoRows 不应该直接抛出，其带有业务含义，应在内部处理
	case errors.Is(err, sql.ErrNoRows):
		return "", errors.Errorf("no user with id %d", id)
	case err != nil:
		return "", errors.Wrapf(err, "query error: %v", err)
	default:
		return username, nil
	}
}
