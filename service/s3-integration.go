package service

// import (
// 	"context"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// )

// // S3 - interface with functions signature
// // to read and write files from S3 bucket
// type S3 interface {
// 	GetFiles() error
// 	WriteFiles() error
// 	getDetails() *S3Details
// }

// // S3Details holds the S3 information
// // struct implments the function in S3 interface
// type S3Details struct {
// 	URL        string
// 	BucketName string
// 	Region     string
// 	PathToFile string
// 	session    *session.Session
// }

// // SetS3Details set all the S3 information
// func SetS3Details(ctx context.Context, url, bucketName, region string) S3 {
// 	return &S3Details{URL: url, BucketName: bucketName, Region: region}
// }

// func ReadFiles(S3) error {

// 	details := S3.getDetails()
// 	// Initialize a session that the SDK will use and load
// 	// credentials from the shared credentials file ~/.aws/credentials.
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(details.Region)},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	details.session = sess

// 	err = details.GetFiles()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *S3Details) getDetails() *S3Details {
// 	return s
// }

// // GetFiles read files from given bucket and
// // returns the file content to handle test scenarios
// func (s *S3Details) GetFiles() error {

// 	s3.New(s.session)

// 	// downloader := s3manager.NewDownloader(s.session)
// 	// downloader.DownloadWithIterator()
// 	// _, err := downloader.Download(file,
// 	// 	&s3.GetObjectInput{
// 	// 		Bucket: aws.String(bucket),
// 	// 	})
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }

// // WriteFiles write result files to given bucket and
// // returns error if someyhing goes wrong
// func (s *S3Details) WriteFiles() error {
// 	return nil
// }
