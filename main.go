package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"scootin/cmd"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go cmd.InitializeRestfulApi(&wg)
	go cmd.InitializeGrpcApi(&wg)
	//go cmd.ExecuteClientBehaviour()
	wg.Wait()
}
