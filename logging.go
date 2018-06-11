package apitool

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const redisLogExpiration = time.Hour * 12

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "api_redis:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var Logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: &logrus.TextFormatter{ForceColors: true},
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

type dbHook struct {
	*redis.Client
}

func init() {

	// Check connection to redis
	if _, err := redisClient.Ping().Result(); err != nil {
		logrus.Error(err)
		return
	}

	// Send logs to redis
	Logger.AddHook(dbHook{redisClient})
}

// Levels add levels that will be shown by hook
func (dbHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func constructKeyFromEntry(e *logrus.Entry) string {

	// Construct fields
	var fields string
	for k, e := range e.Data {
		fields += fmt.Sprintf("{%v:%v}", k, e)
	}

	return fmt.Sprintf("{time:%v}{level:%v}{fields:%v}", e.Time.String(), e.Level.String(), fields)
}

// Fire is called every time something is logged by this logger
func (db dbHook) Fire(e *logrus.Entry) (err error) {

	// Save to redis
	db.Set(constructKeyFromEntry(e), e.Message, redisLogExpiration)

	return
}
