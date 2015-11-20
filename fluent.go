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

func getTagAndDel(entry *logrus.Entry) string {
	var v interface{}
	var ok bool
	if v, ok = entry.Data[TagField]; !ok {
		return DefaultTag
	}

	var val string
	if val, ok = v.(string); !ok {
		return DefaultTag
	}
	delete(entry.Data, TagField)
	return val
}

func setLevelString(entry *logrus.Entry) {
	entry.Data["level"] = entry.Level.String()
}

func setMessage(entry *logrus.Entry) {
	entry.Data[MessageField] = entry.Message
}

func (hook *fluentHook) Fire(entry *logrus.Entry) error {
	setLevelString(entry)
	tag := getTagAndDel(entry)
	setMessage(entry)
	delete(entry.Data, PreviousLevelField)

	data := ConvertFields(entry.Data)
	data["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	return hook.Logger.PostWithTime(tag, entry.Time, data)
}

func (hook *fluentHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *fluentHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}
