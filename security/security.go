package security

import (
	"flag"
	"fmt"
)

var (
	TYPE = flag.String("type", "props", "file type")
	IN   = flag.String("in", "x.in", "input file")
	OUT  = flag.String("out", "x.out", "output file")
)

func init() {
	flag.Parse()
}

func main() {
	//var err error
	fmt.Println(*TYPE)
	switch *TYPE {
	case "props":
		props()
	case "db":
		db()
	default:
		fmt.Println("Not found the type.")
	}
}

func db() {

}

func props() {

}
