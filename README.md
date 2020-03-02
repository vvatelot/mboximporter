# MboxImporter

A basic tool to import your `.mbox` file into MongoDB for query.

Work in progress. Additional info and examples you can see in [Python version](https://github.com/Rpsl/mongodb-gmail) 

## Build

```
go mod download
go build -o mbox-importer ./
```

## Usage

```
Usage of ./mbox-importer:
  -body
    	Parse and insert body of the emails
  -database string
    	The Database name to use in MongoDB (default "mbox-importer")
  -filename string
    	Name of the filename to import
  -headers
    	Parse and insert all headers of the emails
  -init
    	Drop if exist collection and create fresh
  -mongo string
    	The Mongo URI to connect to MongoDB (default "root:example@127.0.0.1")
```



