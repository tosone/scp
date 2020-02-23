Copy files over SCP with Go
=============================
[![Go Report Card](https://goreportcard.com/badge/bramvdbogaerde/go-scp)](https://goreportcard.com/report/bramvdbogaerde/go-scp) [![](https://godoc.org/github.com/bramvdbogaerde/go-scp?status.svg)](https://godoc.org/github.com/bramvdbogaerde/go-scp)

This package makes it very easy to copy files over scp in Go.
It uses the golang.org/x/crypto/ssh package to establish a secure connection to a remote server in order to copy the files via the SCP protocol.

### Example usage

```go
package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/tosone/go-scp"
)

func main() {
	var sshClient = &ssh.ClientConfig{
		User: "tosone",
		Auth: []ssh.AuthMethod{
			ssh.Password("123456"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var client = scp.Client{
		Host:         net.JoinHostPort("192.168.1.100", "22"),
		ClientConfig: sshClient,
	}

	var err error

	const filename = "test.txt"
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(filename); err != nil {
		log.Fatal(err)
	}
	var fileReader *os.File
	if fileReader, err = os.Open(filename); err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(); err != nil {
		log.Fatal(err)
	}

	if err = client.Copy(fileReader, "/home/test/test.txt", 0644, fileInfo.Size()); err != nil {
		log.Fatal(err)
	}
	if err = client.Close(); err != nil {
		log.Fatal(err)
	}
}
```

### License

This library is licensed under the Mozilla Public License 2.0.    
A copy of the license is provided in the `LICENSE.txt` file.

Copyright (c) 2020 Bram Vandenbogaerde
