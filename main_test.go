package main

import (
	"os"
	"testing"
	_ "github.com/the-go-dragons/fake-airline-info-service/config"
	_ "github.com/the-go-dragons/fake-airline-info-service/server"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	return m.Run()
}
