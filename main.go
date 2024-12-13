package main

import (
	db "github.com/Mans397/eLibrary/Database"
	sc "github.com/Mans397/eLibrary/serverConnection"
)

func main() {
	db.Init()

	sc.ConnectToServer()
}
