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

	c, err := ftp.Dial(cfg.Host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(cfg.Username, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	return &Provider{}
}
