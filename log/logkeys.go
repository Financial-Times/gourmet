package log

const (
	DefaultKeyError         = "error"
	DefaultKeyTime          = "time"
	DefaultKeyCaller        = "caller"
	DefaultKeyMessage       = "message"
)

type KeyNamesConfig struct {
	KeyError   string
	KeyTime    string
	KeyCaller  string
	KeyMessage string
}

// GetFullKeyNameConfig returns KeyNamesConfig that has all key names from the input conf and
// if there are key names missing from the input conf, the default key names are used.
func NewKeyNamesConfig(conf *KeyNamesConfig) *KeyNamesConfig {
	defaultConfig := NewDefaultKeyNamesConfig()

	if conf.KeyError == "" {
		conf.KeyError = defaultConfig.KeyError
	}
	if conf.KeyTime == "" {
		conf.KeyTime = defaultConfig.KeyTime
	}
	if conf.KeyCaller == "" {
		conf.KeyCaller = defaultConfig.KeyCaller
	}
	if conf.KeyMessage == "" {
		conf.KeyMessage = defaultConfig.KeyMessage
	}
	return conf
}

func NewDefaultKeyNamesConfig() *KeyNamesConfig {
	return &KeyNamesConfig{
		KeyError:   DefaultKeyError,
		KeyTime:    DefaultKeyTime,
		KeyCaller:  DefaultKeyCaller,
		KeyMessage: DefaultKeyMessage,
	}
}
