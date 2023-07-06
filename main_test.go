package main

import (
	"os"
	"testing"
	"github.com/the-go-dragons/fake-airline-info-service/config"
	"github.com/the-go-dragons/fake-airline-info-service/server"
)

func TestMain(m *testing.M) {
	os.Exit(testMain())
}

func testMain(m *testing.M) int {
	return m.Run()
}
