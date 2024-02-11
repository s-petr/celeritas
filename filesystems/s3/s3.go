package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/s-petr/celeritas/filesystems"
)

type S3 struct {
	Key      string
	Secret   string
	Region   string
	Endpoint string
	Bucket   string
}

func (s *S3) getSession() *session.Session {
	c := credentials.NewStaticCredentials(s.Key, s.Secret, "")

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         &s.Endpoint,
		Region:           &s.Region,
		Credentials:      c,
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
	}))

	return sess
}

func (s *S3) Put(fileName, folder string) error {
	sess := s.getSession()

	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	var size = fileInfo.Size()

	buffer := make([]byte, size)

	if _, err = f.Read(buffer); err != nil {
		return err
	}

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", folder, path.Base(fileName))),
		Body:        fileBytes,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(fileType),
		Metadata: map[string]*string{
			"Key": aws.String("MetadataValue"),
		},
	})

	return err
}

func (s *S3) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing

	sess := s.getSession()

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println("AWS Error:", s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println("AWS Error:", aerr.Error())
			}
		} else {
			fmt.Println("Other Error:", err.Error())
		}
		return nil, err
	}

	for _, key := range result.Contents {
		b := float64(*key.Size)
		kb := b / 1024
		mb := kb / 1024
		current := filesystems.Listing{
			Etag:         *key.ETag,
			LastModified: *key.LastModified,
			Key:          *key.Key,
			Size:         mb,
		}
		listing = append(listing, current)
	}

	return listing, nil
}

func (s *S3) Delete(itemsToDelete []string) bool {
	sess := s.getSession()

	svc := s3.New(sess)

	for _, item := range itemsToDelete {
		input := &s3.DeleteObjectsInput{
			Bucket: &s.Bucket,
			Delete: &s3.Delete{
				Objects: []*s3.ObjectIdentifier{
					{Key: aws.String(item)},
				},
				Quiet: aws.Bool(false),
			},
		}

		_, err := svc.DeleteObjects(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					fmt.Println("AWS Error:", s3.ErrCodeNoSuchBucket, aerr.Error())
				default:
					fmt.Println("AWS Error:", aerr.Error())
				}
			} else {
				fmt.Println("Other Error:", err.Error())
			}
			return false
		}
	}
	return true
}

func (s *S3) Get(destination string, items ...string) error {
	sess := s.getSession()

	for _, item := range items {
		err := func() error {
			file, err := os.Create(fmt.Sprintf("%s/%s", destination, item))
			if err != nil {
				return err
			}
			defer file.Close()

			downloader := s3manager.NewDownloader(sess)
			_, err = downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(s.Bucket),
					Key:    aws.String(item),
				})
			return err
		}()
		return err
	}

	return nil
}
