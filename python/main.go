package main

import (
	"fmt"
	"regexp"
)

func main() {
	val := "aws_account_id.dkr.ecr.england.amazonaws.com"

   kk := regexp.MustCompile(`(?:[^\.]*\.){3}([^\.]*)`).FindStringSubmatch(val)[1]
   fmt.Print(kk)
}

// func iter(link string) (string) {
// 	var dots int
// 	var start int
// 	var finish int
// 	for i, v := range link {
// 		if string(v) == "." {
// 			dots++
// 		} 
// 		if dots == 3 {
// 			start = i
// 		} else if dots == 4 {
// 			finish = i
// 		}
		
// 	}
// 	return link[start:finish]
// }
