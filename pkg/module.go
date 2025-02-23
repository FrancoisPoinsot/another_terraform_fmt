package pkg

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	// "github.com/hashicorp/hcl/v2/hclsyntax"
)

type Module struct {
	Children    []*Module
	Name        string
	Source      string
	PathFromCWD string // may not apply
	Files       map[string]*hcl.File
}

// relative to parent
func (m *Module) IsLocal() bool {
	// TODO definitly not that. Check later
	return strings.Contains(m.Source, "git")
}

// Not super useful.
// Would be easier to print json representation.
func (m *Module) Print() {
	printModule(m, 0)
}

func printModule(m *Module, level int) {
	// just because I am using a fake child for now
	// should remove that later
	if m.Name == "" {
		return
	}

	fmt.Printf("%*s[%s]\n", level*2, "", m.Name)
	for fileName, fileContent := range m.Files {
		base_file_name := path.Base(fileName)
		fmt.Printf("%*s%s\n", (level+1)*2, "", base_file_name)

		forEachFile(fileContent)
	}

	for _, child := range m.Children {
		printModule(child, level+1)
	}
}

func forEachFile(hclFile *hcl.File) {

	// gohcl.DecodeBody(fileContent.Body)
	switch body := hclFile.Body.(type) {
	case *hclsyntax.Body:
		spew.Dump(body)

		// body.Blocks
	default:
		fmt.Println("non HCL format non supported")
		os.Exit(1)
	}

	// attrs, diags := body.JustAttributes()
	// if attrs == nil {
	// 	fmt.Println("---diags---")
	// 	fmt.Println(diags)
	// 	return
	// }

	// fmt.Println("---attrs---")
	// fmt.Println(attrs)
}
