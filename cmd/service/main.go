package main

import (
	"github.com/s21platform/friends-service/internal/config"
)

func main() {
	//чтение конфига
	_ := config.MustLoad()
}
