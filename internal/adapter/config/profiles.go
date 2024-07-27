package config

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

func GetProfiles() []string {

	// Get the list of profiles
	filename := config.DefaultSharedCredentialsFilename()

	f, err := ini.Load(filename)

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	profiles := []string{}

	for _, v := range f.Sections() {
		if len(v.Keys()) != 0 {
			profiles = append(profiles, v.Name())
		}
	}

	return profiles
}
