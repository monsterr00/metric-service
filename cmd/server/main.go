package main

import (
	"flag"

	"github.com/monsterr00/metric-service.gittest_client/cmd/server/applayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/httplayer"
	"github.com/monsterr00/metric-service.gittest_client/cmd/server/storelayer"
	"github.com/monsterr00/metric-service.gittest_client/internal/util"
)

func main() {
	flag.Parse()
	util.SetFlags()

	// create store layer
	storeLayer := storelayer.New()

	// create app layer
	appLayer := applayer.New(storeLayer)

	// create http layer
	api := httplayer.New(appLayer)

	api.Engage()
}
