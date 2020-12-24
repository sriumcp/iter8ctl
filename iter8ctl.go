package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Useful for unit testing the CLI
var out io.Writer = os.Stdout

// Rules for valid k8s resource name and namespace are here: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
var nameRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-|\.){0,253}([[:lower:]]|[[:digit:]])$`)
var namespaceRegex *regexp.Regexp = regexp.MustCompile(`^([[:lower:]]|[[:digit:]])([[:lower:]]|[[:digit:]]|-){0,63}([[:lower:]]|[[:digit:]])$`)

// Exact match required for apiVersion: v2alpha1 is the only allowed value
var apiVersionRegex *regexp.Regexp = regexp.MustCompile(`\bv2alpha1\b`)

func main() {

	describeCmd := flag.NewFlagSet("describe", flag.ExitOnError)
	experimentName := describeCmd.String("name", "", "experiment name")
	experimentNamespace := describeCmd.String("namespace", "default", "experiment namespace")
	apiVersion := describeCmd.String("apiVersion", "v2alpha1", "experiment api version")

	if len(os.Args) < 2 {
		fmt.Fprintln(out, "expected 'describe' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "describe":
		describeCmd.Parse(os.Args[2:])

		if len(*experimentName) == 0 || len(*experimentName) > 253 || !nameRegex.MatchString(*experimentName) || len(*experimentNamespace) == 0 || len(*experimentNamespace) > 63 || !namespaceRegex.MatchString(*experimentNamespace) {
			fmt.Fprintln(out, "Expected a valid value for (experiment) name and namespace.")
			fmt.Fprintln(out, "Name should: contain no more than 253 characters and only lowercase alphanumeric characters, '-' or '.'; start and end with an alphanumeric character.")
			fmt.Fprintln(out, "Namespace should: contain no more than 63 characters and only lowercase alphanumeric characters, or '-'; start and end with an alphanumeric character.")

			fmt.Fprintln(out, "\nYou supplied...")
			fmt.Fprintln(out, "  name:", *experimentName)
			fmt.Fprintln(out, "  namespace:", *experimentNamespace)

			os.Exit(1)
		}

		if !apiVersionRegex.MatchString(*apiVersion) {
			fmt.Fprintln(out, "Expected a valid value for (experiment) api version.")
			fmt.Fprintln(out, "Only allowed value for api version is 'v2alpha1'.")

			fmt.Fprintln(out, "\nYou supplied...")
			fmt.Fprintln(out, "  apiVersion:", *apiVersion)

			os.Exit(1)
		}

	default:
		fmt.Fprintln(out, "expected 'describe' subcommand")
		os.Exit(1)
	}
}
