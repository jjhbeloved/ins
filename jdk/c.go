package jdk

import (
	"asiainfo.com/ins/utils"
	"fmt"
)

func main()  {
	err := utils.Zip("/Users/david/test","/tmp/c.zip")
	fmt.Println(err)
	err = utils.Unzip("/tmp/c.zip","/tmp/t2")
	fmt.Println(err)
}
