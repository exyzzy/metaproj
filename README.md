# metaproj (v4.0)

Automatically create a metaapi project, requires metaapi

New for v4:
* Support BYTEA

New for v3:
* Project setup for a fully functioning web server with Vue/Vuetify (-type=vue)

v3.0 matches Medium article:
* Automatic Applications in Go (in progress)

## Previous versions:

Pull the v2.0 tag metaproj commit to have code that matches the Medium article:
* https://levelup.gitconnected.com/automatic-testing-in-go-ce581238eb57

See also:
* http://github.com/exyzzy/metaapi
* http://github.com/exyzzy/metasplice
* http://github.com/exyzzy/pipe

## Most Common Scenario:

```
#assume project and database name: todo (can be anything)
#assume sql file: events.sql (from examples, but can be anything)
createuser -P -d todo <pass: todo>
createdb todo
go get github.com/exyzzy/metaapi
go install $GOPATH/src/github.com/exyzzy/metaapi
go get github.com/exyzzy/metaproj
go install $GOPATH/src/github.com/exyzzy/metaproj
cp $GOPATH/src/github.com/exyzzy/metaapi/examples/events.sql .
metaproj -sql=events.sql -proj=todo -type=vue
cd todo
go generate
go install
go test
```

## Legacy:

First create your PostgreSQL project database:
```
createuser -P -d myproj <pass: myproj>
createdb myproj
```

Create a default internal metaapi project:
```
go get github.com/exyzzy/metaapi
go install $GOPATH/src/github.com/exyzzy/metaapi
go get github.com/exyzzy/metaproj
go install $GOPATH/src/github.com/exyzzy/metaproj
cp $GOPATH/src/github.com/exyzzy/metaapi/examples/alltypes.sql .
# or your own postgreSQL table definition
rm -rf myproj
#clean out the old directory if needed
metaproj -sql=alltypes.sql -proj=myproj 
cd myproj
go generate
go test
```
Create an external metaapi project, where you want to develop custom templates:
```
go get github.com/exyzzy/metaapi
go install $GOPATH/src/github.com/exyzzy/metaapi
go get github.com/exyzzy/metaproj
go install $GOPATH/src/github.com/exyzzy/metaproj
go get github.com/exyzzy/pipe
go install $GOPATH/src/github.com/exyzzy/pipe
cp $GOPATH/src/github.com/exyzzy/metaapi/examples/alltypes.sql .
metaproj -sql=alltypes.sql -proj=myproj -type=external 
# or your own postgreSQL table definition
rm -rf myproj
#clean out the old directory if needed
cd myproj
go install
go generate
go test
```