package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/nk2ge5k/master"
	"github.com/pkg/errors"
)

func main() {

	var (
		num = flag.Int("n", 1, "num proccess")
		rep = flag.Bool("r", false, "restart process after exit")
		out = flag.String("out", "", "process stdout")
		err = flag.String("err", "", "process stderr")
	)
	flag.Parse()

	outw, err := openWriter(*out, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	errw, err := openWriter(*err, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	args := flag.Args()

	if len(args) < 1 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig // cancel context on interrupt signal
		cancel()
	}()

	slave := &master.Slave{
		Path:   args[0],
		Args:   args[1:],
		Stdout: outw,
		Stderr: errw,
	}

	master.Run(ctx, slave, *num, *rep)
}

func openTeeWriter(path string, fallback io.Writer) (io.Writer, error) {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fallback, errors.Wrap(err, "open writer")
	}
	return io.MultiWriter(w, fallback), nil
}
