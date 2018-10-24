package main

import (
	"context"
	"fmt"
	"github.com/carapace/cellar/.e2e/app"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

// setting up viper configs
func init() {
	viper.BindEnv("DURATION")
	viper.SetDefault("DURATION", 1*time.Hour)
}

func main() {
	ctx := context.Background()
	ctx, cf := context.WithTimeout(ctx, viper.GetDuration("DURATION"))
	now := time.Now()
	go func() {
		if time.Now().After(now.Add(viper.GetDuration("DURATION"))) {
			cf()
		}
		runtime.Gosched()
	}()

	err := app.Routine(ctx, viper.GetDuration("DURATION"))
	if err != nil {
		panic(fmt.Sprintf("e2e routine returned error %s", err))
	}
	fmt.Println("test success!")
}
