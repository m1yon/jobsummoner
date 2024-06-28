package main

import (
	"flag"
)

type config struct {
	useLocalDB  bool
	dsn         string
	proxyConfig proxyConfig
}

func getConfigFromFlags() *config {
	useLocalDB := flag.Bool("local-db", true, "Use a local sqlite DB")
	dsn := flag.String("dsn", "", "Database connection string")

	proxyHostname := flag.String("proxy-hostname", "", "hostname of the proxy server")
	proxyPort := flag.String("proxy-port", "", "port of the proxy server")
	proxyUsername := flag.String("proxy-username", "", "username of the proxy server")
	proxyPassword := flag.String("proxy-password", "", "password of the proxy server")

	flag.Parse()

	proxy := proxyConfig{
		Hostname: *proxyHostname,
		Port:     *proxyPort,
		Username: *proxyUsername,
		Password: *proxyPassword,
	}

	return &config{useLocalDB: *useLocalDB, dsn: *dsn, proxyConfig: proxy}
}
