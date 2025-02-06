package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

var testModule = "./test/basic_terraform_config"

func main() {

	// TODO expect terraform init done?

	parser := hclparse.NewParser()
	root := Module{
		Name:        "root",
		PathFromCWD: testModule,
		Files:       make(map[string]*hcl.File),
	}
	err := ParseModule(parser, &root)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// need module tree here?
	// not for formating
	// But I still want my tree!
	// it's probably a waste to trace that tree just for formating
	//
	// parser.Files() list the files

	// TODO make a nice module tree
	// then think about what to do with it
}

// no cycle detection, yolo
// in fact we specifically want to allow reparsing the same .tf file because they might be used from techincally different modules
// parser act as a cache, so reading a file twice should not be a problem
//
// recursively: as in module per module. Not folder per folder
func ParseModuleRecusrively(parser *hclparse.Parser, currentModule *Module) error {
	err := ParseModule(parser, currentModule)
	if err != nil {
		return err
	}
	for _, module := range currentModule.Children {
		// parallelise here?
		err := ParseModuleRecusrively(parser, module)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseModule(parser *hclparse.Parser, currentModule *Module) error {
	// I can ignore all these cases for now, and still have something interesting
	// case: currentModule is git ref
	// case: source is relative path, but current module is git
	// anymore case?

	// case: currentModule is local file
	entries, err := os.ReadDir(currentModule.PathFromCWD)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), ".tf") {
			continue
		}
		// symlinks?

		filename := filepath.Join(currentModule.PathFromCWD, info.Name())
		hclFile, diag := parser.ParseHCLFile(filename)
		if diag.HasErrors() {
			return diag
		}
		currentModule.Files[filename] = hclFile

		modules, err := getModuleChildren(currentModule, hclFile)
		if err != nil {
			return err
		}
		currentModule.Children = append(currentModule.Children, modules...)
	}
	return nil
}

// search for modules blocks and return a list of those
// for a single file
func getModuleChildren(currentModule *Module, hclFile *hcl.File) ([]*Module, error) {
	// TODO

	return []*Module{
		{
			Name:        "",
			Source:      "",
			PathFromCWD: "", // currentModule.FileSystemPath + source + path resolution
			Files:       make(map[string]*hcl.File),
			// Child unknown yet at this point
		},
	}, nil
}

type Module struct {
	Children    []*Module
	Name        string
	Source      string
	PathFromCWD string // may not apply
	Files       map[string]*hcl.File
}

// relative to parent
func (m Module) IsLocal() bool {
	// TODO definitly not that. Check later
	return strings.Contains(m.Source, "git")
}
