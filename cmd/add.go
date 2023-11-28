/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"flag"
	"log"

	"github.com/spf13/cobra"

	"fmt"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"rsc.io/pdf"
	"strings"
)

var (
	addr                  = flag.String("addr", "localhost:6334", "the address to connect to")
	collectionName        = "test_collection"
	vectorSize     uint64 = 4
	distance              = pb.Distance_Cosine
)

// addCmd represents the add command

func loadAndParsePDF(path string) (content []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := pdf.NewReader(file, file.Stat().Size())
	if err != nil {
		return nil, err
	}

	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text := page.Content().Text
		content = append(content, text)
	}

	return content, nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		flag.Parse()
		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		collections_client := pb.NewCollectionsClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = collections_client.Delete(ctx, &pb.DeleteCollection{
			CollectionName: collectionName,
		})
		if err != nil {
			log.Fatalf("Could not delete collection: %v", err)
		} else {
			log.Println("Collection", collectionName, "deleted")
		}

		// Create collection
		var defaultSegmentNumber uint64 = 2
		_, err = collections_client.Create(ctx, &pb.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     vectorSize,
					Distance: distance,
				},
			}},
			OptimizersConfig: &pb.OptimizersConfigDiff{
				DefaultSegmentNumber: &defaultSegmentNumber,
			},
		})
		if err != nil {
			log.Fatalf("Could not create collection: %v", err)
		} else {
			log.Println("Collection", collectionName, "created")
		}

		r, err := collections_client.List(ctx, &pb.ListCollectionsRequest{})
		if err != nil {
			log.Fatalf("could not get collections: %v", err)
		}
		log.Printf("List of collections: %s", r.GetCollections())

		// Create points grpc client
		pointsClient := pb.NewPointsClient(conn)

		// Create keyword field index
		fieldIndex1Type := pb.FieldType_FieldTypeKeyword
		fieldIndex1Name := "city"
		_, err = pointsClient.CreateFieldIndex(ctx, &pb.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      fieldIndex1Name,
			FieldType:      &fieldIndex1Type,
		})
		if err != nil {
			log.Fatalf("Could not create field index: %v", err)
		} else {
			log.Println("Field index for field", fieldIndex1Name, "created")
		}

		// Create integer field index
		fieldIndex2Type := pb.FieldType_FieldTypeInteger
		fieldIndex2Name := "count"
		_, err = pointsClient.CreateFieldIndex(ctx, &pb.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      fieldIndex2Name,
			FieldType:      &fieldIndex2Type,
		})
		if err != nil {
			log.Fatalf("Could not create field index: %v", err)
		} else {
			log.Println("Field index for field", fieldIndex2Name, "created")
		}

		waitUpsert := true
		upsertPoints := []*pb.PointStruct{
			{
				// Point Id is number or UUID
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{Uuid: "12"},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.05, 0.61, 0.76, 0.74}}}},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.19, 0.81, 0.75, 0.11}}}},
			},

			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 5},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.24, 0.18, 0.22, 0.44}}}},
				Payload: map[string]*pb.Value{
					"count": {
						Kind: &pb.Value_ListValue{
							ListValue: &pb.ListValue{
								Values: []*pb.Value{
									{
										Kind: &pb.Value_IntegerValue{IntegerValue: 0},
									},
								},
							},
						},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 6},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
				Payload: map[string]*pb.Value{},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{Uuid: "58384991-3295-4e21-b711-fd3b94fa73e3"},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
				Payload: map[string]*pb.Value{},
			},
		}

		_, err = pointsClient.Upsert(ctx, &pb.UpsertPoints{
			CollectionName: collectionName,
			Wait:           &waitUpsert,
			Points:         upsertPoints,
		})
		if err != nil {
			log.Fatalf("Could not upsert points: %v", err)
		} else {
			log.Println("Upsert", len(upsertPoints), "points")
		}

		unfilteredSearchResult, err := pointsClient.Search(ctx, &pb.SearchPoints{
			CollectionName: collectionName,
			Vector:         []float32{0.2, 0.1, 0.9, 0.7},
			Limit:          1,
			// Include all payload and vectors in the search result
			WithVectors: &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: true}},
			WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
		})
		if err != nil {
			log.Fatalf("Could not search points: %v", err)
		} else {
			log.Printf("Found points: %s", unfilteredSearchResult.GetResult())
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
