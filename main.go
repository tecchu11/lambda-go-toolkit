package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tecchu11/lambda-invoker-go/internal"
)

func main() {
	var (
		p = flag.Int("p", 9000, "Port number. Default 9000.")
		f = flag.String("f", "", "Event file location. Default event.json.")
	)
	flag.Parse()

	client, err := internal.New(*p)
	if err != nil {
		log.Fatal(err)
	}
	b, err := os.ReadFile(*f)
	if err != nil {
		log.Fatalf("read event file\n%v", err)
	}
	res, err := client.Do(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}
