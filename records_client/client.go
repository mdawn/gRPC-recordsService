package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../recordspb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Records client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := recordspb.NewRecordsServiceClient(cc)

	// create the record
	fmt.Println("Creating records")
	records := &recordspb.Records{
		Phone:   "8084917848",
		Carrier: "Charter",
		Score:   "720",
	}
	createRecordsRes, err := c.CreateRecords(context.Background(), &recordspb.CreateRecordsRequest{Records: records})

	if err != nil {
		log.Fatalf("Unexpected erro: %v", err)
	}
	fmt.Printf("Record has been created: %v", createRecordsRes)
	recordsID := createRecordsRes.GetRecords().GetId()

	// read record
	fmt.Println("Reading records")

	_, err2 := c.ReadRecords(context.Background(), &recordspb.ReadRecordsRequest{RecordsId: "5bdc29e661b75adcac496cf4"})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	readRecordsReq := &recordspb.ReadRecordsRequest{RecordsId: recordsID}
	readRecordsRes, readRecordsErr := c.ReadRecords(context.Background(), readRecordsReq)
	if readRecordsErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readRecordsErr)
	}

	fmt.Printf("Records were read: %v \n", readRecordsRes)

	// update records
	newRecords := &recordspb.Records{
		Id:      recordsID,
		Phone:   "8084917848",
		Carrier: "Carrier updated",
		Score:   "710",
	}
	updateRes, updateErr := c.UpdateRecords(context.Background(), &recordspb.UpdateRecordsRequest{Records: newRecords})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Records were updated: %v\n", updateRes)

	// delete Records
	deleteRes, deleteErr := c.DeleteRecords(context.Background(), &recordspb.DeleteRecordsRequest{RecordsId: recordsID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", updateErr)
	}
	fmt.Printf("Records were deleted: %v \n", deleteRes)

	// list Records

	stream, err := c.ListRecords(context.Background(), &recordspb.ListRecordsRequest{})
	if err != nil {
		log.Fatalf("error while calling ListRecords RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetRecords())
	}
}
