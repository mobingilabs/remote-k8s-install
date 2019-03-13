package ssh

import (
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

func (c *Client) Close() error {
	return c.Close()
}
