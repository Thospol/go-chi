package config

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// CF -> for use configs model
	CF = &Configs{}
)

// Environment environment
type Environment string

const (
	// LOCAL environment local
	LOCAL Environment = "local"
	// DEV environment develop
	DEV Environment = "dev"
	// PROD environment production
	PROD Environment = "prod"
)

// Configs models
type Configs struct {
	UniversalTranslator *ut.UniversalTranslator
	Validator           *validator.Validate
	App                 struct {
		ProjectID  string `mapstructure:"project_name"`
		WebBaseURL string `mapstructure:"web_base_url"`
		APIBaseURL string `mapstructure:"api_base_url"`
		Release    bool   `mapstructure:"release"`
		Port       int    `mapstructure:"port"`
	} `mapstructure:"app"`
	HTTPServer struct {
		ReadTimeout       time.Duration `mapstructure:"read_timeout"`
		WriteTimeout      time.Duration `mapstructure:"write_timeout"`
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	} `mapstructure:"http_server"`
	SwaggerInfo struct {
		Host        string `mapstructure:"host"`
		Title       string `mapstructure:"title"`
		Version     string `mapstructure:"version"`
		Description string `mapstructure:"description"`
		Enable      bool   `mapstructure:"enable"`
	} `mapstructure:"swagger_info"`
	Mysql struct {
		Host           string `mapstructure:"host"`
		Port           string `mapstructure:"port"`
		Username       string `mapstructure:"username"`
		Password       string `mapstructure:"password"`
		DatabaseName   string `mapstructure:"database_name"`
		DriverName     string `mapstructure:"driver_name"`
		ConnectionPool struct {
			MaxOpenConns int           `mapstructure:"max_open_conns"`
			MaxIdleConns int           `mapstructure:"max_idle_conns"`
			MaxLifeTime  time.Duration `mapstructure:"max_lift_time"`
		} `mapstructure:"connection_pool"`
	} `mapstructure:"mysql"`
	Storage struct {
		BucketName               string `mapstructure:"bucket_name"`
		BaseURL                  string `mapstructure:"base_url"`
		DefaultProfilePictureURL string `mapstructure:"default_profile_picture_url"`
	} `mapstructure:"storage"`
	JWT struct {
		SecretKey  string `mapstructure:"secret_key"`
		ExpireTime struct {
			Day    time.Duration `mapstructure:"day"`
			Hour   time.Duration `mapstructure:"hour"`
			Minute time.Duration `mapstructure:"minute"`
		} `mapstructure:"expire_time"`
	} `mapstructure:"jwt"`
}

// InitConfig init config
func InitConfig(configPath string, environment string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(fmt.Sprintf("config.%s", CF.parseEnvironment(environment)))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := bindingConfig(v, CF); err != nil {
			logrus.Error("binding error:", err)
			return
		}
	})

	return nil
}

// bindingConfig binding config
func bindingConfig(vp *viper.Viper, cf *Configs) error {
	if err := vp.Unmarshal(&cf); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}

	cf.Validator = validator.New()

	if err := cf.Validator.RegisterValidation("maxString", validateString); err != nil {
		logrus.Error("cannot register maxString Validator config error:", err)
		return err
	}

	if err := cf.Validator.RegisterValidation("intPercent", validateIntPercent); err != nil {
		logrus.Error("cannot register intPercent Validator config error:", err)
		return err
	}

	en := en.New()
	cf.UniversalTranslator = ut.New(en, en)
	enTrans, _ := cf.UniversalTranslator.GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(cf.Validator, enTrans); err != nil {
		logrus.Error("cannot add english translator config error:", err)
		return err
	}
	_ = cf.Validator.RegisterTranslation("maxString", enTrans, func(ut ut.Translator) error {
		return ut.Add("maxString", "{0} must have number of characters less than 255", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("maxString", fe.Field())
		return t
	})

	_ = cf.Validator.RegisterTranslation("intPercent", enTrans, func(ut ut.Translator) error {
		return ut.Add("intPercent", "{0} must have value between 0 and 100", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("intPercent", fe.Field())
		return t
	})

	return nil
}

// validateString implements validator.Func for max string by rune
func validateString(fl validator.FieldLevel) bool {
	if lengthOfString := utf8.RuneCountInString(fl.Field().String()); lengthOfString > 255 {
		return false
	}

	return true
}

// validateIntPercent implements validator.Func for int percent 0 to 100
func validateIntPercent(fl validator.FieldLevel) bool {
	if percentInInt := fl.Field().Int(); percentInInt < 0 || percentInInt > 100 {
		return false
	}

	return true
}

func (c Configs) parseEnvironment(environment string) Environment {
	switch environment {
	case "local":
		return LOCAL

	case "dev":
		return DEV

	case "prod":
		return PROD
	}

	return DEV
}
