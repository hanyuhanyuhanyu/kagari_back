package persistence

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3Client(ctx context.Context) (cli *s3.Client, err error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return
	}

	cli = s3.NewFromConfig(cfg)
	return
}

const KagariMarkdownBucket = "kagari-markdown"
