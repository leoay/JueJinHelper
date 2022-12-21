package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

//ctx, cancel := context.WithCancel(context.Background())
//cmd := exec.CommandContext(ctx, "./b")
//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//cmd.Stdout = os.Stdout
//cmd.Start()
//time.Sleep(10 * time.Second)
//fmt.Println("退出程序中...", cmd.Process.Pid)
//cancel()
//cmd.Wait()

func CmdWithCtx(commandName string, params []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, commandName, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: false}
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	cancel()
	err = cmd.Wait()
	return err
}
