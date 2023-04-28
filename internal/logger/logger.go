package logger

import (
	"custodyEthereum/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"
)

/**
logger.Debug()
logger.Error()
logger.Info()
**/

// file names, modify if needed more logs
const (
	User        string = "log_user.log"
	Match       string = "log_match.log"
	Deck        string = "log_deck.log"
	Nft         string = "log_nft.log"
	MatchMaking string = "log_matchMak.log"
	Persist     string = "log_persist.log"
	Request     string = "log_request.log"
	DBLog       string = "log_db.log"
)

func NewLogger(env string) (*zap.Logger, error) {
	logDst := "v1/internal/logs/" + env

	debug := configs.GlobalViper.GetBool("server.debug")

	if debug {
		envSplit := strings.Split(env, ".")
		if len(envSplit) == 2 {
			env = envSplit[0] + "_test" + envSplit[1]
		}
	}

	// configs
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	// fileEncoder := zapcore.NewJSONEncoder(config) // Uncomment for JSON encoding, and add it to errorLog and infoLog level
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// create files
	logFile, err := os.OpenFile(logDst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Panicf("logFile creation error %s", err)
	}
	writer := zapcore.AddSync(logFile)

	// defining level
	defaultLogLevel := zapcore.DebugLevel // can be changed if needed
	errorLogLevel := zapcore.ErrorLevel
	infoLogLevel := zapcore.InfoLevel

	// define core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, writer, errorLogLevel), // only errors to log file
		zapcore.NewCore(consoleEncoder, writer, infoLogLevel),  // other info to log file
		zapcore.NewCore(consoleEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel), // development env on console
	)

	// instantiate logger
	logger := zap.New(core, zap.AddCaller())

	return logger, nil
}

/* TASKS FOR IMPLEMENTATION: TBD

- Start with Login/Register DONE
	- Login DONE
	- Register DONE
	- Middleware DONE

- Start with deckGame
	- deckGame
	- tests

- Start with Match
	- matchGame
	- tests

*/
