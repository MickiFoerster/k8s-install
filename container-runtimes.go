package main

import (
	"log"
	"strings"
)

var cmds = []Command{
	{cmdline: `sudo apt-get update`},
	{cmdline: `sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2`},
	{cmdline: `curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -`},
	{cmdline: `sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`},
	{cmdline: `sudo apt-get update `},
	{cmdline: `sudo apt-get install -y containerd.io=1.2.13-2`},
	{cmdline: `sudo apt-get install -y docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) `},
	{cmdline: `sudo apt-get install -y docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)`},
	{
		cmdline: `sudo cat`,
		stdin: strings.NewReader(`
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
`),
		outputfile: "/etc/docker/daemon.json",
	},
	{cmdline: `sudo mkdir -p /etc/systemd/system/docker.service.d`},
	{cmdline: `sudo systemctl daemon-reload`},
	{cmdline: `sudo systemctl restart docker`},
	{cmdline: `sudo systemctl enable docker`},
}

func container_runtimes() {
	log.Println("Install Docker CE")
	for _, c := range cmds {
		err := c.exec()
		if err != nil {
			log.Fatal(err)
		}
	}
}
