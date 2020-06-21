package main

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
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
		fmt.Println("  metproj -proj=yourproj -sql=yoursql.sql -type=vue")
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

	var projTypes = map[string]int{"internal": 0, "external": 1, "vue": 2}
	pIndex, ok := projTypes[pp.ProjType]
	if !ok {
		log.Panic(errors.New(fmt.Sprintf("Invald project type, use: %v", keys(projTypes))))
	}

	err := os.MkdirAll(pp.ProjName, os.FileMode(0755))
	if err != nil {
		return err
	}

	//Resource is flat directory file name of resource, each must be unique
	//Target is target name with directory structure, from project root
	//CopyFn is the copy function to apply
	type FileList struct {
		Resource string
		Target   string
		CopyFn   func([]byte, interface{}, string) error
	}

	//rename projname to sqlname
	var files = [][]FileList{{{"sqlname.txt", "sqlname.go", generateFileWithSQL}, {"configlocaldb.txt", "configlocaldb.json", generateFile}},
		{{"sqlname2.txt", "sqlname.go", generateFileWithSQL}, {"configlocaldb.txt", "configlocaldb.json", generateFile}, {"generate.txt", "generate.go", generateFile}, {"main.txt", "main.go", generateFile}, {"api.txt", "api.txt", writeFile}, {"api_test.txt", "api_test.txt", writeFile}},
		{{"configlocaldb.txt", "data/configlocaldb.json", generateFile}, {"db_util3.txt", "data/db_util.go", writeFile}, {"sqlname3.txt", "sqlname.go", generateFileWithSQL}, {"home3.txt", "templates/home.html", generateFile}, {"homev3.txt", "templates/home.vue.js", writeFile}, {"layout3.txt", "templates/layout.html", generateFile}, {"toolbar3.txt", "templates/toolbar.public.html", generateFile}, {"configapp3.txt", "configapp.json", generateFile}, {"main_route3.txt", "main_route.go", writeFile}, {"main3.txt", "main.go", generateFile}, {"license3.txt", "LICENSE.MD", writeFile}}}

	for _, f := range files[pIndex] {
		dir := filepath.Dir(f.Target)
		dir = filepath.Join(pp.ProjName, dir)

		//make the subdir if necessary
		err := os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}

		dat, err := data.Asset(f.Resource)
		if err != nil {
			return err
		}
		err = f.CopyFn(dat, pp, f.Target)
		if err != nil {
			return err
		}
	}
	_ = pp.DataPath()
	//integrate into FileList
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

//get all the name before the final "."
func filePrefix(name string) string {
	dot := strings.LastIndex(name, ".")
	if dot > 0 {
		return name[:dot]
	} else {
		return name
	}
}

// //get all the path before the final "/"
// func prefixpath(name string) string {
// 	slash := strings.LastIndex(name, "/")
// 	if slash > 0 {
// 		return name[:slash]
// 	} else {
// 		return ""
// 	}
// }

// //get all the path after the final "/", and before the final "."
// func prefix(name string) string {
// 	slash := strings.LastIndex(name, "/")
// 	if slash < 0 {
// 		slash = 0
// 	}
// 	dot := strings.LastIndex(name, ".")
// 	if dot < 0 {
// 		dot = len(name)
// 	}
// 	return name[slash:dot]
// }

//replace default target with SQL file prefix + .go
func generateFileWithSQL(templatesrc []byte, ifc interface{}, dest string) error {
	dir := filepath.Dir(dest)
	newdest := filepath.Join(dir, filePrefix(ifc.(*ProjData).SqlFile)+".go")
	return generateFile(templatesrc, ifc, newdest)
}

//apply template, change name to target
func generateFile(templatesrc []byte, ifc interface{}, dest string) error {
	tt := template.Must(template.New("file").Delims("[[", "]]").Parse(string(templatesrc)))
	// tt := template.Must(template.New("file").Parse(string(templatesrc)))
	file, err := os.Create(filepath.Join(ifc.(*ProjData).ProjName, dest))
	if err != nil {
		return err
	}
	err = tt.Execute(file, ifc)
	file.Close()
	return err
}

//do not apply template, but still change name to target
func writeFile(templatesrc []byte, ifc interface{}, dest string) error {
	err := ioutil.WriteFile(filepath.Join(ifc.(*ProjData).ProjName, dest), templatesrc, 0644)
	if err != nil {
		return err
	}
	return nil
}

//straight copy file to project root
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

//receiver methods

func (pp *ProjData) Package() string {
	proj := os.Getenv("GOPACKAGE")
	if proj == "" {
		proj = "main"
	}
	return proj
}

func (pp *ProjData) CapProjName() string {
	return strings.Title(strings.ToLower(pp.ProjName))
}

func (pp *ProjData) CapSqlFile() string {
	return strings.Title(strings.ToLower(filePrefix(pp.SqlFile)))
}

func (pp *ProjData) DataPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	gopath = filepath.Join(gopath, "src") + "/"
	wd, _ := os.Getwd()
	projpath := filepath.Join(strings.TrimPrefix(wd, gopath), pp.ProjName, "data")
	return projpath
}
