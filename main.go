package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/DHowett/go-plist"
)

type Sequel struct {
	ContentFilters struct{} `plist:"ContentFilters"`
	AutoConnect    bool     `plist:"auto_connect"`
	Data           Data     `plist:"data"`
	Encrypted      bool     `plist:"encrypted"`
	Format         string   `plist:"format"`
	Version        int      `plist:"version"`
	QueryFavorites []string `plist:"queryFavorites"`
	QueryHistory   []string `plist:"queryHistory"`
}

type Data struct {
	Connection interface{} `plist:"connection"`
}

type SSHConnection struct {
	Host                              string `plist:"host"`
	User                              string `plist:"user"`
	Password                          string `plist:"password"`
	SSHHost                           string `plist:"ssh_host"`
	SSHKeyLoocation                   string `plist:"ssh_keyLocation"`
	SSHKeyLoocationEnabled            int    `plist:"ssh_keyLocationEnabled"`
	SSHPassword                       string `plist:"ssh_password"`
	SSHUser                           string `plist:"ssh_user"`
	SSLCACertFileLocation             string `plist:"sslCACertFileLocation"`
	SSLCACertFileLocationEnabled      int    `plist:"sslCACertFileLocationEnabled"`
	SSLCertificateFileLocation        string `plist:"sslCertificateFileLocation"`
	SSLCertificateFileLocationEnabled int    `plist:"sslCertificateFileLocationEnabled"`
	SSLKeyFileLocation                string `plist:"sslKeyFileLocation"`
	SSLKeyFileLocationEnabled         int    `plist:"sslKeyFileLocationEnabled"`
	Type                              string `plist:"type"`
	UseSSL                            int    `plist:"useSSL"`
}

var (
	host     = flag.String("host", "", "host")
	user     = flag.String("user", "", "user")
	password = flag.String("password", "", "password")

	docker = flag.Bool("docker", false, "docker")

	sshPassword = flag.String("ssh-password", "", "ssh password")
	sshUser     = flag.String("ssh-user", "", "ssh user")
	sshHost     = flag.String("ssh-host", "", "ssh host")
)

func main() {
	flag.Parse()

	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	if *docker {
		if err := overwriteViaDocker(); err != nil {
			return err
		}
	}

	f, err := ioutil.TempFile("", "docker-sequel-open")
	if err != nil {
		return err
	}
	// require .spf extension
	if err := os.Rename(f.Name(), f.Name()+".spf"); err != nil {
		return err
	}
	f.Close()

	f, err = os.Create(f.Name() + ".spf")
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if err := buildPlist(f); err != nil {
		return err
	}

	cmd := exec.Command("open", "-b", "com.sequelpro.SequelPro", f.Name())
	if err := cmd.Run(); err != nil {
		return err
	}

	//TODO cant detect finish to connect
	time.Sleep(time.Second * 3)

	return nil
}

func buildPlist(f *os.File) error {
	enc := plist.NewEncoder(f)
	sequel := &Sequel{
		AutoConnect: true,
		Data: Data{
			Connection: SSHConnection{
				Host:     *host,
				User:     *user,
				Password: *password,

				SSHHost:     *sshHost,
				SSHUser:     *sshUser,
				SSHPassword: *sshPassword,

				Type: "SPSSHTunnelConnection",
			},
		},
		Encrypted: false,
		Format:    "connection",
		Version:   1,
	}
	enc.Indent("\t")
	return enc.Encode(sequel)
}

func overwriteViaDocker() error {
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)
	*sshHost = re.FindString(os.Getenv("DOCKER_HOST"))
	if *sshHost == "" {
		return fmt.Errorf("cant detect docker host ip from DOCKER_HOST")
	}

	cmd := exec.Command("docker", "inspect", "--format='{{.NetworkSettings.IPAddress}}'", *host)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	*host = strings.TrimRight(string(output), "\n")

	return nil
}
