package logs

import (
	"context"
	"go-cw-live/internal/adapter/config"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetLogGroups(cfg *config.Cfg, filter string) []string {

	client := cloudwatchlogs.NewFromConfig(cfg.AwsConfig)

	var prefixFilter *string

	if strings.TrimSpace(filter) != "" {
		prefixFilter = &filter
	} else {
		prefixFilter = nil
	}

	logGroups := []string{}

	for {

		var nextToken *string

		response, err := client.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: prefixFilter,
			NextToken:          nextToken,
		})

		if err != nil {
			panic("Failed to get log groups, " + err.Error())
		}

		for _, logGroup := range response.LogGroups {
			logGroups = append(logGroups, string(*logGroup.LogGroupArn))
		}

		if response.NextToken == nil {
			break
		}

		nextToken = response.NextToken
	}

	return logGroups
}

func GetStreamEventLive(cfg *config.Cfg, logGroup string) {

	client := cloudwatchlogs.NewFromConfig(cfg.AwsConfig)

	request := &cloudwatchlogs.StartLiveTailInput{
		LogGroupIdentifiers: []string{logGroup},
	}

	response, err := client.StartLiveTail(context.TODO(), request)

	if err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}

	stream := response.GetStream()

	handleEventStream(stream)
}

func handleEventStream(stream *cloudwatchlogs.StartLiveTailEventStream) {

	signalChan := make(chan os.Signal, 2)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func(stream *cloudwatchlogs.StartLiveTailEventStream) {
		<-signalChan
		stream.Close()
		os.Exit(1)
	}(stream)

	eventsChan := stream.Events()
	for {
		event := <-eventsChan
		switch e := event.(type) {
		case *types.StartLiveTailResponseStreamMemberSessionStart:
			log.Println("Received SessionStart event")
		case *types.StartLiveTailResponseStreamMemberSessionUpdate:
			for _, logEvent := range e.Value.SessionResults {
				log.Println(*logEvent.Message)
			}
		default:
			// Handle on-stream exceptions
			if err := stream.Err(); err != nil {
				log.Fatalf("Error occured during streaming: %v", err)
			} else if event == nil {
				log.Println("Stream is Closed")
				return
			} else {
				log.Fatalf("Unknown event type: %T", e)
			}
		}
	}
}
