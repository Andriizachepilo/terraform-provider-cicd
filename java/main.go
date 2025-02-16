package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "gcr.io/andruha"
    check := regexp.MustCompile(`^[^/]+/(.+)$`).FindStringSubmatch(str)
    fmt.Print(check[1])
}


// func checking (anuthing string) string {
// 	for _, v := range anuthing {
// 		if v
// 	}
// }