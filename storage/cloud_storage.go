package cloud_storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ngtrvu/zen-go/log"
)

type CloudStorageInterface interface {
	LoadJSON(bucketName string, filePath string, v any) error
	LoadJSONL(bucketName string, filePath string) ([][]byte, error)
	CreateFile(ctx context.Context, bucketName string, uploadPath string, file []byte, fileName string) (err error)
	ReadFile(ctx context.Context, bucketName string, uploadPath string, fileName string) (data []byte, err error)
	SignedURL(ctx context.Context, bucketName string, uploadPath string, fileName string, ttl int) (url string, err error)
}

type CloudStorage struct {
	client *storage.Client
}

func NewCloudStorage() *CloudStorage {
	cs := new(CloudStorage)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("storage.NewClient: %v", err)
	}
	defer client.Close()

	_, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	cs.client = client

	return cs
}

func NewCloudStorageWithClient(client *storage.Client) *CloudStorage {
	cs := new(CloudStorage)
	cs.client = client

	return cs
}

// LoadJSON load json file from cloud storage
//
//	bucketName: GCS bucket name
//	filePath: file path in bucket
//	v: Output
func (cs *CloudStorage) LoadJSON(bucketName string, filePath string, v any) error {
	bucket := cs.client.Bucket(bucketName)

	ctx := context.Background()
	rc, err := bucket.Object(filePath).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", filePath, err)
	}
	defer rc.Close()

	dataJson, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}
	return json.Unmarshal(dataJson, &v)
}

// LoadJSONL loads newline-delimited JSON file from cloud storage
//
//	bucketName: GCS bucket name
//	filePath: file path in bucket
//	v: Output
func (cs *CloudStorage) LoadJSONL(bucketName string, filePath string) ([][]byte, error) {
	bucket := cs.client.Bucket(bucketName)

	ctx := context.Background()
	rc, err := bucket.Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", filePath, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var v [][]byte
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		v = append(v, line)
	}

	return v, err
}

// CreateFile Upload file to cloud storage
//
//	bucketName: GCS bucket name
//	uploadPath: file path in bucket
//	file: file content
//	fileName: file name
func (cs *CloudStorage) CreateFile(ctx context.Context, bucketName string, uploadPath string, file []byte, fileName string) (err error) {
	reader := bytes.NewReader(file)
	writer := cs.client.Bucket(bucketName).Object(filepath.Join(uploadPath, fileName)).NewWriter(ctx)
	if _, err = io.Copy(writer, reader); err != nil {
		return
	}

	if err = writer.Close(); err != nil {
		return
	}
	return
}

// ReadFile Read file to cloud storage
//
//	bucketName: GCS bucket name
//	uploadPath: file path in bucket
//	fileName: file name
func (cs *CloudStorage) ReadFile(ctx context.Context, bucketName string, uploadPath string, fileName string) (data []byte, err error) {
	reader, err := cs.client.Bucket(bucketName).Object(uploadPath + fileName).NewReader(ctx)
	if err != nil {
		return
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// SignedURL generate a shared url
//
//	bucketName: GCS bucket name
//	uploadPath: file path in bucket
//	fileName: file name
//	ttl: time to live in second
func (cs *CloudStorage) SignedURL(ctx context.Context, bucketName string, uploadPath string, fileName string, ttl int) (url string, err error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  http.MethodGet,
		Expires: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	return cs.client.Bucket(bucketName).SignedURL(uploadPath+fileName, opts)
}
