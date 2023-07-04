package database

import (
	"time"
)

// Note: All the explicit column/table names are explicitly set to match the current DB structure

var allTables = []any{&BadEgg{}, &Block{}, &StaffLog{}, &MonthLog{}, &Watching{}, &UserReplyThread{}}

type BadEgg struct {
	DbId     int       `gorm:"primaryKey;column:dbid"`
	UserId   int       `gorm:"column:id"`
	Username string    `gorm:"column:username"`
	Number   int       `gorm:"column:num"`
	Date     time.Time `gorm:"column:date;type:DATE"`
	Message  string    `gorm:"column:message"`
	Staff    string    `gorm:"column:staff"`
	Post     int       `gorm:"column:post"`
}

func (BadEgg) TableName() string {
	return "badeggs"
}

type Block struct {
	Id string `gorm:"primaryKey;column:id"`
}

func (Block) TableName() string {
	return "blocks"
}

type StaffLog struct {
	Staff string `gorm:"primaryKey;column:staff"`
	Bans  int    `gorm:"column:bans"`
	Warns int    `gorm:"column:warns"`
}

func (StaffLog) TableName() string {
	return "staffLogs"
}

type MonthLog struct {
	Staff string `gorm:"primaryKey;column:staff"`
	Bans  int    `gorm:"column:bans"`
	Warns int    `gorm:"column:warns"`
}

func (MonthLog) TableName() string {
	return "monthLogs"
}

type Watching struct {
	UserId string `gorm:"primaryKey;column:id"`
}

func (Watching) TableName() string {
	return "monthLogs"
}

type UserReplyThread struct {
	UserId   int `gorm:"primaryKey;column:userid"`
	ThreadId int `gorm:"primaryKey;uniqueIndex;column:threadid"`
}

func (UserReplyThread) TableName() string {
	return "userReplyThreads"
}
