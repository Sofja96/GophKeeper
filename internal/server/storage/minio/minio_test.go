//nolint:errcheck
package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

func TestFile_Integration(t *testing.T) {
	ctx := context.Background()

	container, minioClient := setupMinIO(t, ctx)
	defer container.Terminate(ctx)

	bucketName := "test-bucket"
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	client := &client{
		Client: minioClient,
		Bucket: bucketName,
	}

	t.Run("successful upload", func(t *testing.T) {
		fileName := "test.txt"
		content := []byte("Hello, MinIO!")
		expectedObjectName := "uploads/test.txt"

		url, err := client.UploadFile(ctx, fileName, content)
		assert.NoError(t, err)

		expectedURL := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, expectedObjectName)
		assert.Equal(t, expectedURL, url)

		_, err = minioClient.StatObject(ctx, bucketName, expectedObjectName, minio.StatObjectOptions{})
		assert.NoError(t, err)
	})

	t.Run("successful get", func(t *testing.T) {
		fileName := "test.txt"
		content := []byte("Hello, MinIO!")
		expectedObjectName := "uploads/test.txt"

		_, err := client.UploadFile(ctx, fileName, content)
		assert.NoError(t, err)

		FileUrl := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, expectedObjectName)

		fileContent, err := client.GetFile(ctx, FileUrl)
		log.Println("file_content", fileContent)
		assert.NoError(t, err)

		expectedContent := []byte("Hello, MinIO!")
		assert.Equal(t, expectedContent, fileContent)

		_, err = minioClient.StatObject(ctx, bucketName, expectedObjectName, minio.StatObjectOptions{})
		assert.NoError(t, err)
	})

	t.Run("successful delete", func(t *testing.T) {
		fileName := "test.txt"
		content := []byte("Hello, MinIO!")
		expectedObjectName := "uploads/test.txt"

		_, err := client.UploadFile(ctx, fileName, content)
		assert.NoError(t, err)

		FileUrl := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, expectedObjectName)

		err = client.DeleteFile(ctx, FileUrl)
		assert.NoError(t, err)

		_, err = minioClient.StatObject(ctx, bucketName, expectedObjectName, minio.StatObjectOptions{})
		assert.Error(t, err)
	})

	t.Run("successful update", func(t *testing.T) {
		fileName := "test.txt"
		content := []byte("Hello, MinIO!")
		newContent := []byte("Hello, New MinIO!")
		expectedObjectName := "uploads/new_file.txt"

		_, err := client.UploadFile(ctx, fileName, content)
		assert.NoError(t, err)

		fileUrl, err := client.UpdateFile(ctx, fileName, "new_file.txt", newContent)
		assert.NoError(t, err)

		expectedURL := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, expectedObjectName)
		assert.Equal(t, expectedURL, fileUrl)

		_, err = minioClient.StatObject(ctx, bucketName, expectedObjectName, minio.StatObjectOptions{})
		assert.NoError(t, err)
	})

	t.Run("successful update with exists filename", func(t *testing.T) {
		fileName := "test.txt"
		content := []byte("Hello, MinIO!")
		newContent := []byte("Hello, New MinIO!")
		expectedObjectName := "uploads/test.txt"

		_, err := client.UploadFile(ctx, fileName, content)
		assert.NoError(t, err)

		fileUrl, err := client.UpdateFile(ctx, fileName, "test.txt", newContent)
		assert.NoError(t, err)

		expectedURL := fmt.Sprintf("%s/%s/%s", minioClient.EndpointURL(), bucketName, expectedObjectName)
		assert.Equal(t, expectedURL, fileUrl)

		_, err = minioClient.StatObject(ctx, bucketName, expectedObjectName, minio.StatObjectOptions{})
		assert.NoError(t, err)
	})
}

func TestNewMinioClient_Success(t *testing.T) {
	ctx := context.Background()

	container, minioClient := setupMinIO(t, ctx)
	defer container.Terminate(ctx)

	bucketName := "existing-bucket"
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	assert.NoError(t, err)

	testSettings := &settings.Settings{
		MinioEndpoint:   getMinioEndpoint(t, ctx, container),
		MinioUser:       "minio",
		MinioPassword:   "minio123",
		MinioUseSsl:     false,
		MinioBucketName: "bucket",
	}

	client, err := NewMinioClient(testSettings)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewMinioClient_CreateBucket(t *testing.T) {
	ctx := context.Background()
	container, _ := setupMinIO(t, ctx)
	defer container.Terminate(ctx)

	testSettings := &settings.Settings{
		MinioEndpoint:   getMinioEndpoint(t, ctx, container),
		MinioUser:       "minio",
		MinioPassword:   "minio123",
		MinioUseSsl:     false,
		MinioBucketName: "new-bucket",
	}

	logOutput := captureLogs(func() {
		client, err := NewMinioClient(testSettings)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	assert.Contains(t, logOutput, "Создан новый bucket: new-bucket")
}

func TestNewMinioClient_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	container, _ := setupMinIO(t, ctx)
	defer container.Terminate(ctx)

	testSettings := &settings.Settings{
		MinioEndpoint:   getMinioEndpoint(t, ctx, container),
		MinioUser:       "invalid",
		MinioPassword:   "invalid",
		MinioUseSsl:     false,
		MinioBucketName: "bucket",
	}

	client, err := NewMinioClient(testSettings)
	assert.ErrorContains(t, err, "The Access Key Id you provided does not exist in our records")
	assert.Nil(t, client)
}

func setupMinIO(t *testing.T, ctx context.Context) (testcontainers.Container, *minio.Client) {
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		Cmd:          []string{"server", "/data"},
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minio",
			"MINIO_ROOT_PASSWORD": "minio123",
		},
		WaitingFor: wait.ForHTTP("/minio/health/ready").WithPort("9000"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	port, err := container.MappedPort(ctx, "9000")
	assert.NoError(t, err)
	endpoint := fmt.Sprintf("localhost:%s", port.Port())

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "minio123", ""),
		Secure: false,
	})
	assert.NoError(t, err)

	return container, mc
}

func getMinioEndpoint(t *testing.T, ctx context.Context, container testcontainers.Container) string {
	port, err := container.MappedPort(ctx, "9000")
	assert.NoError(t, err)
	return fmt.Sprintf("localhost:%s", port.Port())
}

func captureLogs(f func()) string {
	rescue := log.Writer()
	defer log.SetOutput(rescue)

	r, w, _ := os.Pipe()
	log.SetOutput(w)

	f()
	w.Close()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
