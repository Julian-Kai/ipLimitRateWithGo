package main

import (
	"ip.limit.rate/internal/api/rest"
	"ip.limit.rate/internal/pkg/log"

	"github.com/sirupsen/logrus"
)

func main() {
	log.Initial()

	if err := rest.NewRouter().Run(); err != nil {
		logrus.Fatal("Failed to run router")
	}
}
