package ssh

import (
	"strings"
	"time"

	stdssh "golang.org/x/crypto/ssh"
)

type Client struct {
	conn *stdssh.Client
}

func NewClient(addr, user, password string) (*Client, error) {
	config := &stdssh.ClientConfig{
		User: user,
		Auth: []stdssh.AuthMethod{
			stdssh.Password(password),
		},
		HostKeyCallback: stdssh.InsecureIgnoreHostKey(),
	}

	client, err := stdssh.Dial("tcp", addr+":22", config)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: client,
	}, nil
}

// Do exec cmd on the romote machine and return std output
// TODO some cmd don't need stdout,some need. split it to two func
func (c *Client) Do(cmd string) (string, error) {
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

func (c *Client) DoWithoutOutput(cmd string) error {
	sess, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	if strings.HasPrefix(cmd, "systemctl") {
		time.Sleep(30 * time.Second)
	}

	return sess.Run(cmd)
}

func (c *Client) Close() error {
	return c.Close()
}
