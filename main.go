package main

import (
	"fmt"
	"os"
)

var (
	routes []string
)

func init() {
	//Проверка необходимых директорий и их создание
	if _, err := os.Stat("./Backups"); os.IsNotExist(err) {
		os.Mkdir("./Backups", 0777)
	}
	if _, err := os.Stat("./Changed"); os.IsNotExist(err) {
		os.Mkdir("./Changed", 0777)
	}
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		_, e := os.Create("config.json")
		if e != nil {
			fmt.Println("Cannot create config.json file. Try to create it manually")
			os.Exit(1)
		}
		fmt.Println("File config.json has been created. Check README.md to find out necessary information")
		os.Exit(1)
	}
}

func main() {
	user := new(UserSession)
	user.Init()
	dashboard := new(Dashboard)
	dashboard.Init("lwg_i7MMk")
	dashboard, _ = dashboard.Restore(user)
}
