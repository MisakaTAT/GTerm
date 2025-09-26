package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"golang.org/x/crypto/ssh"
)

type Result struct {
	stdout   string
	stderr   string
	err      error
	exitCode int
}

type Adapter struct {
	client  *ssh.Client
	timeout time.Duration
}

var (
	DefaultTimeout = 30 * time.Second
	ErrTimeout     = errors.New("command execution timeout")
)

func New(client *ssh.Client) *Adapter {
	return &Adapter{
		client:  client,
		timeout: DefaultTimeout,
	}
}

func (a *Adapter) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		a.timeout = timeout
	}
}

func (a *Adapter) Run(cmd string) *Result {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	session, err := a.client.NewSession()
	if err != nil {
		return &Result{err: fmt.Errorf("failed to create session: %v", err)}
	}
	defer func() {
		_ = session.Close()
	}()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	errCh := make(chan error, 1)
	go func() {
		errCh <- session.Run(cmd)
	}()

	select {
	case err = <-errCh:
		result := &Result{
			stdout: strings.TrimSpace(stdout.String()),
			stderr: strings.TrimSpace(stderr.String()),
			err:    err,
		}

		if err != nil {
			var exitErr *ssh.ExitError
			if errors.As(err, &exitErr) {
				result.exitCode = exitErr.ExitStatus()
			}
		}

		return result
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		time.Sleep(100 * time.Millisecond)
		_ = session.Signal(ssh.SIGKILL)
		return &Result{err: ErrTimeout}
	}
}

func (a *Adapter) RunWithTimeout(cmd string, timeout time.Duration) *Result {
	oldTimeout := a.timeout
	a.timeout = timeout
	defer func() { a.timeout = oldTimeout }()
	return a.Run(cmd)
}

func (a *Adapter) RunScript(script string) *Result {
	return a.Run(fmt.Sprintf("bash -c '%s'", strings.ReplaceAll(script, "'", "'\\''")))
}

func (r *Result) Unwrap() string {
	if r.err != nil {
		return ""
	}
	return r.stdout
}

func (r *Result) Error() error {
	if r.err == nil {
		return nil
	}
	if r.stderr != "" {
		return fmt.Errorf("%v: %s", r.err, r.stderr)
	}
	return r.err
}

func (r *Result) ExitCode() int {
	return r.exitCode
}

func (r *Result) StdOut() string {
	return r.stdout
}

func (r *Result) StdErr() string {
	return r.stderr
}

func (r *Result) Success() bool {
	return r.err == nil && r.exitCode == 0
}

func (r *Result) String() string {
	if r.Success() {
		return r.stdout
	}
	if r.stderr != "" {
		return fmt.Sprintf("error: %v, stderr: %s", r.err, r.stderr)
	}
	return fmt.Sprintf("error: %v", r.err)
}

func NewExec(conf *commonssh.Config) (*ssh.Client, error) {
	return commonssh.NewSSHClient(conf)
}
