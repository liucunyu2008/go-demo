package main

import (
	"os"
	"testing"
)

//func main() {
//var m *testing.M
//	TestMain(m)
//}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}