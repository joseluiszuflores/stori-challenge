package main

import (
	"context"
	"github.com/golang/glog"
	"github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/joseluiszuflores/stori-challenge/internal/bootstrap"
	"github.com/joseluiszuflores/stori-challenge/internal/config"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/file"
	s32 "github.com/joseluiszuflores/stori-challenge/internal/platform/s3"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//nolint:revive
func handler(ctx context.Context, s3Event events.S3Event) error {
	glog.Info("Starting event")
	if err := config.Init(); err != nil {
		return err
	}
	s3D, err := s32.NewS3Reader(config.Config.AWSRegion)
	if err != nil {
		glog.Error("error in s3 Reader", err)

		return err
	}
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		buff, err := s3D.GetBytes(bucket, key)
		if err != nil {
			glog.Error("error getting the data from the file", key, err)

			return err
		}
		transactions, err := file.NewServiceWithBuff(buff)
		if err != nil {
			log.Print("err ToTransactionsFromBytes:", err)

			continue
		}
		//nolint:contextcheck
		if err := bootstrap.Setup(transactions, internal.ToIntFromFile(key)); err != nil {
			glog.Error("err ToTransactionsFromBytes:", err)

			continue
		}

		glog.Info("Finish event")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
