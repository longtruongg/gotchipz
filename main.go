package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"main/cfg"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	
	provider, err := ethclient.Dial(cfg.ArcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}

	param := &cfg.ParamHub{
		Ctx: ctx,
		Key: prik,
		Provider: provider,
	}
	
	e, err := cfg.Counter(param,cfg.ARC_COUNTER_METHODD)
	if err != nil {
		log.Printf("can not send phrs %s", err)
	}
	fmt.Print(e)
	//c := cron.New(cron.WithLogger(
	//	cron.DefaultLogger))
	//rand.NewSource(time.Now().UnixNano())
	//x := gachaGoldenTime()
	//c.AddFunc(x, func() {
	//	for _, val := range cfg.ReadAddrs() {
	//		tx, err := cfg.SendNativePhrs(param, val)
	//		if err != nil {
	//			log.Fatalf("can not send phrs %s", err)
	//		}
	//		fmt.Printf(" goldentime{ %s -> %s \n ", tx, x)
	//	}
	//})
	//c.Run()

}

func gachaGoldenTime() string {
	var dummyTime = []string{
		"@every 1m",
		"@every 12m",
		"@every 5m",
		"@every 3m",
		"@every 7m",
		"@every 9m",
		"@every 30m",
		"@every 15m",
		"@every 7m",
	}
	rand.NewSource(time.Now().UnixNano())
	x := dummyTime[rand.Intn(len(dummyTime))]
	return x
}
