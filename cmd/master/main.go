package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/nk2ge5k/master"
)

func main() {

	num := flag.Int("n", 1, "num proccess")
	rep := flag.Bool("r", false, "restart process after exit")
	flag.Parse()

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
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	master.Run(ctx, slave, *num, *rep)
}
