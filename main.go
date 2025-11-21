package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"main/cfg"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron/v3"
)

func main() {

	provider, err := ethclient.Dial(cfg.ATLANTIC)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}

	param := &cfg.ParamHub{
		Ctx:      ctx,
		Key:      prik,
		Provider: provider,
	}
	err=cfg.LoadBearToken(); if err!=nil{
		log.Fatalf("cannot load bear %w",err)
	}
	xList:=[]string{}
	c := cron.New(cron.WithLogger(
		cron.DefaultLogger))
	rand.NewSource(time.Now().UnixNano())
	x := gachaGoldenTime()
	c.AddFunc(x, func() {
		for _, val := range cfg.ReadAddrs() {
			tx, err := cfg.SendNativeToken(param, val)
			if err != nil {
				log.Fatalf("can not send phrs %s", err)
			}
			xList=tx
			time.Sleep(10*time.Second)
			
		}
		for _,val:=range xList{
		tx,err:=cfg.VerifyTransferTask(param.Ctx,val)
		if err!=nil{
			log.Fatalf("cannot verify %w",err)
		}
		fmt.Println(tx)
	}
	})
    
	c.Run()

}

func gachaGoldenTime() string {
	var dummyTime = []string{
		"@every 1m",
		// "@every 12m",
		"@every 5m",
		"@every 3m",
		// "@every 7m",
		// "@every 9m",
		// "@every 30m",
		// "@every 15m",
		// "@every 7m",
	}
	rand.NewSource(time.Now().UnixNano())
	return dummyTime[rand.Intn(len(dummyTime))]
}
