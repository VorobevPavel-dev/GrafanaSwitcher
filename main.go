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
	_, err := user.GetDashboardModel("lwg_i7MMk")
	if err != nil {
		fmt.Println(err)
	}
	_, err = user.ChangeTag("lwg_i7MMk", `sbis3mon.prod.linux.memory.ins-db2.interval-28sec.mem.ram.used.avg`, "semi-dark-orange")
	if err != nil {
		fmt.Println(err)
	}
}
