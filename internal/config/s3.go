package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func NewS3Client(v *viper.Viper) *s3.Client {
	cfg := aws.Config{
		Region: v.GetString("s3.region"),
		Credentials: credentials.NewStaticCredentialsProvider(
			v.GetString("s3.access_key"),
			v.GetString("s3.secret_key"),
			"",
		),
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(v.GetString("s3.endpoint"))
		o.UsePathStyle = true

		// Untuk S3 Compatible
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
}

// func NewS3Client(viper *viper.Viper) *s3.Client {
// 	accessKey := viper.GetString("s3.access_key")
// 	secretKey := viper.GetString("s3.secret_key")
// 	region := viper.GetString("s3.region")
// 	// bucket := viper.GetString("s3.bucket")
// 	endpoint := viper.GetString("s3.endpoint")
// 	// publicURL := viper.GetString("s3.public_url")

// 	opts := []func(*awsconfig.LoadOptions) error{
// 		awsconfig.WithRegion(region),
// 		awsconfig.WithCredentialsProvider(
// 			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
// 		),
// 	}

// 	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), opts...)
// 	if err != nil {
// 		return nil
// 	}

// 	clientOpts := []func(*s3.Options){}
// 	clientOpts = append(clientOpts, func(o *s3.Options) {
// 		o.BaseEndpoint = aws.String(endpoint)
// 		o.UsePathStyle = true

// 		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
// 		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
// 	})
// 	if endpoint != "" {
// 		clientOpts = append(clientOpts, func(o *s3.Options) {
// 			o.BaseEndpoint = aws.String(endpoint)
// 			o.UsePathStyle = true // diperlukan untuk MinIO
// 			// Disable strict payload checksum — fixes XAmzContentSHA256Mismatch
// 			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
// 			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
// 		})
// 	}

// 	return s3.NewFromConfig(cfg, clientOpts...)
// }
