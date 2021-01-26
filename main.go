package main

import (
	"awesomeProject/internal/api/rest"
	"awesomeProject/internal/pkg/log"

	"github.com/sirupsen/logrus"
)

func main() {
	log.Initial()

	if err := rest.NewRouter().Run(); err != nil {
		logrus.Fatal("Failed to run router")
	}
}
