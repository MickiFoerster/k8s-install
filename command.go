package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	cmdline    string
	stdin      io.Reader
	outputfile string
}

func (c *Command) pipeexec() error {
	log.Println("start pipe execution")
	pipe := strings.Split(c.cmdline, "|")
	if len(pipe) != 2 {
		log.Fatal("error: pipe longer than 2 commands is not supported")
	}

	left := strings.Split(strings.Trim(pipe[0], " "), " ")
	right := strings.Split(strings.Trim(pipe[1], " "), " ")
	var cmdleft *exec.Cmd
	var cmdright *exec.Cmd
	if len(left) > 1 {
		cmdleft = exec.Command(left[0], left[1:]...)
	} else {
		cmdleft = exec.Command(left[0])
	}
	if len(left) > 1 {
		cmdright = exec.Command(right[0], right[1:]...)
	} else {
		cmdright = exec.Command(right[0])
	}

	pipein, err := cmdleft.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error while creating pipe for %q: %v\n", c.cmdline, err)
	}
	pipeout, err := cmdright.StdinPipe()
	if err != nil {
		return fmt.Errorf("error while creating pipe for %q: %v\n", c.cmdline, err)
	}
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := pipein.Read(buf)
			if err != nil {
				if err == io.EOF {
					pipeout.Close()
					return
				}
				log.Fatalf("error: could not read from pipe: %v\n", err)
			}
			if n > 0 {
				written := 0
				for {
					m, err := pipeout.Write(buf[written:n])
					if err != nil {
						log.Fatalf("error: could not write to pipe: %v\n", err)
					}
					written += m
					if written >= n {
						break
					}
				}
			}
		}
	}()
	cmdleft.Stderr = os.Stderr
	cmdright.Stderr = os.Stderr
	cmdright.Stdout = os.Stdout

	if err := cmdleft.Start(); err != nil {
		return fmt.Errorf("error: starting left command of pipe %v failed: %v\n", left, err)
	}
	if err := cmdright.Start(); err != nil {
		return fmt.Errorf("error: starting right command of pipe %v failed: %v\n", right, err)
	}

	if err := cmdleft.Wait(); err != nil {
		return fmt.Errorf("error: left command of pipe %v failed: %v\n", left, err)
	}
	if err := cmdright.Wait(); err != nil {
		return fmt.Errorf("error: right command of pipe %v failed: %v\n", right, err)
	}

	return nil
}

func (c *Command) exec() error {
	if strings.Index(c.cmdline, "|") != -1 {
		return c.pipeexec()
	}
	l := strings.Split(c.cmdline, " ")
	fmt.Println(l)
	cmd := exec.Command(l[0], l[1:]...)
	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}
	if c.outputfile != "" {
		f, err := os.Create(c.outputfile)
		if err != nil {
			return fmt.Errorf("error: could not open file %q: %v", c.outputfile, err)
		}
		cmd.Stdout = f
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error: starting command %q failed: %v\n", c.cmdline, err)
	}
	log.Printf("command %q started\n", c.cmdline)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error: command %q failed: %v\n", c.cmdline, err)
	}
	log.Printf("command %q finished successful\n", c.cmdline)

	return nil
}
