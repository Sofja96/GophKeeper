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
	envKeyDbSource        = "DB_SOURCE" //todo возможно удалить
	envKeyMinioUser       = "MINIO_ROOT_USER"
	envKeyMinioPassword   = "MINIO_ROOT_PASSWORD"
	envKeyEncryptionKey   = "ENCRYPTION_KEY"
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
	EncryptionKey   string
	PathCert        string
	PathKey         string
}

func GetSettings() (*Settings, error) {
	setEnvFunc := []func() error{
		setEnv(envKeyDebug, false),
		setEnv(envKeyServerHost, "0.0.0.0"),
		setEnv(envKeyServerPort, "8080"),
		setEnv(envKeyDbDsn, ""),
		setEnv(envKeyDbAutoMigration, true),
		setEnv(envKeyDbSource, ""),
		setEnv(envKeyMinioUser, ""),
		setEnv(envKeyMinioPassword, ""),
		setEnv(envKeyEncryptionKey, ""),
		setEnv(envKeyPathCert, ""),
		setEnv(envKeyPathKey, ""),
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

	settings, err := newSettings()

	return settings, err

}

func newSettings() (*Settings, error) {
	return &Settings{
		Debug:           viper.GetBool(envKeyDebug),
		Host:            viper.GetString(envKeyServerHost),
		Port:            viper.GetString(envKeyServerPort),
		DbDsn:           viper.GetString(envKeyDbDsn),
		DbSource:        viper.GetString(envKeyDbSource),
		DbAutoMigration: viper.GetBool(envKeyDbAutoMigration),
		MinioUser:       viper.GetString(envKeyMinioUser),
		MinioPassword:   viper.GetString(envKeyMinioPassword),
		EncryptionKey:   viper.GetString(envKeyEncryptionKey),
		PathCert:        viper.GetString(envKeyPathCert),
		PathKey:         viper.GetString(envKeyPathKey),
	}, nil
}

func setEnv(key string, defaultValue interface{}) func() error {
	return func() error {
		viper.SetDefault(key, defaultValue)
		return viper.BindEnv(key)
	}
}
