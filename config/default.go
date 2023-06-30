package config

const defaultJson string = `{
	"service": {
		"address": ":7075",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "trace",
		"file": "logs/server.log",
		"std": false
	},
	"database": {
		"name": "rgsCloud",
		"ip": "192.168.1.10",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	},
	"basic": {
		"synonym": 6,
		"tag": 6
	}
}
`
