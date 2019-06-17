package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/net/http2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const url = "https://simplifiednetworks.co:443"

var httpVersion = flag.Int("version", 2, "HTTP version")

func testWebCall() string {
	flag.Parse()
	client := &http.Client{}

	// Use the proper transport in the client
	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			// TLSClientConfig: tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			// TLSClientConfig: tlsConfig,
		}
	}

	// Perform the request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed get: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body: %s", err)
	}
	fmt.Printf(
		"Got response %d: %s %s\n",
		resp.StatusCode, resp.Proto, string(body))

	// return string(body)
	return string(body)
}

func run() {
	// Do a fork
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv("ENV", "PRODUCTION")

	checkIfError(cmd.Run())
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/* get full path to binary */
	/*
		binary, lookErr := exec.LookPath(os.Args[2])
		if lookErr != nil {
			panic(lookErr)
		}

		checkIfError(syscall.Exec(binary, args, env)) */

	checkIfError(cmd.Run())
}

func getSecret(name string, region string) (*secretsmanager.GetSecretValueOutput, error) {

	// This example assumes that you're connecting to ap-southeast-1 region
	// For a full list of endpoints, you can refer to this site -> https://godoc.org/github.com/aws/aws-sdk-go/aws/endpoints
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		/*
			 To address specific error, you can import this package:
				"github.com/aws/aws-sdk-go/aws/awserr"
			and use this example:
		*/
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	return result, nil
}

func checkIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	/*
		resp := TestWebCall()
		fmt.Println("About to print response")
		fmt.Printf(resp) */

	if len(os.Args) < 2 {
		println("Error: Invalid command usage")
		return
	}

	secretKeys, err := getSecret("tt-test-secret", endpoints.ApEast1RegionID)

	if err != nil {
		println("Error: Something bad happened getting the keys")
		return
	}

	if secretKeys == nil {
		println("Error: Something bad happened getting the keys. Keys are empty")
		return
	}

	fmt.Println(*secretKeys.SecretString)

	switch os.Args[1] {
	case "child":
		child()
	default:
		run()
	}
}
