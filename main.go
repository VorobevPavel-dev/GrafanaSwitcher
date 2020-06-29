package main

func main() {
	user := new(UserSession)
	user.Init()
	user.PostDashboardModel("output.json")

}
