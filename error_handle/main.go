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
	name, err := getUsernameWithId(123)
	if nil != err {
		if isNoUser(err) {
			// do something
			fmt.Printf("%v\n", err)
			return
		}
		// do other thing
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Println("success: ", name)
}

type ErrNoUser struct {
	userId int
}

func (e *ErrNoUser) Error() string {
	return fmt.Sprintf("no user with id %d", e.userId)
}

func (e *ErrNoUser) isNoUser() bool {
	return true
}

type noUser interface {
	isNoUser() bool
}

func isNoUser(err error) bool {
	if nu, ok := err.(noUser); ok {
		return nu.isNoUser()
	} else {
		return false
	}
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
		return "", &ErrNoUser{userId: id}
	case err != nil:
		return "", errors.Wrapf(err, "query error: %v", err)
	default:
		return username, nil
	}
}
