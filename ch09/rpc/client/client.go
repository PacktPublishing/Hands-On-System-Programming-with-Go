package main

import (
	"log"
	"net/rpc"
	"os"

	"github.com/PacktPublishing/Hands-On-Systems-Programming-with-Go/ch09/rpc/common"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address.")
	}
	client, err := rpc.DialHTTP("tcp", os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	const hp = "H.P. Lovecraft"
	var books = []common.Book{
		{ISBN: "1540335534", Author: hp, Title: "The Call of Cthulhu", Pages: 36},
		{ISBN: "1980722803", Author: hp, Title: "The Dunwich Horror ", Pages: 53},
		{ISBN: "197620299X", Author: hp, Title: "The Shadow Over Innsmouth", Pages: 40},
		{ISBN: "1540335534", Author: hp, Title: "The Case of Charles Dexter Ward", Pages: 176},
	}

	callClient(client, "ReadingService.GetProgress", books[0].ISBN, new(int))
	callClient(client, "ReadingService.AddBook", books[0], new(bool))
	callClient(client, "ReadingService.AddBook", books[0], new(bool))
	callClient(client, "ReadingService.GetProgress", books[0].ISBN, new(int))
	callClient(client, "ReadingService.AddBook", books[1], new(bool))
	callClient(client, "ReadingService.AddBook", books[2], new(bool))
	callClient(client, "ReadingService.AddBook", books[3], new(bool))
	callClient(client, "ReadingService.SetProgress", common.Progress{
		ISBN:  books[3].ISBN,
		Pages: 10,
	}, new(bool))
	callClient(client, "ReadingService.GetProgress", books[3].ISBN, new(int))
	callClient(client, "ReadingService.AdvanceProgress", common.Progress{
		ISBN:  books[3].ISBN,
		Pages: 40,
	}, new(bool))
	callClient(client, "ReadingService.GetProgress", books[3].ISBN, new(int))
}

func callClient(client *rpc.Client, method string, in, out interface{}) {
	var r interface{}
	if err := client.Call(method, in, out); err != nil {
		out = err
	}
	switch v := out.(type) {
	case error:
		r = v
	case *int:
		r = *v
	case *bool:
		r = *v
	}
	log.Printf("%s: [%+v] -> %+v", method, in, r)
}
