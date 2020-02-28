package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/exyzzy/metaproj/data"
)

//This struct passed into all templates for generation
//also use receiver methods such as Project, below
type ProjData struct {
	ProjName string
	SqlFile  string
	ProjType string
}

// Create a metaapi base project from a .sql table definition, example:
//  metaproj -proj=newtodo -sql=todo.sql -type=external  (or, for internal:  metaproj -proj=newtodo -sql=todo.sql )
//  cd newtodo
//  go generate
func main() {
	projPtr := flag.String("proj", "myproj", "project to create")
	sqlPtr := flag.String("sql", "", ".sql input file to parse")
	typePtr := flag.String("type", "internal", "project type to create")

	var p ProjData

	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(" valid usage is:")
		fmt.Println("  metproj -proj=yourproj -sql=yoursql.sql")
		fmt.Println("  metproj -proj=yourproj -sql=yoursql.sql -type=external")
		os.Exit(1)
	}

	p.ProjName = strings.ToLower(*projPtr)
	p.SqlFile = strings.ToLower(*sqlPtr)
	p.ProjType = strings.ToLower(*typePtr)

	if (p.SqlFile != "") && (!strings.HasSuffix(p.SqlFile, ".sql")) {
		log.Panic("Invalid .sql File")
	}

	err := createProj(&p)
	if err != nil {
		log.Panic(err)
	}
}

func createProj(pp *ProjData) error {
	if (pp.ProjName == "") || (pp.SqlFile == "") {
		log.Panic(errors.New("Must have projName and sqlFile"))
	}

	var projTypes = map[string]int{"internal": 0, "external": 1}
	pIndex, ok := projTypes[pp.ProjType]
	if !ok {
		log.Panic(errors.New(fmt.Sprintf("Invald project type, use: %v", keys(projTypes))))
	}

	err := os.MkdirAll(pp.ProjName, os.FileMode(0755))
	if err != nil {
		return err
	}

	type FileList struct {
		Name       string
		IsGenerate bool //if false just copy the file with no template actions applied
	}

	var files = [][]FileList{{{"sqlname.txt", true}, {"configlocaldb.txt", true}},
		{{"sqlname2.txt", true}, {"configlocaldb.txt", true}, {"generate.txt", true}, {"main.txt", true}, {"api.txt", false}, {"api_test.txt", false}}}

	for _, f := range files[pIndex] {
		dat, err := data.Asset(f.Name)
		if err != nil {
			return err
		}
		if f.IsGenerate {
			err = generateFile(dat, pp, getDest(pp, f.Name))
			if err != nil {
				return err
			}
		} else {
			err = writeFile(dat, getDest(pp, f.Name))
		}
	}
	err = copyFile(pp, pp.SqlFile)
	return err
}

func keys(ms map[string]int) []string {
	kys := make([]string, len(ms))
	i := 0
	for k := range ms {
		kys[i] = k
		i++
	}
	return kys
}

func getDest(pp *ProjData, name string) string {
	var dest string
	switch name {
	case "sqlname.txt", "sqlname2.txt":
		dest = prefix(pp.SqlFile) + ".go"
	case "configlocaldb.txt":
		dest = "configlocaldb.json"
	case "generate.txt":
		dest = "generate.go"
	case "main.txt":
		dest = "main.go"
	case "api.txt":
		dest = "api.txt"
	case "api_test.txt":
		dest = "api_test.txt"
	}
	return filepath.Join(pp.ProjName, dest)
}

func prefix(name string) string {
	dot := strings.Index(name, ".")
	if dot > 0 {
		return name[:dot]
	} else {
		return name
	}
}

func generateFile(templatesrc []byte, data interface{}, dest string) error {
	tt := template.Must(template.New("file").Parse(string(templatesrc)))
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	err = tt.Execute(file, data)
	file.Close()
	return err
}

func copyFile(pp *ProjData, namesrc string) error {
	dat, err := ioutil.ReadFile("./" + namesrc)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(pp.ProjName, namesrc), dat, 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(templatesrc []byte, dest string) error {
	err := ioutil.WriteFile(dest, templatesrc, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (pp *ProjData) Package() string {
	proj := os.Getenv("GOPACKAGE")
	if proj == "" {
		proj = "main"
	}
	return proj
}
