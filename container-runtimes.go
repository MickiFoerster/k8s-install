package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var cmds = []string{
	"sudo apt-get update",
	"sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2",
	"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
	`sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`,
	`sudo apt-get update `,
	`sudo apt-get install -y containerd.io=1.2.13-2`,
	`sudo apt-get install -y docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) `,
	`sudo apt-get install -y docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)`,
}

func container_runtimes() {
	log.Println("Install Docker CE")
	for _, cmdstr := range cmds {
		l := strings.Split(cmdstr, " ")
		fmt.Println(l)
		cmd := exec.Command(l[0], l[1:]...)

		if stdout, err := cmd.StdoutPipe(); err == nil {
			go io.Copy(os.Stdout, stdout)
		}
		if stderr, err := cmd.StderrPipe(); err == nil {
			go io.Copy(os.Stderr, stderr)
		}

		if err := cmd.Start(); err != nil {
			log.Fatalf("error: starting command %q failed: %v\n", cmdstr, err)
		}
		log.Printf("command %q started\n", cmdstr)

		if err := cmd.Wait(); err != nil {
			log.Fatalf("error: executing command %q failed: %v\n", cmdstr, err)
		}
		log.Printf("command %q finished successful\n", cmdstr)
	}
}
