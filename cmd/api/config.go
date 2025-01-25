package main

import "os"

type config struct {
	addr string
	db   dbconfig
	mail mailconfig
	auth authconfig
}

type authconfig struct {
	name     string
	password string
}

type mailconfig struct {
	sendgridcfg
}

type sendgridcfg struct {
	apiKey    string
	fromEmail string
}

type dbconfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func Config() *config {
	return &config{
		addr: os.Getenv("addr"),
		db: dbconfig{
			addr:         os.Getenv("DBaddr"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		mail: mailconfig{
			sendgridcfg{
				apiKey:    "api_key",
				fromEmail: "from email",
			},
		},
		auth: authconfig{
			name:     "arshia",
			password: "1234",
		},
	}
}
