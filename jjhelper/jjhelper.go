package jjhelper

import (
	"JueJinHelper/task"
	"JueJinHelper/util"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
)

func InitHelper() (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.DisableGPU,
			chromedp.WindowSize(1920, 1080),
		)...,
	)
	ctx, _ = chromedp.NewContext(
		ctx, chromedp.WithLogf(log.Printf),
	)
	return ctx, cancel
}

// Login 登录
func Login(ctx context.Context) error {
	if err := chromedp.Run(ctx, task.Login()); err != nil {
		return err
	}
	return nil
}

// Sign 签到
func Sign(ctx context.Context) error {
	signed := false
	fmt.Println("检查是否签到")
	if err := chromedp.Run(ctx, task.IsSigned(&signed)); err != nil {
		return err
	}
	return nil
}

// PinTask 沸点任务
func PinTask(ctx context.Context) error {
	fmt.Println("发布一个沸点")
	//timestr, err := util.GetNetTime()
	var timeStr, dateStr string
	if err := chromedp.Run(ctx, util.GetBeiJingTime(&timeStr, &dateStr)); err != nil {
		return err
	}
	fmt.Println(dateStr + " " + timeStr)

	if err := chromedp.Run(ctx, task.SetHuangLi()); err != nil {
		return err
	}

	if err := chromedp.Run(ctx, task.SendPin("大家好，JYM, 这是今天的第一个沸点，由机器人【DAO】自动发布 "+dateStr+" "+timeStr+"！", true)); err != nil {
		return err
	}
	return nil
}

// PinDaily 掘金日报
func PinDaily(ctx context.Context) error {

	return nil
}
