package main

import "project_yandex_lms/lms-calc-with-gorutine/pkg/application"

func main() {
	app := application.New()
	app.Setup()
	app.Run()
}
