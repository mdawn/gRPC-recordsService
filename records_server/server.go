package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mdawn/pindrop-GRPC/recordspb"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
}

type recordsItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Phone   string             `bson:"phone"`
	Carrier string             `bson:"content"`
	Score   string             `bson:"title"`
}

func (*server) CreateRecords(ctx context.Context, req *recordspb.CreateRecordsRequest) (*recordspb.CreateRecordsResponse, error) {
	fmt.Println("Create records request")
	// parse the data
	records := req.GetRecords()

	// map data to a recordsItem
	data := recordsItem{
		Phone:   records.GetPhone(),
		Carrier: records.GetCarrier(),
		Score:   records.GetScore(),
	}

	// pass it to the mongodb driver
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can't convert to OID"),
		)
	}

	return &recordspb.CreateRecordsResponse{
		Records: &recordspb.Records{
			Id:      oid.Hex(),
			Phone:   records.GetPhone(),
			Carrier: records.GetCarrier(),
			Score:   records.GetScore(),
		},
	}, nil
}

func (*server) ReadRecords(ctx context.Context, req *recordspb.ReadRecordsRequest) (*recordspb.ReadRecordsResponse, error) {
	recordsID := req.GetRecordsId()
	oid, err := primitive.ObjectIDFromHex(recordsID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Can't parse ID"),
		)

	}
	// create empty struct
	data := &recordsItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Can't find result with given ID: %v", err),
		)
	}

	return &recordspb.ReadRecordsResponse{
		Records: dataToRecordsPb(data),
	}, nil
}

func dataToRecordsPb(data *recordsItem) *recordspb.Records {
	return &recordspb.Records{
		Id:      data.ID.Hex(),
		Phone:   data.Phone,
		Carrier: data.Carrier,
		Score:   data.Score,
	}
}

func (*server) UpdateRecords(ctx context.Context, req *recordspb.UpdateRecordsRequest) (*recordspb.UpdateRecordsResponse, error) {
	fmt.Println("Update records request")
	records := req.GetRecords()
	oid, err := primitive.ObjectIDFromHex(records.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create empty struct
	data := &recordsItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find records with given ID: %v", err),
		)
	}

	// update internal struct
	data.Phone = records.GetPhone()
	data.Carrier = records.GetCarrier()
	data.Score = records.GetScore()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can't update object in MongoDB: %v", updateErr),
		)
	}

	return &recordspb.UpdateRecordsResponse{
		Records: dataToRecordsPb(data),
	}, nil

}

func (*server) DeleteRecords(ctx context.Context, req *recordspb.DeleteRecordsRequest) (*recordspb.DeleteRecordsResponse, error) {
	fmt.Println("Delete records request")
	oid, err := primitive.ObjectIDFromHex(req.GetRecordsId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_id": oid}

	res, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete records in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Can't find record in MongoDB: %v", err),
		)
	}

	return &recordspb.DeleteRecordsResponse{RecordsId: req.GetRecordsId()}, nil
}

func (*server) ListRecords(req *recordspb.ListRecordsRequest, stream recordspb.RecordsService_ListRecordsServer) error {
	fmt.Println("List records request")

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &recordsItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)

		}
		stream.Send(&recordspb.ListRecordsResponse{Records: dataToRecordsPb(data)})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

func main() {
	// log to tell us file name & line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	// open mongodb client
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Records Service Started")
	collection = client.Database("mydb").Collection("records")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	recordspb.RegisterRecordsServiceServer(s, &server{})
	reflection.Register(s)

	// begin graceful shutdown
	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stopping")
	s.Stop()
	fmt.Println("Closing listener")
	lis.Close()

	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
