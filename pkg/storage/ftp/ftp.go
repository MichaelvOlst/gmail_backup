package ftp

import (
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

const name = "ftp"

// Config holds the config the ftp option
type Config struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Provider implements storage.Provider for the ftp file storage.
type Provider struct {
	client *ftp.ServerConn
}

// Name returns ftp
func (p *Provider) Name() string {
	return name
}

// New initializer for Provider struct ftp
func New(cfg Config) *Provider {

	c, err := ftp.Dial(cfg.Host, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(cfg.Username, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	// data := bytes.NewBufferString("Hello World")
	// err = c.Stor("Public/test-file.txt", data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// f, err := os.Open("ftp.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = c.Stor("/Public/ftp.txt", f)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%+v\n", c)

	// err = c.ChangeDir("/Public")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = c.MakeDir("/ftp")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// c.ChangeDirToParent()
	// err = c.ChangeDir("/Public")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = c.MakeDir("/Public/bla")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// entries, err := c.List("/Public/")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// entries, _ := c.List("/Public")
	// fmt.Printf("%+v\n", entries)
	// pwd, _ := c.CurrentDir()
	// fmt.Println(pwd)

	return &Provider{c}
}
