package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/zeebo/errs"
)

func handle(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"cat"}
	}

	var (
		cmd *exec.Cmd
		wc  io.WriteCloser
	)

	initCmd := func() (err error) {
		cmd = exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		wc, err = cmd.StdinPipe()
		if err != nil {
			return errs.Wrap(err)
		}
		return errs.Wrap(cmd.Start())
	}

	handle(initCmd())

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(bytes.TrimSpace(scanner.Bytes())) == 0 {
			wc.Close()
			cmd.Wait()
			fmt.Println()
			handle(initCmd())
			continue
		}
		_, err := fmt.Fprintln(wc, scanner.Text())
		handle(err)
	}
	handle(errs.Wrap(scanner.Err()))
}
