package master

import (
	"context"
	"io"
	_ "net/http/pprof"
	"os/exec"
	"sync"
)

// Slave represents exec.Cmd struct without underlying process
type Slave struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If Stdout and Stderr are the same writer, and have a type that can be compared with ==,
	// at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	sema chan struct{}
}

// command creates new exec.Cmd form the slave
func (s *Slave) command(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, s.Path, s.Args...)

	cmd.Env = s.Env
	cmd.Dir = s.Dir

	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr

	return cmd
}

// Run runs command produced by slave N times and repeats if necessary
func Run(ctx context.Context, s *Slave, n int, repeat bool) {
	queue := make(chan *exec.Cmd, n)
	defer close(queue)

	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		// populate queue with n commands
		queue <- s.command(ctx)

		wg.Add(1)
		go func() {
			defer wg.Done()

		run:
			for {
				select {
				case <-ctx.Done():
					break run
				case cmd := <-queue:
					// run command blocking
					if err := cmd.Run(); err != nil {
						// if was any errors than stop runner
						break run
					}

					// if command shoud be repeated than we will create new
					// command and send it into the queue
					if repeat {
						queue <- s.command(ctx)
					} else {
						break run
					}
				}
			}
		}()
	}

	wg.Wait()
}
