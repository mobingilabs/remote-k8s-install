package machine

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"mobingi/ocean/pkg/util/ssh"
)

type CheckFunc func(output string) bool

// TODO we will use job to combine some commands together
type Job struct {
	Name     string
	Commands []Command
}

type CommandList []Command

func (cl *CommandList) Add(cmd string, check CheckFunc) {
	command := Command{cmd, check}
	*cl = append([]Command(*cl), command)
}

func (cl *CommandList) AddAnother(another CommandList) {
	*cl = append([]Command(*cl), []Command(another)...)
}

type Command struct {
	Cmd   string
	Check CheckFunc
}

type Machine interface {
	AddCommand(command Command) error
	AddCommandList(commandList CommandList) error
	Run() error
	Reset() error
	DisConnect() error
	SCP(localPath, remotePath string) error
}

type machine struct {
	c ssh.Client

	sync.Mutex
	commands []Command

	run int32
}

func NewMachine(addr, user, password string) (Machine, error) {
	c, err := ssh.NewClient(addr, user, password)
	if err != nil {
		return nil, err
	}

	return &machine{
		c:        c,
		commands: make([]Command, 0),
	}, nil
}

func (m *machine) AddCommand(command Command) error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	m.Lock()
	m.commands = append(m.commands, command)
	m.Unlock()

	return nil
}

func (m *machine) AddCommandList(commandList CommandList) error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	for _, command := range []Command(commandList) {
		m.commands = append(m.commands, command)
	}

	return nil
}

func (m *machine) Run() error {
	if err := m.checkRunState(); err != nil {
		return err
	}
	defer m.Reset()

	atomic.StoreInt32(&m.run, 1)
	defer atomic.StoreInt32(&m.run, 0)
	for _, command := range m.commands {
		if err := m.docmdAndCheck(command); err != nil {
			return err
		}
	}

	return nil
}

func (m *machine) DisConnect() error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	return m.c.Close()
}

func (m *machine) docmdAndCheck(command Command) error {
	output, err := m.c.Do(command.Cmd)
	if err != nil {
		return fmt.Errorf("cmd:%s,err:%s", command.Cmd, err.Error())
	}

	if !command.Check(output) {
		return fmt.Errorf("check failed, output is :%s", output)
	}

	return nil
}

func (m *machine) Reset() error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	m.commands = m.commands[0:0]
	return nil
}

func (m *machine) checkRunState() error {
	if atomic.LoadInt32(&m.run) == 1 {
		return errors.New("this machine is runing")
	}

	return nil
}

func (m *machine) SCP(localPath, remotePath string) error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	return m.c.SCP(localPath, remotePath)
}
