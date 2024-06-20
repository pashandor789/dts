package appconfig

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Config interface {
}

func MustParseAppConfig[T Config]() T {
	configFile := flag.String("config-path", "/opt/app/config/application.conf", "Application config file")

	flag.Parse()

	cfg, err := Bind[T](*configFile)
	if err != nil {
		log.Fatalf("can't unmarshal application config: %v", err)
	}

	return cfg
}

func Bind[T Config](configPath string) (T, error) {
	var cfg T

	f, err := os.Open(configPath)
	if err != nil {
		return cfg, fmt.Errorf("can't bind config: %w", err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return cfg, fmt.Errorf("can't read config from %s: %w", configPath, err)
	}

	if err = unmarshalAndSetup(bytes, &cfg); err != nil {
		return cfg, fmt.Errorf("can't unmarshal config: %w", err)
	}

	return cfg, nil
}

func unmarshalAndSetup(data []byte, cfg interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: parseHook,
		Result:     cfg,
		TagName:    "yaml",
	})
	if err != nil {
		return fmt.Errorf("mapstructure.NewDecoder: %w", err)
	}

	var tmpConfig map[string]interface{}
	if err = yaml.Unmarshal(data, &tmpConfig); err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	if err = decoder.Decode(tmpConfig); err != nil {
		return fmt.Errorf("decoder.Decode: %w", err)
	}

	return nil
}

func parseHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String {
		stringData, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("%+v not string (parseHook)", data)
		}

		switch to {
		case reflect.TypeOf(time.Duration(0)):
			return time.ParseDuration(stringData)
		case reflect.TypeOf(time.Time{}):
			return time.Parse(time.RFC3339, stringData)
		case reflect.TypeOf(""):
			return os.ExpandEnv(stringData), nil
		}
	}
	return data, nil
}
