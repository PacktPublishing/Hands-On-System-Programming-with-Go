package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()

	var cmds = []*exec.Cmd{
		exec.Command("cat", "book_list.txt"),
		exec.Command("grep", "Game"),
		exec.Command("wc", "-l"),
	}

	cmds[1].Stdin, cmds[0].Stdout = r1, w1
	cmds[2].Stdin, cmds[1].Stdout = r2, w2
	cmds[2].Stdout = os.Stdout

	for i := range cmds {
		if err := cmds[i].Start(); err != nil {
			log.Fatalln("Start", i, err)
		}
	}

	for i, closer := range []io.Closer{w1, w2, nil} {
		if err := cmds[i].Wait(); err != nil {
			log.Fatalln("Wait", i, err)
		}
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			log.Fatalln("Close", i, err)
		}
	}
}
