package models

// Settings is the model for settings of this app
type Settings struct {
	Ftp FtpOption `json:"ftp"`
}

// FtpOption holds the config for Ftp
type FtpOption struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}
