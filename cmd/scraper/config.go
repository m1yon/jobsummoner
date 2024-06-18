package main

import "flag"

type config struct {
	useLocalDB bool
}

func getConfigFromFlags() *config {
	useLocalDB := flag.Bool("local-db", true, "Use a local sqlite DB")

	flag.Parse()

	return &config{useLocalDB: *useLocalDB}
}
