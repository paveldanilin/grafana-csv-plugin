package main

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/sftp"
	"os"
	"time"
)

func main () {
	conn := sftp.ConnectionConfig{
		Host:          "",
		Port:          "22",
		User:          "",
		Password:      "",
		Timeout:       time.Second * 20,
		IgnoreHostKey: true,
	}

	workdir, _ := os.Getwd()
	println("/tmp -> " + workdir)

	dwf, err := sftp.GetFile(conn, "/tmp/batman.csv", workdir)
	if err != nil {
		println(err.Error())
		return
	}

	println("-->" + dwf)
}
