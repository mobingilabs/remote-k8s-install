package ssh

import (
	"fmt"
	"io"
	"mobingi/ocean/pkg/log"
	"os"
	"time"

	stdssh "golang.org/x/crypto/ssh"
)

const (
	dialTimeout = 10 * time.Second
)

// Client is a conn to remote machine
type Client interface {
	Do(cmd string) (string, error)
	Close() error
	SCP(localPath, remotePath string) error
}

type client struct {
	conn *stdssh.Client
}

func NewClient(addr, user, password string) (Client, error) {
	config := &stdssh.ClientConfig{
		User: user,
		Auth: []stdssh.AuthMethod{
			stdssh.Password(password),
		},
		HostKeyCallback: stdssh.InsecureIgnoreHostKey(),
		Timeout:         dialTimeout,
	}

	sshClient, err := stdssh.Dial("tcp", addr+":22", config)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: sshClient,
	}, nil
}

// Do exec cmd on the romote machine and return std output
// TODO some cmd don't need stdout,some need. split it to two func
func (c *client) Do(cmd string) (string, error) {
	sess, err := c.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	output, err := sess.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) SCP(localPath, remotePath string) error {
	fi, err := os.Stat(localPath)
	log.Info(fi)
	if err != nil {
		return err
	}

	sess, err := c.conn.NewSession()
	log.Info("sess new")
	if err != nil {
		return err
	}
	defer sess.Close()

	w, err := sess.StdinPipe()
	log.Info("ssh stdin")
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "C0644", fi.Size(), fi.Name()); err != nil {
		return err
	}
	log.Info("write file meta info to remote")

	f, err := os.Open(localPath)
	log.Info("open local file")
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, f); err != nil {
		log.Infof("%s is copying to remote", fi.Name())
		return err
	}

	if _, err := fmt.Fprintln(w, "\x00"); err != nil {
		return nil
	}

	return nil
}
