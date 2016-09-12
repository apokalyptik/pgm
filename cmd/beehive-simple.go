// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"strings"

	"googlemaps.github.io/maps"

	"github.com/apokalyptik/pgm/beehive"
	"github.com/apokalyptik/pgm/debug"
	"github.com/apokalyptik/pgm/ll"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// beehive-simpleCmd represents the beehive-simple command
var beehiveSimpleCmd = &cobra.Command{
	Use:   "simple",
	Short: "A simple beehive",
	Run: func(cmd *cobra.Command, args []string) {
		if gm, err := maps.NewClient(maps.WithAPIKey(viper.GetString("maps"))); err != nil {
			panic(err)
		} else {
			ll.Geo = gm
		}

		start, err := ll.NewCoord(viper.GetString("location"), viper.GetInt("jitter"))
		if err != nil {
			log.Fatal(err)
		}

		var accounts = [][]string{}
		for _, account := range viper.GetStringSlice("accounts") {
			accounts = append(accounts, strings.Split(account, ":"))
		}
		beehive.Feed = new(debug.Feed)
		beehive.Start(start, viper.GetInt("steps"), accounts)
		<-make(chan struct{})
	},
}

func init() {
	beehiveCmd.AddCommand(beehiveSimpleCmd)
}
