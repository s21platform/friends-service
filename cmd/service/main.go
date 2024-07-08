package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/config"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	fmt.Printf("%+v\n", cfg)
}
