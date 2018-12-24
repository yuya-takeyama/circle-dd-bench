package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	flags "github.com/jessevdk/go-flags"
	datadog "github.com/zorkian/go-datadog-api"
)

const appName = "circle-dd-bench"

type options struct {
	Tags        []string `short:"t" long:"tag" description:"Tag to send to Datadog with TAG:VALUE format"`
	ShowVersion bool     `short:"v" long:"version" description:"Show version"`
}

var opts options

var metricName = "circleci.benchmarks.command"
var datadogClient = datadog.NewClient(os.Getenv("DATADOG_API_KEY"), "")

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = appName
	parser.Usage = "[OPTIONS] -- COMMAND"

	args, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	if opts.ShowVersion {
		fmt.Printf("%s v%s, build %s\n", appName, Version, GitCommit)
		return
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	elapsed, err := runCommand(args)
	if err != nil {
		log.Fatalf("%s: failed to run command: %s", appName, err)
	}

	log.Printf("%s: took %.0f seconds", appName, *elapsed)
	metrics := createMetrics(elapsed, opts)

	err = datadogClient.PostMetrics(metrics)
	if err != nil {
		log.Fatalf("%s: failed to send metric to Datadog: %s", appName, err)
	}

	log.Printf("%s: sent a metric to Datadog", appName)
}

func runCommand(args []string) (*float64, error) {
	cmdName := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	begin := time.Now()
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	elapsed := end.Sub(begin)
	duration := float64(elapsed / time.Millisecond)

	return &duration, nil
}

func createMetrics(duration *float64, opts options) []datadog.Metric {
	now := float64(time.Now().Unix())
	tags := append(baseTags(), opts.Tags...)

	return []datadog.Metric{
		datadog.Metric{
			Metric: &metricName,
			Points: []datadog.DataPoint{{&now, duration}},
			Tags:   tags,
		},
	}
}

func baseTags() []string {
	return []string{
		"username:" + os.Getenv("CIRCLE_PROJECT_USERNAME"),
		"reponame:" + os.Getenv("CIRCLE_PROJECT_REPONAME"),
		"branch:" + os.Getenv("CIRCLE_BRANCH"),
		"job:" + os.Getenv("CIRCLE_JOB"),
	}
}
