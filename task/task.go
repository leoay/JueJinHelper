package task

import (
	"JueJinHelper/util/cmd"
	"context"
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/tuotoo/qrcode"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

func ClosePlugin() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		if err := chromedp.WaitVisible(".ion-close", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		if err := chromedp.Click(".ion-close", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		return nil
	}
}

// SignIn 签到
func SignIn() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		//先点击图像
		if err := chromedp.Click(".avatar-wrapper", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		if err := chromedp.Click(".ore", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		//点击签到按钮
		if err := chromedp.WaitVisible(".code-calender", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := chromedp.Click(".code-calender", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := chromedp.WaitVisible(".btn-area > .btn", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		//去抽奖
		if err := chromedp.Click(".btn-area > .btn", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		fmt.Println("准备抽奖......")
		time.Sleep(2 * time.Second)
		if err := chromedp.Click("#turntable-item-0", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		//time.Sleep(10 * time.Second)
		if err := chromedp.WaitVisible(".submit", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := chromedp.Click(".submit", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		fmt.Println("抽奖完毕，准备退出......")
		time.Sleep(2 * time.Second)
		return nil
	}
}

// ReadPins 浏览沸点，并评论
func ReadPins() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		//跳转到沸点页面
		if err := chromedp.Navigate("https://juejin.cn/pins").Do(ctx); err != nil {
			return err
		}

		if err := chromedp.Click(".action-title-box", chromedp.ByQuery).Do(ctx); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}
}

// SendPin 写文本沸点
func SendPin(msg string, withpic bool) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		if err := chromedp.Navigate("https://juejin.cn/pins").Do(ctx); err != nil {
			return err
		}
		if err := chromedp.WaitVisible(".auth-card > .rich-editor", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := chromedp.SendKeys(".auth-card > .rich-editor", msg, chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}

		if withpic {
			//如果有图片，则从剪贴板粘贴
			fmt.Println("粘贴黄历到沸点")
			if err := chromedp.KeyEvent(kb.Paste).Do(ctx); err != nil {
				return err
			}
		}

		time.Sleep(2 * time.Second)
		if err := chromedp.Click("//button[normalize-space() = '发布']", chromedp.BySearch).Do(ctx); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
		return nil
	}
}

func addNewTabListener(ctx context.Context) <-chan target.ID {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	return chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return info.URL != ""
	})
}

//func printQRCode(code []byte) (err error) {
//	img, _, err := image.Decode(bytes.NewReader(code))
//	if err != nil {
//		return
//	}
//	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
//	if err != nil {
//		return
//	}
//	res, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
//	if err != nil {
//		return
//	}
//	return
//}

// 加载Cookies
func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {

		//先判断文件夹是否存在
		if _, _err := os.Stat("tmp"); os.IsNotExist(_err) {
			fmt.Println("tmp目录不存在，正在创建")
			os.Mkdir("tmp", os.ModePerm)
			return
		}

		// 如果cookies临时文件不存在则直接跳过
		if _, _err := os.Stat("tmp/cookies.tmp"); os.IsNotExist(_err) {
			fmt.Println("cookie文件不存在")
			return
		}
		cookiesData, err := os.ReadFile("tmp/cookies.tmp")
		if err != nil {
			return
		}
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			return
		}
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

// 微信登录
func enterWechatLogin() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click("//button[normalize-space() = '登录']", chromedp.BySearch),
		chromedp.WaitVisible(`form[class="auth-form"]`),
		chromedp.Click(`span[class="clickable"]`),
		chromedp.Click(`img[title="微信"]`),
	}
}

// 获取微信二维码图片
func getCode() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 1. 用于存储图片的字节切片
		var code []byte
		// 2. 截图
		// 注意这里需要注明直接使用ID选择器来获取元素（chromedp.ByID）
		chromedp.WaitVisible(`img[class="web_qrcode_img"]`)
		if err = chromedp.Screenshot(`img[class="web_qrcode_img"]`, &code, chromedp.NodeVisible).Do(ctx); err != nil {
			return
		}

		// 3. 保存文件
		if err = os.WriteFile("tmp/code.png", code, 0755); err != nil {
			return
		}
		//printQRCode(code)

		fi, err := os.Open("tmp/code.png")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer fi.Close()
		qrmatrix, err := qrcode.Decode(fi)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(qrmatrix.Content)

		obj := qrcodeTerminal.New()
		obj.Get(qrmatrix.Content).Print()

		return
	}
}

// 保存Cookies
func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		fmt.Println("开始保持Cookies")
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}
		if err = os.WriteFile("tmp/cookies.tmp", cookiesData, 0755); err != nil {
			return
		}
		fmt.Println("保持Cookies成功")
		return
	}
}

// 等待扫码跳转
func waitScan() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		fmt.Println("等待扫码跳转")
		if err := chromedp.WaitVisible(`img[class="lazy avatar avatar immediate"]`, chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		fmt.Println("登录成功")
		return nil
	}
}

// Login 登录
func Login() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ch := addNewTabListener(ctx)
		if err := chromedp.Run(ctx, loadCookies()); err != nil {
			return err
		}
		if err := chromedp.Navigate("https://juejin.cn").Do(ctx); err != nil {
			return err
		}
		var body string
		//判断
		if err := chromedp.OuterHTML(".main-nav", &body, chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}

		if !strings.Contains(body, "登录") {
			fmt.Println("已登录")
		} else {
			fmt.Println("未登录")

			//未登录，则微信登录
			if err := chromedp.Run(ctx, enterWechatLogin()); err != nil {
				return err
			}

			newCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
			if err := chromedp.Run(newCtx, getCode()); err != nil {
				return err
			}

			if err := chromedp.Run(ctx, waitScan()); err != nil {
				return err
			}

			if err := chromedp.Run(ctx, saveCookies()); err != nil {
				return err
			}
		}
		return nil
	}
}

// IsSigned 是否完成签到
func IsSigned(signed *bool) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		if err := chromedp.WaitVisible(".first-line", chromedp.ByQuery).Do(ctx); err != nil {
			log.Fatal(err)
		}
		body := ""
		if err := chromedp.OuterHTML(".first-line", &body, chromedp.ByQuery).Do(ctx); err != nil {
			log.Fatal(err)
		}
		if strings.Contains(body, "去签到") {
			fmt.Println("未签到")
			*signed = false
		} else {
			fmt.Println("已签到")
			*signed = true
		}
		return nil
	}
}

func SetHuangLi() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		if err := chromedp.Navigate("https://www.huangli.com/huangli/").Do(ctx); err != nil {
			return err
		}
		var code []byte
		if err := chromedp.WaitVisible(".lunar-info", chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := chromedp.Screenshot(".lunar-info", &code, chromedp.ByQuery).Do(ctx); err != nil {
			return err
		}
		if err := os.WriteFile("tmp/tmphl.png", code, 0755); err != nil {
			return err
		}
		go func() {
			cmd.CmdWithCtx("gclip", []string{"-copy", "-f", "tmp/tmphl.png"})
		}()

		return nil
	}
}
