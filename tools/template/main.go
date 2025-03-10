package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type options struct {
	templates     []string
	data          string
	dataFormat    string
	output        string
	configuration bool
	delims        string
}

func parseFlags() options {
	var opts options
	var templatesList string

	flag.StringVar(&templatesList, "template", "", "The template files to use (comma-separated)")
	flag.StringVar(&opts.data, "data", "", "The data file to use")
	flag.StringVar(&opts.dataFormat, "data-format", "json", "The data format to use (json or yaml)")
	flag.StringVar(&opts.output, "output", "", "The output file to create")
	flag.StringVar(&opts.delims, "delims", "{{.}}", "Template delimiters format (e.g., '[[.]]' or '{{.}}')")

	flag.Parse()

	if templatesList == "" {
		log.Fatal("--template is required")
	}
	if opts.output == "" {
		log.Fatal("--output is required")
	}

	opts.templates = strings.Split(templatesList, ",")
	return opts
}

func parseDelims(format string) (left, right string, err error) {
	dotIndex := strings.Index(format, ".")
	if dotIndex == -1 {
		return "", "", fmt.Errorf("delimiter format must contain a '.' (e.g., '[[.]]')")
	}
	left = format[:dotIndex]
	right = format[dotIndex+1:]
	if left == "" || right == "" {
		return "", "", fmt.Errorf("both left and right delimiters must be specified")
	}
	return left, right, nil
}

func main() {
	opts := parseFlags()

	// Parse delimiters
	leftDelim, rightDelim, err := parseDelims(opts.delims)
	if err != nil {
		log.Fatalf("invalid delimiter format: %v", err)
	}

	// Use to do operations once and only once.
	cache := map[string]bool{}
	doOnce := func(key string) bool {
		if _, ok := cache[key]; ok {
			return false // Already done.
		}
		cache[key] = true
		return true
	}

	// Read the template file
	funcMap := sprig.TxtFuncMap()
	funcMap["doOnce"] = doOnce
	for k, v := range customFuncMap {
		funcMap[k] = v
	}
	// Parse the template
	tmpl := template.New("template").Funcs(funcMap).Delims(leftDelim, rightDelim)
	for _, templatePath := range opts.templates {
		bytes, err := os.ReadFile(templatePath)
		if err != nil {
			log.Fatalf("reading template file: %v", err)
		}
		tmpl, err = tmpl.Parse(string(bytes))
		if err != nil {
			log.Fatalf("parsing template: %v", err)
		}
	}

	// Read the data file
	data := map[string]any{}
	if opts.data != "" {
		dataBytes, err := os.ReadFile(opts.data)
		if err != nil {
			log.Fatalf("reading data file: %v", err)
		}
		fixedDataBytes := bytes.ReplaceAll(dataBytes, []byte("True"), []byte("true"))
		fixedDataBytes = bytes.ReplaceAll(fixedDataBytes, []byte("False"), []byte("false"))

		// Unmarshal the data into a map
		switch opts.dataFormat {
		case "json":
			if err := json.Unmarshal(fixedDataBytes, &data); err != nil {
				log.Fatalf("unmarshaling json data: %v", err)
			}
		case "yaml":
			if err := yaml.Unmarshal(fixedDataBytes, &data); err != nil {
				log.Fatalf("unmarshaling yaml data: %v", err)
			}
		default:
			log.Fatalf("unknown data format: %s", opts.dataFormat)
		}
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("executing template: %v", err)
	}
	// Write the result to the output file
	if err := os.WriteFile(opts.output, buf.Bytes(), 0644); err != nil {
		log.Fatalf("writing output file: %v", err)
	}
	log.Printf("Successfully processed template and data")
}
