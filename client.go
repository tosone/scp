// Package scp ...
package scp

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/tosone/logging"
)

// Binary ..
var Binary = "/usr/bin/scp"

type Client struct {
	Host         string            // scp target host
	ClientConfig *ssh.ClientConfig // ssh client
	session      *ssh.Session      // ssh connection session
	conn         ssh.Conn          // ssh connection
}

// Connect connect to the remote SSH server, returns error if it couldn't establish a session to the SSH server
func (c *Client) Connect() (err error) {
	var client *ssh.Client
	if client, err = ssh.Dial("tcp", c.Host, c.ClientConfig); err != nil {
		return
	}

	c.conn = client.Conn
	if c.session, err = client.NewSession(); err != nil {
		return
	}
	return
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// CopyFile ..
func (c *Client) CopyFile(filename string, remotePath string, permission os.FileMode) (err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(filename); err != nil {
		return
	}
	var fileReader *os.File
	if fileReader, err = os.Open(filename); err != nil {
		return
	}
	err = c.Copy(fileReader, remotePath, permission, fileInfo.Size())
	return
}

// CopyFileTimeout ..
func (c *Client) CopyFileTimeout(filename string, remotePath string, permission os.FileMode, timeout time.Duration) (err error) {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = c.CopyFile(filename, remotePath, permission)
	}()
	if waitTimeout(wg, timeout) {
		err = fmt.Errorf("123")
	}
	return
}

// CopyTimeout ..
func (c *Client) CopyTimeout(reader io.Reader, remotePath string, permission os.FileMode, size int64, timeout time.Duration) (err error) {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = c.Copy(reader, remotePath, permission, size)
	}()
	if waitTimeout(wg, timeout) {
		err = fmt.Errorf("123")
	}
	return
}

// Copy Copies the contents of an io.Reader to a remote location
func (c *Client) Copy(reader io.Reader, remotePath string, permission os.FileMode, size int64) (err error) {
	var filename = path.Base(remotePath)

	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		var writer io.WriteCloser
		if writer, err = c.session.StdinPipe(); err != nil {
			return
		}

		defer func() {
			if err := writer.Close(); err != nil {
				logging.Error(err)
			}
		}()

		if _, err = fmt.Fprintln(writer, "C0"+strconv.FormatUint(uint64(permission), 10), size, filename); err != nil {
			return
		}

		if _, err = io.Copy(writer, reader); err != nil {
			return
		}

		_, err = fmt.Fprint(writer, "\x00")
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		err := c.session.Run(fmt.Sprintf("%s -rqt %s", Binary, remotePath))
		if err != nil {
			return
		}
	}()

	wg.Wait()

	return
}

// Close ..
func (c *Client) Close() (err error) {
	if err = c.conn.Close(); err != nil {
		return
	}
	if err = c.session.Close(); err != io.EOF && err != nil {
		return
	}
	return
}
