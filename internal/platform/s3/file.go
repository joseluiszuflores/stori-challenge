package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Reader struct {
	bucket string
	key    string
	sess   *session.Session
}

func NewS3Reader(region string) (*S3Reader, error) {
	//nolint:exhaustruct
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	return &S3Reader{
		sess: sess,
	}, nil
}

func (s *S3Reader) GetBytes(bucket, item string) ([]byte, error) {
	s3St := s3.New(s.sess)
	out, err := s3St.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(out.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
