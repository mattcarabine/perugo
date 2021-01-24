package main

import (
	"flag"
	"github.com/mattcarabine/perugo/internal/server"
	"go.uber.org/zap"
	"log"
)

var (
	addr   = flag.String("addr", "localhost:8080", "http service address")
	secret = flag.String("secret", "IAmNotSecure", "signing secret for JWT tokens")
)

// All roll dice and see rolls
// Make bet in order
// Raise - increase number of dice or the actual number
// Liar - previous bet is incorrect, show all dice
// Spot on - previous bet is exactly correct
// All dice, out
// five dice to start

// Let's start with a game where users connect to a game and then roll dice and see who wins with the rolls
// Flow:
// User creates room, given a room ID
// User enters room ID to join
// Turn-based rolling

func initLogging() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	logger.Info("initialized zap logger")
	return logger
}
func main() {
	logger := initLogging()
	zap.ReplaceGlobals(logger)
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("unable to flush logger: %s", err.Error())
		}
	}()

	flag.Parse()
	logger.Fatal("failed to setup server", zap.Error(server.SetupServer(*addr, *secret)))
}
