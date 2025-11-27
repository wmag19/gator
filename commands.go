package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("commmand not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) { //register is a method on commands that is passed a name and a function.
	if c.registeredCommands == nil {
		c.registeredCommands = make(map[string]func(*state, command) error)
	}
	c.registeredCommands[name] = f
}
