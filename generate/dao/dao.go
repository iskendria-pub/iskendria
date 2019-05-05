package main

import (
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "generate-dao", log.Flags())
	dao.Init("ref.db", logger)
}
