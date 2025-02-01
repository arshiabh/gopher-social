package main

import (
	"os"
	"time"
)

type config struct {
	addr  string
	db    dbconfig
	redis redisconfig
	mail  mailconfig
	auth  authconfig
}

type authconfig struct {
	baseconfig
	jwtconfig
}

type jwtconfig struct {
	secret string
	exp    time.Duration
}

type baseconfig struct {
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

type redisconfig struct {
	addr     string
	password string
	db       int
	enable   bool
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
			baseconfig: baseconfig{
				name:     "arshia",
				password: "1234",
			},
			jwtconfig: jwtconfig{
				secret: "supersecret",
				exp:    time.Hour * 24 * 3,
			},
		},
		redis: redisconfig{
			addr:     "localhost:6379",
			password: "",
			db:       0,
			enable:   false,
		},
	}
}
