package main

import (
	"fmt"
	config "gator/internal/config"
)

func main() {
	cfg, _ := config.Read()
	cfg.SetUser("will")
	fmt.Println(cfg)

}
