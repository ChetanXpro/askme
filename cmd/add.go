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

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr                  = flag.String("addr", "localhost:6334", "the address to connect to")
	collectionName        = "test_collection_2"
	vectorSize     uint64 = 4
	distance              = pb.Distance_Dot
)

// addCmd represents the add command
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
