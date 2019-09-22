package main

import (
	"github.com/iskendria-pub/iskendria/dao"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "generate-dao", log.Flags())
	dao.Init("ref.db", logger)
}
