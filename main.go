package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"

	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// TODO: Support profiles by different .env file names

var requiredInputs = []string{"username", "password", "clientId", "clientSecret", "audience", "issuer"}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {

	// This will store all input arguments
	inputsMap := make(map[string]*string)

	// Load the profile if one was provided
	profileName := flag.String("profile", "", "Profiles are structered like .env files.")
	flag.Parse()
	if len(*profileName) > 0 {
		if err := godotenv.Load(*profileName); err != nil {
			log.Printf("Error loading the .env file: %v", err)
		}
	}

	// Initialize all the required flags
	for i := 0; i < len(requiredInputs); i++ {
		inputsMap[requiredInputs[i]] = flag.String(requiredInputs[i], "", "")
	}
	flag.Parse()

	// Validate input args
	for i := 0; i < len(requiredInputs); i++ {
		if len(*inputsMap[requiredInputs[i]]) > 0 {
			// The value was provided on the command line
			fmt.Printf(fmt.Sprintf("Reading %s argument from %s: %s...\n", color.BlueString(requiredInputs[i]), color.CyanString("command line"), color.GreenString((*inputsMap[requiredInputs[i]])[:3])))
		} else if profileInput := os.Getenv(strings.ToUpper(requiredInputs[i])); len(profileInput) > 0 {
			// The input exists in the .env file
			fmt.Printf(fmt.Sprintf("Reading %s argument from %s: %s...\n", color.BlueString(requiredInputs[i]), color.MagentaString(*profileName), color.GreenString(profileInput[:3])))
			inputsMap[requiredInputs[i]] = &profileInput
		}
	}

	for i := 0; i < len(requiredInputs); i++ {
		if len(*(inputsMap[requiredInputs[i]])) <= 0 {
			// Whatever is not provided by .env or input args, should be provided now
			fmt.Printf("Please enter a %s: ", color.YellowString(requiredInputs[i]))
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			// Remove newline
			text = strings.Replace(text, "\n", "", -1)
			inputsMap[requiredInputs[i]] = &text
			fmt.Printf(fmt.Sprintf("You entered: %s...\n", color.GreenString((text)[:3])))
		}
	}

	// generate a bearer token
	GetBearerToken(inputsMap)

	// TODO: generate an m2m token

}

func GetBearerToken(inputsMap map[string]*string) {

	reqBody := url.Values{}
	reqBody.Set("username", *inputsMap["username"])
	reqBody.Set("password", *inputsMap["password"])
	reqBody.Set("client_id", *inputsMap["clientId"])
	reqBody.Set("client_secret", *inputsMap["clientSecret"])
	reqBody.Set("audience", *inputsMap["audience"])
	reqBody.Set("grant_type", "http://auth0.com/oauth/grant-type/password-realm")
	reqBody.Set("realm", "Username-Password-Authentication")

	resp, err := http.Post(fmt.Sprintf("https://%s%s", *inputsMap["issuer"], "/oauth/token"), "application/x-www-form-urlencoded", strings.NewReader(reqBody.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := OAuthResponse{}

	json.Unmarshal(b, &data)

	fmt.Printf("\nYour JWT is:\n\n%s\n\n", color.YellowString(data.AccessToken))
}
