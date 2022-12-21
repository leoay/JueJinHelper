package main

import (
	"JueJinHelper/jjhelper"
	"github.com/robfig/cron"
	"log"
)

func main() {
	ctx, cancel := jjhelper.InitHelper()
	defer cancel()
	if err := jjhelper.Login(ctx); err != nil {
		log.Fatal("登录异常", err)
	}

	//if err := jjhelper.Sign(ctx); err != nil {
	//	log.Fatal("签到异常", err)
	//}

	//if err := jjhelper.PinTask(ctx); err != nil {
	//	log.Fatal("发布沸点异常", err)
	//}

	c := cron.New() // 新建一个定时任务对象
	//0 28 0 * * ?   每天0点28分
	c.AddFunc("0 */1 * * * ?", func() {
		ctx, cancel := jjhelper.InitHelper()
		defer cancel()
		if err := jjhelper.Login(ctx); err != nil {
			log.Fatal("登录异常", err)
		}
		//每天的凌晨00：:28分发布一个签到沸点，带黄历
		if err := jjhelper.PinTask(ctx); err != nil {
			log.Fatal("发布沸点异常", err)
		}
	})
	c.Start()
	select {}
}
