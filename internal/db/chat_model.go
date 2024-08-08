package db

import "fmt"

type TgChat struct {
	Id        int64
	Type      string
	UserName  string
	FirstName string
	LastName  string
}

func (c TgChat) String() string {
	return fmt.Sprintf("{ id:%v, type:%v, username:%v, firstname:%v, lastname:%v }", c.Id, c.Type, c.UserName, c.FirstName, c.LastName)
}
