package main

import (
	"math/rand"
	"testing"
)

func TestGetAllFile(t *testing.T) {
	file, err := GetAllFile(".")
	if err != nil {
		t.Error(err)
	}
	println("TestGetAllFile: ", file[0])
}

func TestSecondToTime(t *testing.T) {
	time := SecondToTime(rand.Int63())

	println("SecondToTime: ", time.String())
}
