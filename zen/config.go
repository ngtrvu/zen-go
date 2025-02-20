package zen

import (
	"github.com/ngtrvu/zen-go/gorm"
	"github.com/ngtrvu/zen-go/httpserver"
)

const (
	MaxUploadSizeInMegabyteDefault = int64(10)
)

type LoggingConfig struct {
	Level    string `config:"LOGGING_LEVEL"`
	Encoding string `config:"LOGGING_ENCODING"`
}

type ZenConfig struct {
	HTTPServerConfig        httpserver.HTTPServerConfig
	LoggingConfig           LoggingConfig
	DBConfig                gorm.DBConfig
	ECDSAPrivateKeyPath     string `config:"ECD_SA_PRIVATE_KEY_PATH"`
	ECDSAPublicKeyPath      string `config:"ECD_SA_PUBLIC_KEY_PATH"`
	MonitoringEnabled       bool   `config:"MONITORING_ENABLED"`
	SecretKey               string `config:"SECRET_KEY"`
	MaxUploadSizeInMegabyte int64  `config:"MAX_UPLOAD_SIZE_IN_MEGABYTE"`
}
