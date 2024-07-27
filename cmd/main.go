package main

import (
	"go-cw-live/internal/adapter/config"
	"go-cw-live/internal/adapter/logs"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {

	profiles := config.GetProfiles()

	selectedProfiles := promptui.Select{
		Label: "Select profile",
		Items: profiles,
	}

	_, selected, err := selectedProfiles.Run()

	if err != nil {
		panic("Prompt failed, " + err.Error())
	}

	cfg := config.LoadConfig(selected)

	promptFilter := promptui.Prompt{
		Label: "Scrivi il prefisso del gruppo di log da cercare",
	}

	filter, err := promptFilter.Run()

	if err != nil {
		panic("Prompt failed, " + err.Error())
	}

	logGroups := logs.GetLogGroups(cfg, filter)

	if len(logGroups) == 0 {
		panic("No log groups found")

	}

	promptLogGroup := promptui.Select{
		Label: "Seleziona il gruppo di log",
		Items: logGroups,
		Searcher: func(input string, index int) bool {
			logGroup := logGroups[index]
			name := strings.Replace(strings.ToLower(logGroup), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		},
	}

	_, selectedLogGroup, err := promptLogGroup.Run()

	if err != nil {
		panic("Prompt failed, " + err.Error())
	}

	logs.GetStreamEventLive(cfg, selectedLogGroup)
}
