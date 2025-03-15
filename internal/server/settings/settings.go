package settings

import (
	"os"

	"github.com/spf13/viper"
)

const (
	envKeyDebug           = "DEBUG"
	envKeyServerHost      = "SERVER_HOST"
	envKeyServerPort      = "SERVER_PORT"
	envKeyDbDsn           = "DB_DSN"
	envKeyDbAutoMigration = "DB_AUTO_MIGRATION"
	envKeyMinioEndpoint   = "MINIO_ENDPOINT"
	envKeyMinioUser       = "MINIO_ROOT_USER"
	envKeyMinioPassword   = "MINIO_ROOT_PASSWORD"
	envMinioBucketName    = "MINIO_BUCKET_NAME"
	envKeyMinioUseSsl     = "MINIO_USE_SSL"
	envKeyPathCert        = "CERT_PATH"
	envKeyPathKey         = "KEY_PATH"
)

type Settings struct {
	Debug           bool
	Host            string
	Port            string
	DbDsn           string
	DbSource        string
	DbAutoMigration bool
	MinioUser       string
	MinioPassword   string
	PathCert        string
	PathKey         string
	MinioEndpoint   string
	MinioUseSsl     bool
	MinioBucketName string
}

// GetSettings загружает настройки из .env файла и переменных окружения,
// если они указаны, и возвращает структуру настроек.
func GetSettings() (*Settings, error) {
	setEnvFunc := []func() error{
		setEnv(envKeyDebug, false),
		setEnv(envKeyServerHost, "0.0.0.0"),
		setEnv(envKeyServerPort, "8080"),
		setEnv(envKeyDbDsn, ""),
		setEnv(envKeyDbAutoMigration, true),
		setEnv(envKeyMinioUser, ""),
		setEnv(envKeyMinioPassword, ""),
		setEnv(envKeyPathCert, ""),
		setEnv(envKeyPathKey, ""),
		setEnv(envKeyMinioEndpoint, "0.0.0.0:9000"),
		setEnv(envKeyMinioUseSsl, false),
		setEnv(envMinioBucketName, ""),
	}

	for _, f := range setEnvFunc {
		err := f()
		if err != nil {
			return nil, err
		}
	}
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	settings := newSettings()

	return settings, nil
}

func newSettings() *Settings {
	return &Settings{
		Debug:           viper.GetBool(envKeyDebug),
		Host:            viper.GetString(envKeyServerHost),
		Port:            viper.GetString(envKeyServerPort),
		DbDsn:           viper.GetString(envKeyDbDsn),
		DbAutoMigration: viper.GetBool(envKeyDbAutoMigration),
		MinioUser:       viper.GetString(envKeyMinioUser),
		MinioPassword:   viper.GetString(envKeyMinioPassword),
		PathCert:        viper.GetString(envKeyPathCert),
		PathKey:         viper.GetString(envKeyPathKey),
		MinioEndpoint:   viper.GetString(envKeyMinioEndpoint),
		MinioUseSsl:     viper.GetBool(envKeyMinioUseSsl),
		MinioBucketName: viper.GetString(envMinioBucketName),
	}
}

func setEnv(key string, defaultValue interface{}) func() error {
	return func() error {
		viper.SetDefault(key, defaultValue)
		return viper.BindEnv(key)
	}
}
