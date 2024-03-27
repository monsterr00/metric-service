package main

import (
	"flag"

	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/httplayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/agent/storelayer"
)

func main() {
	flag.Parse()

	// create store layer
	storeLayer := storelayer.New()

	// create app layer
	appLayer := applayer.New(storeLayer)

	// create http layer
	api := httplayer.New(appLayer)

	api.Engage()
}
