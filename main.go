package main

func main() {
	user := UserSession{}
	user.Init()
	user.GetDashboardModel("000000081")
}
