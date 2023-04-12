package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Starting tests!")
	os.Setenv("DB_SCHEMA", "test")
	Db = *SetUpDb()
	defer Db.Close()
	exitVal := m.Run()
	log.Println("Ending tests!")

	os.Exit(exitVal)
}
