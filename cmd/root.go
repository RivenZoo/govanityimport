// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/RivenZoo/govanityimport/config"
)

var cfgFile string
var (
	listenAddress *string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github.com/RivenZoo/govanityimport",
	Short: "Golang vanity import server and client.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.govanityimport.yaml)")
	listenAddress = rootCmd.PersistentFlags().StringP("listen", "L", ":8080", "server listen address")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP("debug.trace", "T", false, "debug trace switch")
	rootCmd.PersistentFlags().StringP("debug.grpctraceaddress", "", ":15031", "grpc trace http server listen address")

	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("debug.trace", rootCmd.PersistentFlags().Lookup("debug.trace"))
	viper.BindPFlag("debug.grpctraceaddress", rootCmd.PersistentFlags().Lookup("debug.grpctraceaddress"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".govanityimport" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".govanityimport")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "read config error: ", err)
		os.Exit(-1)
	}
	cfg := config.GetConfig()
	err := viper.Unmarshal(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal config error: ", err)
		os.Exit(-1)
	}
	fmt.Fprintln(os.Stdout, "config: %s", cfg)
}
