// main.go

package main

func main() {

	a := App{}
	a.Initialize(
		"postgres",
		"postgres",
		"postgres",
		"postgres.default.svc.cluster.local",
		"postgres")

	a.Run(":8010")
}
