package util

import (
	"JueJinHelper/util/requester"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

// SubTime 取当前时间到其他时间的差值
func SubTime(ts string, now time.Time) string {
	_, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		fmt.Printf("parse string err:%v\n", err)
		return err.Error()
	}
	tlocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("Parse a string according to the time zone format of Dongba district err:%v\n", err)
		return err.Error()
	}
	// 按照指定的时区解析时间
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ts, tlocal)
	if err != nil {
		fmt.Printf("Resolve the time according to the specified time zone:%v\n", err)
		return err.Error()
	}
	// 计算时间的差值
	reverseTime := t.Sub(now)
	return reverseTime.String()
}

type TimeSt struct {
	T int64 `json:"t"`
}

// 通过浏览器获取北京时间
func GetBeiJingTime(timeStr *string, dateStr *string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if err = chromedp.Navigate("\thttps://cn.bing.com/search?q=%E5%8C%97%E4%BA%AC%E6%97%B6%E9%97%B4").Do(ctx); err != nil {
			return
		}

		if err = chromedp.Text("#digit_time", timeStr, chromedp.ByQuery).Do(ctx); err != nil {
			return
		}

		if err = chromedp.Text("#digit_date", dateStr, chromedp.ByQuery).Do(ctx); err != nil {
			return
		}
		return nil
	}
}

// GetNetTime 获取网络时间
func GetNetTime() (string, error) {
	_, body, _ := requester.Fetch("GET", "https://vv.video.qq.com/checktime?otype=json", nil, map[string]string{})

	fmt.Println(string(body))

	var timeSt TimeSt
	err := json.Unmarshal(body, &timeSt)
	if err != nil {
		return "", err
	}

	time1 := time.UnixMilli(timeSt.T).Format("2006-01-02 15:04:05")
	if err != nil {
		return "", err
	}
	fmt.Println(time1)
	return time1, nil
}
