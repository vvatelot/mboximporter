# Mbox Importer

Parse your "gmail takeout file" or any another `.mbox` file and indexing mail messages into MongoDB. 

After that you can use some aggregation functions for insights or analytics your inbox

### Prerequisites

1) Go [here](https://www.google.com/settings/takeout/custom/gmail) and download your Gmail mailbox, depending on the amount of emails you have accumulated this might take a while. The downloaded archive is in the [mbox format](http://en.wikipedia.org/wiki/Mbox).
2) Build importer from source or download [release](https://github.com/Rpsl/mboximporter/releases)

### Build

```
$ git clone git@github.com:Rpsl/mboximporter.git && cd ./mboximporter
$ go mod download
$ go build -o mbox-importer ./
```

### Usage

Run MongoDB, you can use docker-compose for starting mongodb and web-view panel: 
```bash
docker-compose up
```

Parse the messages:

```bash
mbox-importer -init -filename ~/path/to/your/mail.mbox
```


Connection to the MongoDB instance:

```
mongo -u root -p example --authenticationDatabase admin

> use mbox-importer
switched to db mbox-importer
```

And exec aggregation functions.
* See [examples](https://github.com/Rpsl/mboximporter/tree/master/examples)  
* See [documentation of MongoDB Aggregation framework](https://docs.mongodb.com/manual/aggregation/)
```
> db.mails.aggregate([
    { $match: { labels: { $in: ['inbox'] } } },
    { $unwind: "$sender" },
    { $group: {_id: "$sender", total: {$sum : 1} } },
    { $sort: {"total": -1 } }
]);
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


### Todo

- [ ] Repair parse body
- [ ] Extract examples (aggregate functions) to the personal classes and execute from cli
- [ ] Add `--report` option for executing the aggregates and generate report files
