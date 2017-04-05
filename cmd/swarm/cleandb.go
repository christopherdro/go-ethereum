// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"gopkg.in/urfave/cli.v1"
)

func cleandb(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 0 {
		utils.Fatalf("Takes no argument")
	}

	hash := storage.MakeHashFunc("SHA3")

	var (
		bzzaccount = ctx.GlobalString(SwarmAccountFlag.Name)
		datadir    = ctx.GlobalString(utils.DataDirFlag.Name)
	)

	bzzdir := fmt.Sprintf("%s/swarm/bzz-%s", path.Clean(datadir), bzzaccount)

	configdata, err := ioutil.ReadFile(bzzdir + "/config.json")
	if err != nil {
		utils.Fatalf("Could not open source config file '%s'", filepath.Join(bzzdir, "/config.json"))
	}

	var sourceconfig api.Config
	err = json.Unmarshal(configdata, &sourceconfig)
	if err != nil {
		utils.Fatalf("Corrupt or invalid source config file '%s'", filepath.Join(bzzdir, "/config.json"))
	}
	log.Trace(fmt.Sprintf("bzzkey %v", sourceconfig.BzzKey))

	basekey := common.HexToHash(sourceconfig.BzzKey[2:])
	dbStore, err := storage.NewDbStore(filepath.Join(bzzdir, "chunks"), hash, 10000000, func(k storage.Key) (ret uint8) { return uint8(storage.Proximity(basekey[:], k[:])) })
	if err != nil {
		utils.Fatalf("Cannot initialise dbstore: %v", err)
	}
	dbStore.Cleanup()
}
