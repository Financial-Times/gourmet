# Config 

Simple go library to load configuration into struct. 

Currenty supports only environment variables as source, but can be easily 
exended.

## Usage 

Define struct that represents your configuration

```go
type appConfig struct {
	CustomConfig string `conf:"CUSTOM_CONFIG"`
	App          struct {
		SystemCode string `conf:"APP_SYSTEM_CODE" required:"true"`
		LogLevel   string `conf:"LOG_LEVEL" default:"INFO"`
	}
	Server struct {
		Enabed          bool `conf:"SERVER_ENABLED" default:"true"`
		Port            int  `conf:"SERVER_PORT" default:"8080"`
		ReadTimeout     int  `conf:"SERVER_READ_TIMEOUT" default:"10"`
		WriteTimeout    int  `conf:"SERVER_WRITE_TIMETOUT" default:"15"`
		IdleTimeout     int  `conf:"SERVER_IDLE_TIMEOUT" default:"20"`
	}
}
```

initialize config loader:

```go
confLoader := config.NewEnvConfigLoader()
// shorter version of 
// confLoader := config.NewLoader(config.EnvConfigProvider{})
// you can implement custom data providers
```

and load the data:

```sh
conf := appConfig{}
err := confLoader.Load(&conf)
```

### Supported data types

Your configuration properties can be only of type `int`, `string` or `bool`


## Credits 

This project is "inspired" by [github.com/codingconcepts/env](https://github.com/codingconcepts/env)
