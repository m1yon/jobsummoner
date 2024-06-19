package main

import "flag"

type config struct {
	useLocalDB bool
	dsn        string
}

func getConfigFromFlags() *config {
	useLocalDB := flag.Bool("local-db", true, "Use a local sqlite DB")
	dsn := flag.String("dsn", "", "Database connection string")

	flag.Parse()

	return &config{useLocalDB: *useLocalDB, dsn: *dsn}
}
