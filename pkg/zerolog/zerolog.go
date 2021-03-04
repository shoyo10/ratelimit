package zerolog

import (
	"fmt"

	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	// Teal ...
	Teal = Color("\033[1;36m%s\033[0m")
	// Yellow ...
	Yellow = Color("\033[35m%s\033[0m")
)

// Color ...
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type severityHook struct{}

// Config ...
type Config struct {
	Debug bool
	Local bool
	AppID string `yaml:"app_id" mapstructure:"app_id"`
	Env   string
}

// Init ...
func Init(c *Config) {
	zerolog.DisableSampling(true)
	zerolog.TimestampFieldName = "local_timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	hostname, _ := os.Hostname()
	lvl := zerolog.InfoLevel
	if c.Debug {
		lvl = zerolog.DebugLevel
	}

	var z zerolog.Logger

	if c.Local {
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("[ %s ]", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", Teal(i))
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatTimestamp = func(i interface{}) string {
			t := fmt.Sprintf("%s", i)
			millisecond, err := strconv.ParseInt(fmt.Sprintf("%s", i), 10, 64)
			if err == nil {
				t = time.Unix(int64(millisecond/1000), 0).Local().Format("2006/01/02 15:04:05")
			}
			return Yellow(t)
		}
		z = zerolog.New(output)
	} else {
		z = zerolog.New(os.Stdout)
	}

	log.Logger = z.With().
		Fields(map[string]interface{}{
			"app_id": c.AppID,
			"env":    c.Env,
		}).
		Str("host", hostname).
		Timestamp().
		Caller().
		Logger().
		Level(lvl)
}
