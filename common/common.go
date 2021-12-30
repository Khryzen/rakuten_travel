package common

import (
	"database/sql"
	"fmt"
	"strings"
)

func Print(level int, msg interface{}, i ...interface{}) {
	message := fmt.Sprint(msg)
	if level >= OK {
		if !strings.HasSuffix(message, "\n") {
			message += "\n"
		}
		fmt.Printf(levelMap[level]+message, i...)
	}
}

func CheckErr(msg string, err error) {
	if err != nil {
		Print(ERROR, "%s | %s", msg, err)
	}
	return
}

func ErrMsg(msg string, err error) (*sql.DB, error) {
	if err != nil {
		Print(ERROR, "%s | %s", msg, err)
	}
	return nil, err
}
