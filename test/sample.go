package main

// +langconv
const (
	CONST1 int32   = 1
	CONST2 string  = "hello hello hello"
	CONST3 float64 = 3.1412
)

// +langconv
type User struct {
	Username string
	Age      int
	Parent   User
	Items    []Item
}

type Item struct {
	Name   string
	Amount int
}
