package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

const (
	LAMBDA_SERVER_PORT    = "9988"
	LAMBDA_INVOKE_PROC_ID = "Function.Invoke"
)

func main() {

	var eventfile, outputfile string
	flag.StringVar(&eventfile, "e", "", "The AWS event template [request] you want to send to the Lamba (must be .json file).")
	flag.StringVar(&outputfile, "o", "", "Print the [response] sent from the AWS Lambda server (local function running).")
	flag.Parse()
	if len(strings.TrimSpace(eventfile)) <= 0 {
		log.Fatalln("You must specify an event [request] to send to the LAMBDA")
	}

	// Dial up to the function over localhost:9988 (reserved port used for Lambas in the AWS SDK)
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%v", LAMBDA_SERVER_PORT))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Translate the event file to raw bytes to send in the request payload to the Lambda
	openEventFile, err := os.Open(eventfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer openEventFile.Close()
	eventFileBytes, err := io.ReadAll(openEventFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Send the request payload to the Lambda over RPC
	lambdaRequest := &messages.InvokeRequest{Payload: eventFileBytes}
	var lambdaResponse messages.InvokeResponse
	err = client.Call(LAMBDA_INVOKE_PROC_ID, lambdaRequest, &lambdaResponse)
	if err != nil {
		log.Fatalln(err)
	}
	if lambdaResponse.Error != nil {
		log.Fatalln(lambdaResponse.Error)
	}

	decodedResponse, err := decodeResponse(lambdaResponse.Payload)
	if err != nil {
		log.Fatalln(err)
	}
	printResponse(decodedResponse, outputfile)
}

func printResponse(responseString string, outputfile string) {
	if len(strings.TrimSpace(responseString)) <= 0 {
		return
	}

	log.Println("LAMBDA RESPONSE:\n", responseString)
	if len(strings.TrimSpace(outputfile)) != 0 {
		file, err := os.Create(outputfile)
		if err != nil {
			log.Printf("Error occurred when creating the output file provided: [%v]\n", err)
		} else {
			file.WriteString(responseString)
		}
		file.Close()
	}
}

func decodeResponse(input []byte) (string, error) {

	marshalledInput, err := json.Marshal(input)
	if err != nil {
		log.Fatalf("could not marshal the lambda reponse payload [%v]\n", err.Error())
	}

	unquotedString, err := strconv.Unquote(string(marshalledInput))
	if err != nil {
		log.Fatalf("could not remove escaped quotes from the lambda reponse [%v]\n", err.Error())
	}

	decoded, err := base64.StdEncoding.DecodeString(unquotedString)
	if err != nil {
		log.Fatalf("could not base64 decode the lambda response string [%v]\n", err.Error())
	}
	return string(decoded), nil
}
