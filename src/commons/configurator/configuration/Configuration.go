package configuration

import (
	"strconv"
	"time"

	"github.com/gofiber/swagger"
)

const serviceName = "GoGallery"

var configuration *Configuration

type Configuration struct {
	args                 map[string]string
	serviceName          string
	sesionId             string
	timestamp            time.Time
	port                 string
	jwtSecret            string
	swaggerConfiguration swagger.Config
}

func Instance(args map[string]string) *Configuration {
	if configuration == nil {
		timestamp := time.Now()
		miliseconds := timestamp.UnixMilli()

		return &Configuration{
			serviceName:          serviceName,
			sesionId:             serviceName + "-" + strconv.FormatInt(miliseconds, 10),
			timestamp:            timestamp,
			args:                 args,
			port:                 args["GO_GALLERY_API_PORT"],
			jwtSecret:            args["JWT_SECRET"],
			swaggerConfiguration: createSwaggerConfiguration(),
		}
	}

	panic("ERROR: Configuration is already intanced")
}

func GetInstance() *Configuration {
	if configuration == nil {
		panic("ERROR: Configuration is not instanced")
	}
	return configuration
}

func createSwaggerConfiguration() swagger.Config {
	return swagger.Config{
		URL: "/api/docs/definition/swagger.json",
	}
}

func (conf *Configuration) GetSwaggerConfiguration() swagger.Config {
	return conf.swaggerConfiguration
}

func (conf *Configuration) GetArgs() map[string]string {
	return conf.args
}

func (conf *Configuration) GetArg(key string) string {
	return conf.args[key]
}

func (conf *Configuration) GetServiceName() string {
	return conf.serviceName
}

func (conf *Configuration) GetSessionId() string {
	return conf.sesionId
}

func (conf *Configuration) GetTimestamp() time.Time {
	return conf.timestamp
}

func (conf *Configuration) GetPort() string {
	return conf.port
}

func (conf *Configuration) GetJWTSecret() string {
	return conf.jwtSecret
}
