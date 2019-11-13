// oauth_example.go provides a simple example implementing Strava OAuth
// using the go.strava library.
//
// usage:
//   > go get github.com/strava/go.strava
//   > cd $GOPATH/github.com/strava/go.strava/examples
//   > go run oauth_example.go -id=youappsid -secret=yourappsecret
//
//   Visit http://localhost:8080 in your webbrowser
//
//   Application id and secret can be found at https://www.strava.com/settings/api
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	strava "github.com/yushihui/go.strava/strava"
)

const port = 8080

var authenticator *strava.OAuthAuthenticator

func main() {
	strava.ClientId = 36533
	strava.ClientSecret = "5be099c9d101dc124ce545e4d9d3aad15c5aafa9"

	authenticator = &strava.OAuthAuthenticator{
		CallbackURL:            fmt.Sprintf("http://localhost:%d/exchange_token", port),
		RequestClientGenerator: nil,
	}

	http.HandleFunc("/", indexHandler)

	path, err := authenticator.CallbackPath()
	if err != nil {
		// possibly that the callback url set above is invalid
		fmt.Println(err)
		os.Exit(1)
	}
	http.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))

	// start the server
	fmt.Printf("Visit http://localhost:%d/ to view the demo\n", port)
	fmt.Printf("ctrl-c to exit")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	token := "18feede74b77132b29b5e7dd080f9eba022f3112"
	getActivitySummary(token, w)
	// you should make this a template in your real application
	// fmt.Fprintf(w, `<a href="%s">`, authenticator.AuthorizationURL(strava.Permissions.ReadAll, true))
	// fmt.Fprint(w, `<img src="http://1n4rcn88bk4ziht713dla5ub-wpengine.netdna-ssl.com/wp-content/uploads/2017/12/press-thumbnail-asset-logo-02.jpg" />`)
	// fmt.Fprint(w, `</a>`)
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SUCCESS:\nAt this point you can use this information to create a new user or link the account to one of your existing users\n")
	fmt.Fprintf(w, "State: %s\n\n", auth.State)
	fmt.Fprintf(w, "Access Token: %s\n\n", auth.AccessToken)

	fmt.Fprintf(w, "The Authenticated Athlete (you):\n")
	content, _ := json.MarshalIndent(auth.Athlete, "", " ")
	fmt.Fprint(w, string(content))
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprint(w, err)
	}
}

func getActivityDetail(accessToken string, activity *strava.ActivitySummary) ([]*strava.Split, error) {
	client := strava.NewClient(accessToken)
	activityDetail, err := strava.NewActivitiesService(client).Get(2782580752).IncludeAllEfforts().Do()
	if err != nil {
		return nil, err
	} else {
		activity.Description = activityDetail.Description
		activity.ParseTHW()
		//var splits = make([]strava.Split, 30)
		for _, split := range activityDetail.SplitsStandard {
			split.Temperature = activity.Temperature
			split.Humidity = activity.Humidity
			split.WindSpeed = activity.WindSpeed
			split.StartDate = activity.StartDate
			split.ActivityId = activity.Id
			//splits=append(splits,*split)
		}
		return activityDetail.SplitsStandard, nil
	}

}

func getActivitySummary(accessToken string, w http.ResponseWriter) {
	client := strava.NewClient(accessToken)
	start := time.Now()
	activities, err := strava.NewCurrentAthleteService(client).ListActivities().Page(3).PerPage(200).Do()
	var splits []*strava.Split
	if err != nil {
		fmt.Fprintln(w, "get activity failed")

	} else {
		for _, activity := range activities {
			activity_splits, err := getActivityDetail(accessToken, activity)
			if err == nil {
				splits = append(splits, activity_splits...)
			}
		}
	}

	write2File(activities, fmt.Sprintf("result/%d/activities.json", strava.ClientId))
	write2File(splits, fmt.Sprintf("result/%d/splits.json", strava.ClientId))
	elapsed := time.Since(start)
	log.Printf("get activity took %s", elapsed)
}

func write2File(activities interface{}, fileName string) {
	log.Println("write to file - " + fileName)
	content, _ := json.Marshal(activities)
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		log.Println(err)
	}

}
