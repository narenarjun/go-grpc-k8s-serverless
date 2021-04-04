package main

import (
	"blog/blogpb"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogitem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func datatoBlog(data *blogitem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create Blog request")
	blog := req.GetBlog()

	data := blogitem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v\n", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID: %v\n", err),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read Blog Request")

	blogID := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot Parse ID"),
		)
	}

	//  empty interface for the blog
	data := &blogitem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with the given ID : %v\n", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: datatoBlog(data),
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update Blog Request")

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot Parse ID"),
		)
	}
	data := &blogitem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with the given ID : %v\n", err),
		)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, Updateerr := collection.ReplaceOne(context.Background(), filter, data)
	if Updateerr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update Object in MongoDB: %v\n", Updateerr),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: datatoBlog(data),
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {

	fmt.Println("Blog Delete request")

	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot Parse ID"),
		)
	}
	filter := bson.M{"_id": oid}

	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete Object in MongoDB: %v\n", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find Blog in MongoDb: %v\n", err),
		)
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}, nil

}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	fmt.Println("List blog request")
	cur, err := collection.Find(context.Background(), primitive.D{{}})

	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown Internal Error: %v\n", err),
		)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &blogitem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding: %v\n", err),
			)
		}

		stream.Send(&blogpb.ListBlogResponse{
			Blog: datatoBlog(data),
		})
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unkown internal Error: %v\n", err),
		)

	}

	return nil
}

func main() {

	// if our program crashes, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	connUri, ok := os.LookupEnv("MONGO_DB_URI")

	if !ok {
		fmt.Println("MONGO_DB_URI must be provided and set a env variable")
		return
	}

	fmt.Println("Connecting to MongoDB")

	// * connection to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(connUri))

	if err != nil {
		log.Fatalf("Error while connecting to mongodb: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error while acquiring connection to database: %v\n", err)
		return
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service started")

	// ! 50051 is the default port for grpc
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failes to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	go func() {
		fmt.Println("starting Server...")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()

	//  wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block untill signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDb connection")
	client.Disconnect(ctx)
	fmt.Println("END of Program")

}
