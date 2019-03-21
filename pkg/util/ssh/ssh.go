package ssh

import (
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
