package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Cfg struct {
	Profile   string
	AwsConfig aws.Config
}

func LoadConfig(profile string) *Cfg {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	return &Cfg{
		Profile:   profile,
		AwsConfig: cfg,
	}
}
