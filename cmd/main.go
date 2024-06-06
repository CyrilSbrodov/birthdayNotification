package main

import "birthdayNotification/internal/app"

func main() {
	srv := app.NewServerApp()
	srv.Run()
}
