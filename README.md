# See & Store & More 

A CRUD gRPC made with protobufs which uses [MongoDB](https://www.mongodb.com/download-center/community) and [Evans](https://github.com/ktr0731/evans)

**_"THIS. WILL GO DOWN. ON YOUR PERMANENT. RECORD."_**
- Violent Femmes, "Kiss Off" ;)

## What You Need
(I assume you're on a Mac.)

- The latest version of [Go](https://golang.org/doc/install#install)
- [MongoDB](https://www.mongodb.com/download-center/community)
- [Evans](https://github.com/ktr0731/evans) 
- This repo

## What it Does

This CRUD API will:

Create, Read, Update, Delete, and List a BSON with an Id (for tracking) and the following fields:
1. Phone
2. Carrier
3. Score

## Getting Operational

**STEP 1**: Open a terminal and start your MongoDB server. 
- `$ cd mongodb-osx-x86_64-4.0.5`
- install the go driver: 
`$ go get github.com/mongodb/mongo-go-driver/mongo`
- run it & set your database path (if you like): 
`$ bin/mongod -dbpath data/db`

**STEP 2**: Open another terminal and clone this repository. Cd into the directory. 
- run the server: `go run records/records_server/server.go`

**STEP 3**: Open your last terminal and homebrew install Evans.
- `$ brew tap ktr0731/evans`
- `$ brew install evans`

Then connect:
- `evans -p 50051 -r`


## Example Scenario
I want to CREATE a record.

Once in Evans (having used `evans -p 50051 -r` above), I show available services:
- `show service`

I view the service I want:
-  `service RecordsService`

I call the RPC:
- `call CreateRecords`

And enter my data:
- records::id (TYPE_STRING) => (just hit enter)
- records::phone (TYPE_STRING) => `867-5309`
- records::carrier (TYPE_STRING) => `Charter`
- records::score (TYPE_STRING) => `720`

Your object will return like so:
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


I'll let Evans take it from here! See Evans's [Basic Usage](https://github.com/ktr0731/evans#basic-usage) for more fun. :)
