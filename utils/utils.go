package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dslipak/pdf"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func ReadPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	// remember close file

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func LoadPDF(path string) ([]schema.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	p := documentloaders.NewPDF(f, finfo.Size())

	docs, err := p.LoadAndSplit(
		context.Background(),
		textsplitter.NewTokenSplitter(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not load pdf: %w", err)
	}

	// fmt.Println(docs)

	return docs, nil
}

type Config struct {
	OpenaiApiKeys       string
	PineconeProjectName string
	PineconeIndexName   string
	PineconeEnvironment string
	PineconeAPIKEY      string
}

func ConstructPineconeUrl(config Config) string {

	url := "https://" + config.PineconeIndexName + "-" + config.PineconeProjectName + ".svc." + config.PineconeEnvironment + ".pinecone.io/describe_index_stats"

	return url
}

type IndexResponse struct {
	NameSpaces map[string]struct {
		VectorCount int `json:"vectorCount"`
	} `json:"namespaces"`
}

func GetIndexStats(config Config) []string {

	url := ConstructPineconeUrl(Config(config))

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		fmt.Println("Start", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", string(config.PineconeAPIKEY))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var response IndexResponse

	errr := json.Unmarshal([]byte(body), &response)

	if errr != nil {
		log.Fatalf("Error unmarshalling response body: %v", err)
	}

	allNameSpace := make([]string, len(response.NameSpaces))

	index := 0
	for namespace := range response.NameSpaces {

		allNameSpace[index] = namespace

		index++

	}

	return allNameSpace

}
