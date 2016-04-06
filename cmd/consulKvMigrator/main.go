package main

import "flag"
import "fmt"
import "os"
import migrator "github.com/tshak/consulKvMigrator"
import consul "github.com/hashicorp/consul/api"

func main() {
	prompt := flag.Bool("prompt", true, "prompt before submitting changes to consul")
	dryRun := flag.Bool("dry-run", false, "don't submit any changes to consul")
	address := flag.String("address", "localhost:8500", "address for consul")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] inputFile\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	inputFile := flag.Arg(0)

	if len(inputFile) == 0 {
		fmt.Printf("Must specify an inputFile.\n")
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Unable to read file: %s\n", inputFile)
		os.Exit(1)
	}

	consulConfig := &consul.Config{
		Address: *address,
	}
	client, err := consul.NewClient(consulConfig)

	if err != nil {
		panic(err)
	}

	if err := migrator.Migrate(client.KV(), inputFile, *prompt, *dryRun); err != nil {
		panic(err)
	}
}
