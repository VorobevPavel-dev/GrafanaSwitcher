package main

import "fmt"

func main() {
	user := UserSession{}
	user.Init()
	UIDlist, err := user.GetUIDList()
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range UIDlist {
		user.GetDashboardModel(i)
	}
}
