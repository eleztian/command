package command

import (
	"context"
	"io"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"
)

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
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

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// Timeout
	Timeout time.Duration
}

// ConcurrenceComE concurrence run command
// if any command has return error, all command will been kill
// return the first error.
func ConcurrenceComE(ctx context.Context, cmds ...*Cmd) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, c := range cmds {
		var (
			ctxC    = ctx
			cancelC context.CancelFunc
		)
		if c.Timeout != 0 {
			ctxC, cancelC = context.WithTimeout(ctx, c.Timeout)
		}
		cmd := exec.CommandContext(ctxC, c.Path, c.Args...)
		cmd.Dir = c.Dir
		cmd.Env = c.Env
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr
		cmd.Stdin = c.Stdin
		eg.Go(func() (err error) {
			err = cmd.Run()
			if cancelC != nil {
				cancelC()
			}
			return
		})
	}
	return eg.Wait()
}

// ConcurrenceComNE concurrence run command
// return the first error.
func ConcurrenceComNE(ctx context.Context, cmds ...*Cmd) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, c := range cmds {
		var (
			ctxC    context.Context
			cancelC context.CancelFunc
		)
		if c.Timeout != 0 {
			ctxC, cancelC = context.WithTimeout(ctx, c.Timeout)
		} else {
			ctxC, cancelC = context.WithCancel(ctx)
		}
		cmd := exec.CommandContext(ctxC, c.Path, c.Args...)
		cmd.Dir = c.Dir
		cmd.Env = c.Env
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr
		cmd.Stdin = c.Stdin
		eg.Go(func() (err error) {
			err = cmd.Run()
			cancelC()
			return
		})
	}
	return eg.Wait()
}

func NewCmd(name string, timeout time.Duration, args ...string) *Cmd {
	return &Cmd{
		Path:    name,
		Args:    args,
		Timeout: timeout,
	}
}
