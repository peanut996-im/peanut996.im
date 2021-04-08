package main

import "fmt"

var (
	BuildVersion, BuildTime, BuildUser, BuildMachine string
)

func main() {
	fmt.Println("Build Info:")
	fmt.Println("Build Version: ", BuildVersion)
	fmt.Println("Build Time: ", BuildTime)
	fmt.Println("Build Machine: ", BuildMachine)
	fmt.Println("Build User: ", BuildUser)
	fmt.Println("This is a starter project for peanut996.im")
	fmt.Println("Please copy and then edit for project.")
}
