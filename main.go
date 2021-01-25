// main.go

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("APP_DB_USERNAME"))
	a := App{}
	a.Initialize(
		"admin",
		"admin123",
		"postgresdb",
		"postgres.ingress-basic.svc.cluster.local")

	a.Run(":8010")
}