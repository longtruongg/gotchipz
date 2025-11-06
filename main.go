package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"main/cfg"
)

func main() {
	xRPC := []string{
		cfg.ATLANTIC,
		cfg.PharosGotChips,
	}
	//provider, err := ethclient.Dial(cfg.ATLANTIC)
	//if err != nil {
	//	log.Fatalf("Failed to connect to Ethereum client: %v", err)
	//}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}

	param := &cfg.ParamHub{
		Ctx: ctx,
		Key: prik,
	}
	service, err := cfg.NewClientService(param, xRPC)
	if err != nil {
		log.Fatalf("can not create service %s", err)
	}
	param.Provider = service.GetCurrentClient()
	_, err = cfg.SendNativePhrs(param, "0x13f78c79df91419edf23db8cd135b220c1964581")
	if err != nil {
		log.Printf("can not send phrs %s", err)
	}
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
