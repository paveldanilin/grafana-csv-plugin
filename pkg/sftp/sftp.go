package sftp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ConnectionConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	Timeout       time.Duration
	IgnoreHostKey bool
}

func Test(conn ConnectionConfig) error {
	c, err := newConnection(conn)
	if err != nil {
		return err
	}
	_ = c.Close()
	return nil
}

func GetFile(conn ConnectionConfig, fileName string, targetDir string) (string, error) {
	c, err := newConnection(conn)
	if err != nil {
		return "", err
	}
	defer c.Close()

	client, err := sftp.NewClient(c)
	if err != nil {
		return "", err
	}
	defer client.Close()

	localFilename := filepath.Base(fileName)

	// Dest file
	dstFile, err := os.Create(filepath.Join(targetDir, localFilename))
	if err != nil {
		return "", err
	}

	srcFile, err := client.Open(fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", err
	}

	err = dstFile.Sync()
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, localFilename), nil
}

func newConnection(conn ConnectionConfig) (*ssh.Client, error) {
	config, err := newShhClientConfig(conn)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", conn.Host, conn.Port), config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newShhClientConfig(conn ConnectionConfig) (*ssh.ClientConfig, error) {
	hostKeyCallback := ssh.InsecureIgnoreHostKey()

	if conn.IgnoreHostKey == false {
		hostKey, err := getHostKey(conn.Host)
		if err != nil {
			return nil, err
		}
		hostKeyCallback = ssh.FixedHostKey(*hostKey)
	}

	return &ssh.ClientConfig{
		User:              conn.User,
		Auth:              []ssh.AuthMethod{
			ssh.Password(conn.Password),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         conn.Timeout,
	}, nil
}

func getHostKey(host string) (*ssh.PublicKey, error) {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
			}
			break
		}
	}

	if hostKey == nil {
		return nil, errors.New(fmt.Sprintf("no hostkey found for %s", host))
	}

	return &hostKey, nil
}
