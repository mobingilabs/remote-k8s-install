package machine

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"mobingi/ocean/pkg/util/ssh"
)

type CheckFunc func(output string) bool

type Command struct {
	Cmd       string
	Check     CheckFunc
	NeedCheck bool // when false,we don't call CheckFunc
}

type Job struct {
	Name     string
	Commands []Command
}

func NewJob(name string) *Job {
	return &Job{
		Name:     name,
		Commands: []Command{},
	}
}

func (j *Job) AddCmdWithCheck(cmd string, check CheckFunc) {
	c := Command{
		Cmd:       cmd,
		Check:     check,
		NeedCheck: true,
	}

	j.Commands = append(j.Commands, c)
}

func (j *Job) AddCmd(cmd string) {
	c := Command{
		Cmd:       cmd,
		NeedCheck: false,
	}

	j.Commands = append(j.Commands, c)
}

type Machine interface {
	// may be we should use context to manage run,so we can stop it, it is return immeditly, so we don't need run a goroutine
	Run(*Job) error
	Close() error
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
		c: c,
	}, nil
}

func (m *machine) Run(j *Job) error {
	if err := m.checkRunState(); err != nil {
		return err
	}
	atomic.StoreInt32(&m.run, 1)

	return m.doJob(j)
}

func (m *machine) Close() error {
	if err := m.checkRunState(); err != nil {
		return err
	}

	return m.c.Close()
}

func (m *machine) doJob(j *Job) error {
	defer atomic.StoreInt32(&m.run, 0)
	for _, v := range j.Commands {
		output, err := m.c.Do(v.Cmd)
		if err != nil {
			return err
		}
		if v.NeedCheck {
			if !v.Check(output) {
				return fmt.Errorf("check failed, output is:%s", output)
			}
		}
	}

	return nil
}

func (m *machine) checkRunState() error {
	if atomic.LoadInt32(&m.run) == 1 {
		return errors.New("this machine is runing")
	}

	return nil
}
