package tasks

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/kubermatic/kubeone/pkg/ssh"
)

// tee mimics the unix `tee` command by piping its
// input through to the upstream writer and also
// capturing it in a buffer.
type tee struct {
	buffer   bytes.Buffer
	upstream io.Writer
}

func (t *tee) Write(p []byte) (int, error) {
	t.buffer.Write(p)

	return t.upstream.Write(p)
}

func (t *tee) String() string {
	return strings.TrimSpace(t.buffer.String())
}

func runCommand(conn ssh.Connection, cmd string, verbose bool) (string, string, int, error) {
	if !verbose {
		return conn.Exec(cmd)
	}

	stdout := &tee{
		upstream: os.Stdout,
	}

	stderr := &tee{
		upstream: os.Stdout,
	}

	exitCode, err := conn.Stream(cmd, stdout, os.Stderr)

	return stdout.String(), stderr.String(), exitCode, err
}