package emporia

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	apiRoot              = "https://api.emporiaenergy.com"
	apiCustomerDevices   = "customers/devices"
	flowUsernamePassword = "USER_PASSWORD_AUTH"
	userPoolRegion       = "us-east-2"
)

// Emporia models components needed.
type Emporia struct {
	RootTempDir string
	Timezone    string
	Username    string
	Password    string
	ClientID    string
	UserPoolID  string
	DeviceGID   int
	Circuits    []Circuit `yaml:"circuits"`
}

// Circuit specific configuration.
type Circuit struct {
	Name    string `yaml:"name"`
	Channel int    `yaml:"channel"`
}

func (e *Emporia) getCustomerDevices(token *string) (err error) {
	_, err = e.getRequest(token, apiCustomerDevices)
	if err != nil {
		return err
	}

	return
}

func (e *Emporia) getRequest(token *string, endpoint string) (resp *string, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", apiRoot, endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"authtoken":    []string{*token},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//Convert the body to type string
	sb := string(body)

	//a.Log.Info("response status code: %s", res.Status)
	//a.Log.Info("response body: %s", sb)

	return &sb, nil
}

func (e *Emporia) GetLogin() (token *string, err error) {
	authFilename := "emporia_auth"

	reAuth, err := e.reAuthenticate(authFilename)
	if err != nil {
		return nil, err
	}

	if !*reAuth {
		token, err := e.ReadStringFromFile(authFilename)
		if err != nil {
			return nil, err
		}

		return &token, nil
	}

	log.Printf("Looks like we need to reauth with Emporia")

	conf := &aws.Config{Region: aws.String(userPoolRegion)}
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	client := cognito.New(sess)

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(e.Username),
		"PASSWORD": aws.String(e.Password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(e.ClientID),
	}

	res, err := client.InitiateAuth(authTry)
	if err != nil {
		return nil, err
	}

	//a.Log.Info("Auth result: %s", res.AuthenticationResult)
	//a.Log.Info("ID Token: %s", res.AuthenticationResult.IdToken)
	e.WriteValueToFile(authFilename, *res.AuthenticationResult.IdToken)

	return res.AuthenticationResult.IdToken, nil
}

func (e *Emporia) reAuthenticate(authFilename string) (proceed *bool, err error) {
	result := true
	if e.FileExists(authFilename) {
		fileTime, err := e.GetModifiedTime(authFilename)
		if err != nil {
			return nil, err
		}

		now, err := e.LocalTime()
		if err != nil {
			return nil, err
		}

		// 50mins out of 60mins
		if now.Sub(*fileTime).Seconds() > 3000 {
			result = true
			return &result, nil
		}

		// Still have time remaining
		result = false
		return &result, nil
	}

	// Auth file does not exist. Assume we need to auth.
	result = true
	return &result, nil
}
