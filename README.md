# See & Store & More 

A CRUD gRPC made with protobufs using [MongoDB](https://www.mongodb.com/download-center/community) and [Evans](https://github.com/ktr0731/evans)


## What We Need

We're on a Mac, right? Of course we are. 

- The latest version of [Go](https://golang.org/doc/install#install)
- [MongoDB](https://www.mongodb.com/download-center/community)
- [Evans](https://github.com/ktr0731/evans) 
- This repo

## What it Does

This CRUD RPC will:

Create, Read, Update, Delete, and List BSON with an Id (for tracking) and the following fields:
1. Phone
2. Carrier
3. Score

## Getting Operational

**STEP 1**: We open a terminal and start our MongoDB server. 
- `$ cd mongodb-osx-x86_64-4.0.5`
- Install the go driver & grpc: 
`$ go get github.com/mongodb/mongo-go-driver/mongo`
- Run it & set our database path: 
`$ bin/mongod -dbpath data/db`

**STEP 2**: We open another terminal and clone this repository. Cd into the directory. 
- Install gRPC: `$ go get google.golang.org/grpc`
- Run the server: `go run records_server/server.go`

**STEP 3**: We open our last terminal and homebrew install Evans.
- `$ brew tap ktr0731/evans`
- `$ brew install evans`

Then connect:
- `evans -p 50051 -r`


## Example Scenario
We want to CREATE a record.

Once in Evans (having used `evans -p 50051 -r` above), we show available services:
- `show service`

We view the service we want:
-  `service RecordsService`

We call the RPC:
- `call CreateRecords`

And enter the data:
- records::id (TYPE_STRING) => (just hit enter)
- records::phone (TYPE_STRING) => `867-5309`
- records::carrier (TYPE_STRING) => `Charter`
- records::score (TYPE_STRING) => `720`

Our object will return like so:
```
{
  "records": {
    "id": "5c56af21074ee2ca70c1bec7",
    "phone": "867-5309",
    "carrier": "Charter",
    "score": "720"
  }
}
```


Evans can take it from here! See Evans's [Basic Usage](https://github.com/ktr0731/evans#basic-usage) for more fun. :)
