package main

import (
	"fmt"
	"time"
	"unsafe"

	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ZeronoFreya/go-hotkey"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

const (
	modelName = "popcapgame1.exe"
)

var (
	pid         uint32
	pHandler    windows.Handle
	baseAddress uint64
	start       = make(chan int)
)

func main() {
	go func() {
		var err error
		for ; err != nil || pid == 0; pid, err = findProcessPidByName(modelName) {
			fmt.Printf("[Waiting] 等待游戏启动..., pid: %d, err: %v\n", pid, err)
			time.Sleep(time.Second * 1)
		}
		fmt.Printf("[Game Started] 目标进程 PID: %d\n", pid)
		pHandler, err = getProcessHandle(pid)
		if err != nil {
			fmt.Println("[Error] 获取目标进程句柄失败", err)
			return
		}
		baseAddress, err = getProcessAddress(pHandler)
		if err != nil {
			fmt.Println("[Error] 获取目标进程模块基址失败", err)
			return
		}
		fmt.Printf("[Game Base Address] 0x%X\n", baseAddress)
		start <- 1
	}()
	myApp := app.New()

	fontPath := "NotoSansSC-Regular.ttf"

	// 设置自定义主题
	myApp.Settings().SetTheme(&MyTheme{
		Theme:    theme.DefaultTheme(),
		fontPath: fontPath,
	})
	myWindow := myApp.NewWindow("Plants vs. Zombies")
	myWindow.Resize(fyne.NewSize(300, 230))
	content := widget.NewLabel("[Waiting for Game Start...]")

	UnlimitedSunshineCheck := widget.NewCheck("[1]无限阳光", setUnlimitedSunshineFunc)
	NoCoolingCheck := widget.NewCheck("[2]无CD", setNoCoolingFunc)
	AllZombieComingCheck := widget.NewCheck("[3]僵尸倾巢", setAllZombieComingFunc)
	KillInstantlyCheck := widget.NewCheck("[4]秒杀僵尸", setKillInstantlyFunc)
	if err := hotkey.Register("1", "press", func() {
		UnlimitedSunshineCheck.SetChecked(!UnlimitedSunshineCheck.Checked)
	}); err != nil {
		fmt.Println("Error registering hotkey:", err)
	}
	defer hotkey.Unregister("1", "press")

	if err := hotkey.Register("2", "press", func() {
		NoCoolingCheck.SetChecked(!NoCoolingCheck.Checked)
	}); err != nil {
		fmt.Println("Error registering hotkey:", err)
	}

	if err := hotkey.Register("3", "press", func() {
		AllZombieComingCheck.SetChecked(!AllZombieComingCheck.Checked)
	}); err != nil {
		fmt.Println("Error registering hotkey:", err)
	}
	defer hotkey.Unregister("3", "press")

	if err := hotkey.Register("4", "press", func() {
		KillInstantlyCheck.SetChecked(!KillInstantlyCheck.Checked)
	}); err != nil {
		fmt.Println("Error registering hotkey:", err)
	}
	defer hotkey.Unregister("4", "press")

	go func() {
		<-start
		content.SetText("[Game Running]")
	}()

	myWindow.SetContent(container.NewVBox(content, UnlimitedSunshineCheck, NoCoolingCheck, AllZombieComingCheck, KillInstantlyCheck))
	myWindow.ShowAndRun()
}

// FindProcessPidByName 返回第一个匹配进程名的 PID（找不到返回 0）
func findProcessPidByName(targetName string) (uint32, error) {
	procs, err := process.Processes()
	if err != nil {
		fmt.Printf("[findProcessPidByName] 获取进程列表失败, err: %v\n", err)
		return 0, err
	}
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if name == targetName {
			return uint32(p.Pid), nil
		}
	}
	return 0, errors.New("process not found")
}

// 打开进程并返回句柄
func getProcessHandle(pid uint32) (windows.Handle, error) {
	// 0x1F0FFF 表示对进程的完全访问权限
	handle, err := windows.OpenProcess(windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_READ|windows.PROCESS_VM_WRITE, false, pid)
	if err != nil {
		fmt.Printf("[getProcessHandle] 打开进程失败, pid: %d, err: %v\n", pid, err)
		return 0, err
	}
	if handle == 0 {
		return 0, errors.New("failed to open process: handle is 0")
	}
	return handle, nil
}

func getProcessAddress(pHandler windows.Handle) (uint64, error) {
	modules := make([]windows.Handle, 1024) // 缓冲区接收模块
	var cbNeeded uint32

	err := windows.EnumProcessModulesEx(
		pHandler,
		&modules[0],
		uint32(len(modules))*uint32(unsafe.Sizeof(modules[0])),
		&cbNeeded,
		windows.LIST_MODULES_DEFAULT,
	)
	if err != nil {
		return 0, err
	}

	if cbNeeded == 0 {
		return 0, fmt.Errorf("没有获取到模块信息")
	}

	baseAddress := uint64(modules[0]) // 第一个模块一般是主模块
	fmt.Printf("[getProcessAddress] 目标进程模块基址: 0x%X\n", baseAddress)

	return baseAddress, nil
}

func IsHandleValid(h windows.Handle) bool {
	const INVALID_HANDLE_VALUE = windows.Handle(^uintptr(0)) // -1
	return h != 0 && h != INVALID_HANDLE_VALUE
}
