package main

import (
	"log"

	"github.com/wanna-beat-by-bit/goMetricService/internal/app"
)

func main() {
	//TODO: сделать конфиг
	//TODO: сделать логгер
	//TODO: сделать работу с хранилищами
	//TODO: сделтаь эндпойнт
	//TODO: запустить сервер

	a, err := app.New()
	if err != nil {
		log.Fatalf("Error occured while starting application: %s", err.Error())
	}

	a.Run()

}
