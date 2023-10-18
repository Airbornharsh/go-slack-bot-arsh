package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type UserPagination struct {
	Users []slack.User
}

func main() {
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	// channelId := os.Getenv("SLACK_CHANNEL_ID")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(client, socketmode.OptionDebug(true), socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {

			// fmt.Println(socketClient.Events)

			select {
			case <-ctx.Done():
				log.Println("Context canceled")
				return

			case evt := <-socketClient.Events:
				fmt.Println("Event received")
				switch evt.Type {
				case socketmode.EventTypeConnecting:
					fmt.Println("Connecting to Slack with Socket Mode...")
				case socketmode.EventTypeConnectionError:
					fmt.Println("Connection failed. Retrying later...")
				case socketmode.EventTypeConnected:
					fmt.Println("Connected to Slack with Socket Mode.")
				case socketmode.EventTypeEventsAPI:
					eventsAPI, ok := evt.Data.(slackevents.EventsAPIEvent)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)
						continue
					}
					socketClient.Ack(*evt.Request)
					fmt.Printf("Event received: %+v\n", eventsAPI)
					err := HandleEventMessage(eventsAPI, client)
					if err != nil {
						log.Fatal(err)
					}

				case socketmode.EventTypeInteractive:
					callback, ok := evt.Data.(slack.InteractionCallback)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)
						continue
					}
					fmt.Printf("Interaction received: %+v\n", callback)
					socketClient.Ack(*evt.Request)
					// Handle the interaction
				case socketmode.EventTypeSlashCommand:
					cmd, ok := evt.Data.(slack.SlashCommand)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)
						continue
					}
					fmt.Printf("Slash command received: %+v\n", cmd)
					socketClient.Ack(*evt.Request)
				default:
					fmt.Printf("Unexpected: %v\n", evt.Type)
				}
			// default:
				// Do other stuff
				// fmt.Println("Default", count)

				// Handle the events
			}
		}
	}(ctx, client, socketClient)

	// client.PostMessage(channelId, slack.MsgOptionText("Hello World", false))

	socketClient.Run()

}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {

	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(evnt, client)

			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unsupported event type: %s", event.Type)
	}
	return nil
}

func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {

	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	text := strings.ToLower(event.Text)

	attachment := slack.Attachment{}
	// client.PostMessage(channelId, slack.MsgOptionText("Hello World", false))

	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Color = "#4af030"
	} else if strings.Contains(text, "weather") {
		attachment.Text = fmt.Sprintf("Weather is sunny today. %s", (user.Profile.FirstName + " " + user.Profile.LastName))
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("I am good. How are you %s?", user.Name)
		attachment.Color = "#4af030"
	}
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}
