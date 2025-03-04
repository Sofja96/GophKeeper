package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSettings(t *testing.T) {
	originalDebug := os.Getenv(envKeyDebug)
	originalHost := os.Getenv(envKeyServerHost)
	originalPort := os.Getenv(envKeyServerPort)
	originalDbDsn := os.Getenv(envKeyDbDsn)
	originalDbAutoMigration := os.Getenv(envKeyDbAutoMigration)
	originalMinioUser := os.Getenv(envKeyMinioUser)
	originalMinioPassword := os.Getenv(envKeyMinioPassword)
	originalPathCert := os.Getenv(envKeyPathCert)
	originalPathKey := os.Getenv(envKeyPathKey)

	defer func() {
		os.Setenv(envKeyDebug, originalDebug)
		os.Setenv(envKeyServerHost, originalHost)
		os.Setenv(envKeyServerPort, originalPort)
		os.Setenv(envKeyDbDsn, originalDbDsn)
		os.Setenv(envKeyDbAutoMigration, originalDbAutoMigration)
		os.Setenv(envKeyMinioUser, originalMinioUser)
		os.Setenv(envKeyMinioPassword, originalMinioPassword)
		os.Setenv(envKeyPathCert, originalPathCert)
		os.Setenv(envKeyPathKey, originalPathKey)
	}()

	t.Run("Default values", func(t *testing.T) {
		os.Unsetenv(envKeyDebug)
		os.Unsetenv(envKeyServerHost)
		os.Unsetenv(envKeyServerPort)
		os.Unsetenv(envKeyDbDsn)
		os.Unsetenv(envKeyDbAutoMigration)
		os.Unsetenv(envKeyMinioUser)
		os.Unsetenv(envKeyMinioPassword)
		os.Unsetenv(envKeyPathCert)
		os.Unsetenv(envKeyPathKey)

		settings, err := GetSettings()
		assert.NoError(t, err)
		assert.Equal(t, false, settings.Debug)
		assert.Equal(t, "0.0.0.0", settings.Host)
		assert.Equal(t, "8080", settings.Port)
		assert.Equal(t, "", settings.DbDsn)
		assert.Equal(t, true, settings.DbAutoMigration)
		assert.Equal(t, "", settings.DbSource)
		assert.Equal(t, "", settings.MinioUser)
		assert.Equal(t, "", settings.MinioPassword)
		assert.Equal(t, "", settings.PathCert)
		assert.Equal(t, "", settings.PathKey)
	})

	t.Run("Environment variables", func(t *testing.T) {
		os.Setenv(envKeyDebug, "true")
		os.Setenv(envKeyServerHost, "127.0.0.1")
		os.Setenv(envKeyServerPort, "9090")
		os.Setenv(envKeyDbDsn, "postgres://user:password@localhost:5432/dbname")
		os.Setenv(envKeyDbAutoMigration, "false")
		os.Setenv(envKeyMinioUser, "minio_user")
		os.Setenv(envKeyMinioPassword, "minio_password")
		os.Setenv(envKeyPathCert, "/path/to/cert")
		os.Setenv(envKeyPathKey, "/path/to/key")

		settings, err := GetSettings()
		assert.NoError(t, err)
		assert.Equal(t, true, settings.Debug)
		assert.Equal(t, "127.0.0.1", settings.Host)
		assert.Equal(t, "9090", settings.Port)
		assert.Equal(t, "postgres://user:password@localhost:5432/dbname", settings.DbDsn)
		assert.Equal(t, false, settings.DbAutoMigration)
		assert.Equal(t, "minio_user", settings.MinioUser)
		assert.Equal(t, "minio_password", settings.MinioPassword)
		assert.Equal(t, "/path/to/cert", settings.PathCert)
		assert.Equal(t, "/path/to/key", settings.PathKey)
	})

	t.Run("Read from .env file", func(t *testing.T) {
		os.Unsetenv(envKeyDebug)
		os.Unsetenv(envKeyServerHost)
		os.Unsetenv(envKeyServerPort)
		os.Unsetenv(envKeyDbDsn)
		os.Unsetenv(envKeyDbAutoMigration)
		os.Unsetenv(envKeyMinioUser)
		os.Unsetenv(envKeyMinioPassword)
		os.Unsetenv(envKeyPathCert)
		os.Unsetenv(envKeyPathKey)

		envContent := `
			DEBUG=true
			SERVER_HOST=192.168.1.1
			SERVER_PORT=7070
			DB_DSN=postgres://user:password@localhost:5432/dbname
			DB_AUTO_MIGRATION=false
			MINIO_ROOT_USER=minio_user
			MINIO_ROOT_PASSWORD=minio_password
			CERT_PATH=/path/to/cert
			KEY_PATH=/path/to/key
			`
		err := os.WriteFile(".env", []byte(envContent), 0644)
		assert.NoError(t, err)
		defer os.Remove(".env")

		settings, err := GetSettings()
		assert.NoError(t, err)
		assert.Equal(t, true, settings.Debug)
		assert.Equal(t, "192.168.1.1", settings.Host)
		assert.Equal(t, "7070", settings.Port)
		assert.Equal(t, "postgres://user:password@localhost:5432/dbname", settings.DbDsn)
		assert.Equal(t, false, settings.DbAutoMigration)
		assert.Equal(t, "minio_user", settings.MinioUser)
		assert.Equal(t, "minio_password", settings.MinioPassword)
		assert.Equal(t, "/path/to/cert", settings.PathCert)
		assert.Equal(t, "/path/to/key", settings.PathKey)
	})
}
