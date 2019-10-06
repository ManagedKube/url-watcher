package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("[INFO] Starting endpoint-runner.")

	/////////////////////////////
	// MongoDB Init
	/////////////////////////////
	//mongoConnectionURL := os.Getenv("MONGO_DB_CONNECTION_URL")
	//if mongoConnectionURL == "" {
	//	fmt.Fprintf(os.Stderr, "MONGO_CONNECTION_URI environment variable must be set.\n")
	//	os.Exit(1)
	//}
	//mongoDBName := os.Getenv("MONGO_DB_NAME")
	//if mongoDBName == "" {
	//	fmt.Fprintf(os.Stderr, "MONGO_DB_NAME environment variable must be set.\n")
	//	os.Exit(1)
	//}
	//mongoCollection := os.Getenv("MONGO_COLLECTION")
	//if mongoDBName == "" {
	//	fmt.Fprintf(os.Stderr, "MONGO_COLLECTION environment variable must be set.\n")
	//	os.Exit(1)
	//}

	//clientMongo, err := mongo.Connect(context.TODO(), mongoConnectionURL)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Println("[INFO] Mongo: Connected to MongoDB!")
	//log.Println("[INFO] Mongo: MONGO_DB_NAME", mongoDBName)
	//log.Println("[INFO] Mongo: MONGO_COLLECTION", mongoCollection)

	//mongoDbClient := clientMongo.Database(mongoDBName)
	//mkSlackDispatch.SetMongoClient(mongoDbClient)
	//
	//mkSlackInteractiveComponents.SetMongoClient(mongoDbClient)

	/////////////////////////////
	// Init env params
	/////////////////////////////
	endpointTestJson := os.Getenv("ENDPOINT_TEST_JSON")
	if endpointTestJson == "" {
		fmt.Fprintf(os.Stderr, "ENDPOINT_TEST_JSON environment variable must be set.\n")
		os.Exit(1)
	}


	log.Println("[INFO] ENDPOINT_TEST_JSON: ", endpointTestJson)



	//mkSlackDispatch.SetGcpProject(gcpProject)
	//mkSlackInteractiveComponents.SetGcpProject(gcpProject)

	/////////////////////////////
	// Slack Events
	/////////////////////////////
	//slackVerificationToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	//if slackVerificationToken == "" {
	//	fmt.Fprintf(os.Stderr, "SLACK_VERIFICATION_TOKEN environment variable must be set.\n")
	//	os.Exit(1)
	//}

	// devops-sf - xoxb-3115514008-521006750918-VgYby6JJBhuQgYsBnSPboABc
	// verification token: 9S7eT2PBGhn2e5ZkPmczulna

	// Setting type to Producer
	//mkSlackDispatch.SetPubSubType(mkSlackDispatch.PubSubTypeProducer)

	//http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
	//	buf := new(bytes.Buffer)
	//	buf.ReadFrom(r.Body)
	//	body := buf.String()
	//	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: slackVerificationToken}))
	//	if e != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//
	//	log.Println("[INFO] Event Received: ")
	//	log.Printf("[INFO] Message: %v\n", body)
	//
	//	if eventsAPIEvent.Type == slackevents.URLVerification {
	//		var r *slackevents.ChallengeResponse
	//		err := json.Unmarshal([]byte(body), &r)
	//		if err != nil {
	//			w.WriteHeader(http.StatusInternalServerError)
	//		}
	//		w.Header().Set("Content-Type", "text")
	//		w.Write([]byte(r.Challenge))
	//	}
	//	if eventsAPIEvent.Type == slackevents.CallbackEvent {
	//
	//		innerEvent := eventsAPIEvent.InnerEvent
	//
	//		switch ev := innerEvent.Data.(type) {
	//		case *slackevents.AppMentionEvent:
	//
	//			log.Printf("[INFO] TeamID: %s | User: %s | Type: %s | Channel: %s | Text: %s\n", eventsAPIEvent.TeamID, ev.User, ev.Type, ev.Channel, ev.Text)
	//
	//			event := mkSlackEvent.Event{
	//				TeamID: eventsAPIEvent.TeamID,
	//				EventType: ev.Type,
	//				User: ev.User,
	//				Text: ev.Text,
	//				TimeStamp: ev.TimeStamp,
	//				Channel: ev.Channel,
	//				EventTimeStamp: ev.EventTimeStamp,
	//			}
	//
	//			var actionResponse mkK8sAction.Response
	//
	//			mkSlackDispatch.Do(event, actionResponse)
	//
	//		}
	//	}
	//})
	//
	//http.HandleFunc("/action-endpoint", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("[INFO] Event Received: /action-endpoint")
	//
	//	mkSlackInteractiveComponents.Dispatch(w, r)
	//
	//})
	//
	//http.HandleFunc("/options-load-endpoint", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("[INFO] Event Received: /options-load-endpoint")
	//})

	/////////////////////////////
	// http server listen
	/////////////////////////////
	log.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)

}