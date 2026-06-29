package profiles

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const MaxImageBytes = 8 * 1024 * 1024

type R2Config struct {
	Region          string
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	PublicBaseURL    string
}

type R2Uploader struct {
	config R2Config
	client *s3.Client
}

func NewR2UploaderFromEnv() (*R2Uploader, error) {
	cfg := R2Config{
		Region:          env("R2_REGION", "auto"),
		Endpoint:        os.Getenv("R2_ENDPOINT"),
		Bucket:          os.Getenv("R2_BUCKET"),
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		PublicBaseURL:    strings.TrimRight(os.Getenv("R2_PUBLIC_BASE_URL"), "/"),
	}

	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.PublicBaseURL == "" {
		return nil, errors.New("cloudflare R2 is not configured")
	}

	client := s3.New(s3.Options{
		Region: cfg.Region,
		BaseEndpoint: aws.String(cfg.Endpoint),
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		),
		UsePathStyle: true,
	})

	return &R2Uploader{
		config: cfg,
		client: client,
	}, nil
}

func env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func cleanStorageKey(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimLeft(value, "/")
	value = regexp.MustCompile(`/+`).ReplaceAllString(value, "/")
	value = regexp.MustCompile(`[^a-zA-Z0-9/_\-./]`).ReplaceAllString(value, "-")
	return value
}

func parseBase64Image(dataURI string) ([]byte, string, string, error) {
	if !strings.HasPrefix(dataURI, "data:image/") {
		return nil, "", "", errors.New("please upload a valid image")
	}

	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", "", errors.New("invalid image data")
	}

	meta := parts[0]
	body := parts[1]

	re := regexp.MustCompile(`^data:(image/(jpeg|jpg|png|webp));base64$`)
	matches := re.FindStringSubmatch(meta)

	if len(matches) < 3 {
		return nil, "", "", errors.New("only JPG, PNG, or WEBP images are allowed")
	}

	raw, err := base64.StdEncoding.DecodeString(body)
	if err != nil || len(raw) == 0 {
		return nil, "", "", errors.New("invalid image")
	}

	if len(raw) > MaxImageBytes {
		return nil, "", "", errors.New("image must be 8MB or smaller")
	}

	contentType := strings.ToLower(matches[1])
	if contentType == "image/jpg" {
		contentType = "image/jpeg"
	}

	exts, _ := mime.ExtensionsByType(contentType)
	ext := "jpg"

	if contentType == "image/png" {
		ext = "png"
	} else if contentType == "image/webp" {
		ext = "webp"
	} else if len(exts) > 0 {
		ext = strings.TrimPrefix(exts[0], ".")
	}

	return raw, contentType, ext, nil
}

func (u *R2Uploader) UploadImage(ctx context.Context, base64Image string, folder string, fileName string) (*Media, error) {
	raw, contentType, ext, err := parseBase64Image(base64Image)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	key := cleanStorageKey(fmt.Sprintf("%s/%s-%d.%s", folder, fileName, now.UnixMilli(), ext))

	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(u.config.Bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(raw),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
		Metadata: map[string]string{
			"provider":   "blomiva",
			"uploadedAt": fmt.Sprintf("%d", now.UnixMilli()),
		},
	})

	if err != nil {
		return nil, err
	}

	publicURL := u.config.PublicBaseURL + "/" + key

	return &Media{
		Provider:    "r2",
		Bucket:      u.config.Bucket,
		Key:         key,
		StorageKey:  key,
		PublicID:    key,
		URL:         publicURL,
		PublicURL:   publicURL,
		ContentType: contentType,
		Bytes:       int64(len(raw)),
		UploadedAt:  &now,
	}, nil
}

func (u *R2Uploader) DeleteImage(ctx context.Context, media *Media) error {
	if media == nil {
		return nil
	}

	key := cleanStorageKey(media.StorageKey)
	if key == "" {
		key = cleanStorageKey(media.Key)
	}

	if key == "" {
		return nil
	}

	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.config.Bucket),
		Key:    aws.String(key),
	})

	return err
}