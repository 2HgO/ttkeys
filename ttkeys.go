package main

import (
	"encoding/json"
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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/viper"
)

const url = "https://simplifiednetworks.co:443"

var httpVersion = flag.Int("version", 2, "HTTP version")

// testWebCall tests whether a web request can be made using the defined http version
func testWebCall() (bool, string) {
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
		return false, fmt.Sprintf("Failed get: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Sprintf("Failed reading response body: %s", err)
	}
	log.Printf(
		"Got response %d: %s %s\n",
		resp.StatusCode, resp.Proto, string(body))

	return true, ""
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
					return nil, fmt.Errorf("%s %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				case secretsmanager.ErrCodeInvalidParameterException:
					return nil, fmt.Errorf("%s %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				case secretsmanager.ErrCodeInvalidRequestException:
					return nil, fmt.Errorf("%s %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				case secretsmanager.ErrCodeDecryptionFailure:
					return nil, fmt.Errorf("%s %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				case secretsmanager.ErrCodeInternalServiceError:
					return nil, fmt.Errorf("%s %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				default:
					return nil, fmt.Errorf("%s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
	}
	return result, nil
}

func setupViper() {
	viper.SetConfigName("ttkeysconfig")  // name of config file (without extension)
	viper.AddConfigPath("/etc/ttkeys/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.ttkeys") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		log.Panicln(fmt.Errorf("fatal error config file: %s", err))
	}
}

func getValForKeyViper(key string) string {
	if !viper.IsSet(key) {
		log.Panic(fmt.Errorf("fatal error: could not find key '%s' in ttkeysconfig file", key))
	}
	return viper.GetString(key)	
}

func init() {
	// Kill program if command line arguments aren't provided
	if len(os.Args) < 2 {
		log.Fatalln("Error: Invalid command usage")
	}

	// Kill program if no internet connection
	//if ok, msg := testWebCall(); !ok {
		//log.Fatalln(msg)
	//}

	// Set up viper config
	setupViper()
}

func main() {
	// Retrieve aws secrets wrt secretname and region specified
	secretKeys, err := getSecret(getValForKeyViper("secretName"), getValForKeyViper("region"))

	switch {
		// Check if error while getting secrets
		case err != nil: log.Panicln(err)

		// Check if secretKeys returned is empty
		case secretKeys == nil: log.Panicln("Error: Something bad happened getting the keys. Keys are empty")
	}

	jsonData := []byte(*secretKeys.SecretString)

	var data map[string]interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Panicln(err)
	}

	// Set process I/O streams as std streams
	process := exec.Command(os.Args[1], os.Args[2:]...)
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	// Inject environmental variables into process env
	for key, val := range data {
		process.Env = append(process.Env, fmt.Sprintf("%s=%s", key, val))
	}

	// Spawn process with process.Start() and panic it error occurs when starting process
	if err := process.Start(); err != nil {
		log.Panicln(err)
	}
}
