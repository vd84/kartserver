// main.go

package main

import "os"
import "fmt"

func main() {
	fmt.Println(os.Getenv("APP_DB_USERNAME"))
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8010")
}