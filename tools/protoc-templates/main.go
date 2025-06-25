package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	_ "proto/aip"
	_ "proto/codegen/admin_api"
	_ "proto/codegen/model"
	_ "proto/codegen/rpc"
)

var (
	opts struct {
		Debug         *bool
		Template      *string
		Configuration *string
	}

	//go:embed ./**/*.tmpl
	templateFS embed.FS
)

type Input struct {
	File          *protogen.File
	Files         []*protogen.File
	GeneratedFile *protogen.GeneratedFile
	Configuration map[any]any
}

func main() {
	var flags flag.FlagSet
	opts.Debug = flags.Bool("debug", false, "verbose output")
	opts.Template = flags.String("template", "", "template to compile")
	opts.Configuration = flags.String("configuration", "", "configuration to inject in context")
	options := protogen.Options{
		ParamFunc: flags.Set,
	}
	options.Run(func(gen *protogen.Plugin) error {
		var configuration map[any]any
		if *opts.Configuration != "" {
			configData, err := os.ReadFile(*opts.Configuration)
			if err != nil {
				return fmt.Errorf("reading configuration file: %w", err)
			}

			if err := json.Unmarshal(configData, &configuration); err != nil {
				return fmt.Errorf("parsing configuration file: %w", err)
			}
		}

		templateFilepaths, err := templateFilepaths()
		if err != nil {
			return fmt.Errorf("parsing template filenames: %w", err)
		}

		for _, templateFilepath := range templateFilepaths {
			templateFilename := filepath.Base(templateFilepath)
			templateFilenameWithoutExtension := strings.TrimSuffix(templateFilename, filepath.Ext(templateFilepath))
			if *opts.Template != templateFilenameWithoutExtension {
				continue
			}
			// Let's grab other files.
			otherFiles := []*protogen.File{}
			for _, f := range gen.Files {
				if !f.Generate {
					otherFiles = append(otherFiles, f)
				}
			}

			for _, f := range gen.Files {
				if !f.Generate {
					continue
				}
				generatedFilename := fmt.Sprintf(
					"%s_%s.pb.go", f.GeneratedFilenamePrefix, templateFilenameWithoutExtension,
				)
				generatedFile := gen.NewGeneratedFile(generatedFilename, "")
				scopedExecution := newScopedExecution(generatedFile)
				template := template.Must(
					template.New(templateFilename).Funcs(scopedExecution.FuncMap()).ParseFS(templateFS, templateFilepath),
				)
				input := &Input{
					File:          f,
					Files:         otherFiles,
					GeneratedFile: generatedFile,
					Configuration: configuration,
				}
				if err := template.Execute(generatedFile, input); err != nil {
					return fmt.Errorf("executing template: %w", err)
				}
			}
		}
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return nil
	})
}

func templateFilepaths() ([]string, error) {
	filepaths := []string{}
	err := fs.WalkDir(templateFS, "templates", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}
		filepaths = append(filepaths, path)
		return nil
	})
	return filepaths, err
}
