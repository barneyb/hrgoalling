package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mrjones/oauth"
	"strings"
)

func main() {
	consumerKey := ""
	consumerSecret := ""
	bytes, err := ioutil.ReadFile("keys.txt")
	if err == nil {
		keys := strings.Split(strings.Trim(string(bytes), " \t\r\n"), "\n")
		consumerKey = strings.Trim(keys[0], " \t")
		if len(keys) > 1 {
			consumerSecret = strings.Trim(keys[1], " \t")
		}
	}

	if len(consumerKey) == 0 || len(consumerSecret) == 0 {
		fmt.Println("You must create a 'keys.txt' file with your key and secret as the first two lines.")
		os.Exit(1)
	}

	c := oauth.NewConsumer(
	consumerKey,
	consumerSecret,
	oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.fitbit.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.fitbit.com/oauth/authorize",
		AccessTokenUrl:    "https://api.fitbit.com/oauth/access_token",
	})

	c.Debug(true)

	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.Get(
		"https://api.fitbit.com/1/user/-/activities/steps/date/today/1w.json",
		map[string]string{},
		accessToken)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	fmt.Println("Your past week of steps are: " + string(bits))

}
