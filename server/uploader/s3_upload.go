package uploader

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	uploader *manager.Uploader
	bucket   string
	writer   *io.PipeWriter
	wg       *sync.WaitGroup
}

func NewS3Uploader(bucket string) *S3Uploader {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Uploader{
		uploader: uploader,
		bucket:   bucket,
		wg:       &sync.WaitGroup{},
	}
}

func (u *S3Uploader) UploadStart(key string) {
	r, w := io.Pipe()
	u.writer = w

	u.wg.Add(1)
	go func() {
		defer u.wg.Done()

		req := s3.PutObjectInput{
			Bucket: &u.bucket,
			Key:    &key,
			Body:   r,
		}
		log.Println("start upload")
		_, err := u.uploader.Upload(context.TODO(), &req)

		log.Println("end upload")
		if err != nil {
			log.Println(err.Error())
		}
	}()
}

func (u *S3Uploader) Upload(input []byte) error {
	_, err := u.writer.Write(input)
	return err
}

func (u *S3Uploader) UploadEnd() {
	u.writer.Close()
	u.wg.Wait()
}
