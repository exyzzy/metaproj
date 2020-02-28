# metaproj

Automatically create a metaapi project

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