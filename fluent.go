package logrus_fluent

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fluent/fluent-logger-golang/fluent"
)

const (
	TagField           = "tag"
	MessageField       = "message"
	DefaultTag         = "log"
	PreviousLevelField = "fields.level"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

type fluentHook struct {
	Logger *fluent.Fluent
	levels []logrus.Level
}

func NewHook(host string, port int) (*fluentHook, error) {
	logger, err := fluent.New(fluent.Config{
		FluentHost: host,
		FluentPort: port,
	})
	if err != nil {
		return nil, err
	}
	return &fluentHook{
		Logger: logger,
		levels: defaultLevels,
	}, nil
}

func getTagAndDel(entry *logrus.Entry, data logrus.Fields) string {
	var v interface{}
	var ok bool
	if v, ok = data[TagField]; !ok {
		return entry.Message
	}

	var val string
	if val, ok = v.(string); !ok {
		return DefaultTag
	}
	delete(data, TagField)
	return val
}

func setLevelString(entry *logrus.Entry, data logrus.Fields) {
	data["level"] = entry.Level.String()
}

func setMessage(entry *logrus.Entry, data logrus.Fields) {
	if _, ok := data[MessageField]; !ok {
		data[MessageField] = entry.Message
	}
}

func (hook *fluentHook) Fire(entry *logrus.Entry) error {
	logger, err := fluent.New(fluent.Config{
		FluentHost: hook.host,
		FluentPort: hook.port,
	})
	if err != nil {
		return err
	}
	defer logger.Close()

	// Create a map for passing to FluentD
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}

	setLevelString(entry, data)
	tag := getTagAndDel(entry, data)
	if tag != entry.Message {
		setMessage(entry, data)
	}

	fluentData := ConvertToValue(data, TagName)
	err = logger.PostWithTime(tag, entry.Time, fluentData)
	return err
}

func (hook *fluentHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *fluentHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}
