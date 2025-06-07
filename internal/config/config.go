package config

import "fmt"

func Load() {
	path := getConfigPath()

	fmt.Printf("using config path: %s\n", path)
}
