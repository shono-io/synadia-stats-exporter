package cmd

import (
    "fmt"
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "os"
    "synadia-stats-exporter/pkg"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string
var natsURL string
var userJwt string
var userSeed string

var interval time.Duration

var rootCmd = &cobra.Command{
    Use:   "synadia-stats-exporter",
    Short: "Export Synadia NATS stats to prometheus",
    RunE: func(cmd *cobra.Command, args []string) error {
        nc, err := nats.Connect(natsURL, nats.UserJWTAndSeed(userJwt, userSeed))
        if err != nil {
            return err
        }

        js, err := jetstream.New(nc)
        if err != nil {
            return err
        }

        w, err := pkg.NewWorker(nc, js, interval)
        if err != nil {
            return err
        }

        return w.Run()
    },
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.synadia-stats-exporter.yaml)")
    rootCmd.PersistentFlags().StringVarP(&natsURL, "nats", "n", "tls://connect.ngs.global", "the nats server url")
    rootCmd.PersistentFlags().StringVarP(&userJwt, "jwt", "j", "", "the nats user jwt")
    rootCmd.PersistentFlags().StringVarP(&userSeed, "seed", "x", "", "the nats user seed")
    rootCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", 5*time.Second, "the interval at which to update the metrics")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        viper.AddConfigPath(home)
        viper.SetConfigType("yaml")
        viper.SetConfigName(".synadia-stats-exporter")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
