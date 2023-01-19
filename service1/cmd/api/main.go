package main

import (
	"fmt"
	serve "taobin-service/serve"
)

func main() {
	err := serve.ServeHTTP()
	if err != nil {
		fmt.Println(err)
	}
}
