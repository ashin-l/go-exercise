package main

import (
	"fmt"
	"regexp"
)

const text = `
My email is hehe@gmail.com
email1 is sss@lll.org
email2 is 		kkk@qq.com
email3 is 		lll@qq.com.cn
`

func main() {
	re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)(\.[a-zA-Z0-9.]+)`)
	match := re.FindAllStringSubmatch(text, -1)
	for _, m := range match {
		fmt.Println(m)
	}
}
