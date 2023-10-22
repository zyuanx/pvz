package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

const (
	modelName = "popcapgame1.exe"
)

var (
	UnlimitedSunshineFlag = false
	NoCDFlag              = false
	AllZombieFlag         = false
)

var (
	pid         int32
	pHandler    windows.Handle
	err         error
	baseAddress uint32
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Plants vs. Zombies")
	myWindow.Resize(fyne.NewSize(300, 230))
	content := widget.NewLabel("Game not running")
	start := make(chan int)

	UnlimitedSunshineCheck := widget.NewCheck("Unlimited Sunshine", func(value bool) {
		UnlimitedSunshineFlag = value
	})
	UnlimitedSunshineCheck.Disable()
	NoCoolingCheck := widget.NewCheck("No Cooling", func(value bool) {
		NoCDFlag = value
	})
	NoCoolingCheck.Disable()
	AllZombieComingCheck := widget.NewCheck("All Zombie Coming", func(value bool) {
		AllZombieFlag = value
	})
	AllZombieComingCheck.Disable()
	KillInstantlyCheck := widget.NewCheck("Kill Instantly", func(value bool) {
		setKillInstantly(pHandler, baseAddress, value)
	})
	KillInstantlyCheck.Disable()

	go func() {
		pid = getProcessPid()
		for pid == 0 {
			pid = getProcessPid()
			time.Sleep(time.Second * 1)
		}
		content.SetText("Game is running")
		pHandler, err = getProcessHandle(pid)
		if err != nil {
			fmt.Println("[-] 获取目标进程句柄失败", err)
			return
		}
		baseAddress, err = getProcessAddress(pHandler)
		if err != nil {
			fmt.Println("[-] 获取目标进程模块基址失败", err)
			return
		}
		UnlimitedSunshineCheck.Enable()
		NoCoolingCheck.Enable()
		AllZombieComingCheck.Enable()
		KillInstantlyCheck.Enable()
		start <- 1
	}()
	go func() {
		<-start
		go setUnlimitedSunshine(pHandler, baseAddress)
		go setNoCooling(pHandler, baseAddress)
		go setAllZombieComing(pHandler, baseAddress)
	}()

	myWindow.SetContent(container.NewVBox(content, UnlimitedSunshineCheck, NoCoolingCheck, AllZombieComingCheck, KillInstantlyCheck))
	myWindow.ShowAndRun()
}

func setUnlimitedSunshine(pHandler windows.Handle, baseAddress uint32) {
	offset := []uint32{0x331C50, 0x868, 0x5578}
	sunshineNum := []byte{0x06, 0x27, 0x00, 0x00}
	tk := time.NewTicker(time.Millisecond * 500)
	for range tk.C {
		if !UnlimitedSunshineFlag {
			continue
		}
		curAddress := baseAddress
		for i, v := range offset {
			curAddress += v
			buffer := make([]byte, 4)
			if i < len(offset)-1 {
				if err := windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil); err != nil {
					log.Println(err, curAddress, buffer)
				}
				curAddress = uint32(binary.LittleEndian.Uint32(buffer))
			} else {
				if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &sunshineNum[0], uintptr(len(sunshineNum)), nil); err != nil {
					log.Println(err, curAddress, sunshineNum)
				}
			}
		}
	}
}

func setNoCooling(pHandler windows.Handle, baseAddress uint32) {
	addrs := make([]uint32, 10)
	offset := []uint32{0x331C50, 0x868, 0x15C, 0x70}
	cdLock := []byte{0x01, 0x00, 0x00, 0x00}
	tk := time.NewTicker(time.Millisecond * 500)
	for range tk.C {
		if !NoCDFlag {
			continue
		}
		p := 0 // 每个格子的偏移量
		for i := 0; i < 10; i++ {
			go func(curAddress uint32, idx int, p int) {
				for i, v := range offset {
					curAddress += v
					buffer := make([]byte, 4)
					if i < len(offset)-1 {
						windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil)
						curAddress = uint32(binary.LittleEndian.Uint32(buffer))
					} else {
						addrs[idx] = curAddress + uint32(p)
						windows.WriteProcessMemory(pHandler, uintptr(addrs[idx]), &cdLock[0], uintptr(len(cdLock)), nil)
					}
				}
			}(baseAddress, i, p)
			p += 0x50
		}
	}

}

func setAllZombieComing(pHandler windows.Handle, baseAddress uint32) {
	offset := []uint32{0x331C50, 0x868, 0x55B4}
	sunshineNum := []byte{0x01, 0x00, 0x00, 0x00}
	tk := time.NewTicker(time.Millisecond * 1)
	for range tk.C {
		if !AllZombieFlag {
			continue
		}
		curAddress := baseAddress
		for i, v := range offset {
			curAddress += v
			buffer := make([]byte, 4)
			if i < len(offset)-1 {
				if err := windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil); err != nil {
					log.Println(err, curAddress, buffer)
				}
				curAddress = uint32(binary.LittleEndian.Uint32(buffer))
			} else {
				if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &sunshineNum[0], uintptr(len(sunshineNum)), nil); err != nil {
					log.Println(err, curAddress, sunshineNum)
				}
			}
		}
	}
}

func setKillInstantly(pHandler windows.Handle, baseAddress uint32, flag bool) {
	// 秒杀
	offset := []uint32{0x145DFA}
	killVal := []byte{0x29, 0xED, 0x90, 0x90}
	if !flag {
		killVal = []byte{0x2B, 0x6C, 0x24, 0x20} // origin
	}
	curAddress := baseAddress
	curAddress += offset[0]
	if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &killVal[0], uintptr(len(killVal)), nil); err != nil {
		log.Println(err, curAddress, killVal)
	}

	// 头戴防具秒杀
	curAddress = baseAddress
	offset = []uint32{0x145B14}
	killVal = []byte{0x29, 0xC9}
	if !flag {
		killVal = []byte{0x2B, 0xC8} // origin
	}
	curAddress += offset[0]
	if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &killVal[0], uintptr(len(killVal)), nil); err != nil {
		log.Println(err, curAddress, killVal)
	}
	// 前挡防具秒杀
	curAddress = baseAddress
	offset = []uint32{0x145771}
	killVal = []byte{0x89, 0x96, 0xDC}
	if !flag {
		killVal = []byte{0x29, 0x86, 0xDC} // origin
	}
	curAddress += offset[0]
	if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &killVal[0], uintptr(len(killVal)), nil); err != nil {
		log.Println(err, curAddress, killVal)
	}
}

func getProcessPid() int32 {
	processes, _ := process.Processes()
	for _, p := range processes {
		name, _ := p.Name()
		if name == modelName {
			return p.Pid
		}
	}
	return 0
}
func getProcessHandle(pid int32) (windows.Handle, error) {
	handle, err := windows.OpenProcess(0x1F0FFF, false, uint32(pid))
	if err != nil {
		return 0, err
	}
	fmt.Println("[+] 目标进程PID", pid)
	fmt.Println("[+] 目标进程句柄", handle)
	return handle, nil

}

func getProcessAddress(pHandler windows.Handle) (uint32, error) {
	var module windows.Handle
	var cbNeeded uint32

	if err := windows.EnumProcessModulesEx(pHandler, &module, uint32(unsafe.Sizeof(module)), &cbNeeded, windows.LIST_MODULES_DEFAULT); err != nil {
		return 0, err
	}
	// exePath, _ := os.Executable()
	// modulePathUTF16 := make([]uint16, len(exePath)+1)

	// if err := windows.GetModuleBaseName(pHandler, module, &modulePathUTF16[0], uint32(len(modulePathUTF16))); err != nil {
	// 	return 0, err
	// }
	// modulePath := windows.UTF16ToString(modulePathUTF16)
	// fmt.Println(modulePath, module)
	baseAddress := uint32(module)
	// var module [1024]windows.Handle
	// var cbNeeded uint32
	// windows.EnumProcessModulesEx(process, &module[0], uint32(unsafe.Sizeof(module)), &cbNeeded, windows.LIST_MODULES_DEFAULT)
	// num := cbNeeded / uint32(unsafe.Sizeof(module[0]))
	// fmt.Println(num)
	// for i := 0; i < int(num); i++ {
	// 	exePath, _ := os.Executable()
	// 	modulePathUTF16 := make([]uint16, len(exePath)+1)
	// 	windows.GetModuleBaseName(process, module[i], &modulePathUTF16[0], uint32(len(modulePathUTF16)))
	// 	modulePath := windows.UTF16ToString(modulePathUTF16)
	// 	fmt.Println(modulePath, module[i])
	// 	if strings.EqualFold(modulePath, modelName) {
	// 		baseAddress = uint32(module[i])
	// 		break
	// 	}
	// }
	fmt.Println("[+] 目标进程模块基址", baseAddress)
	return baseAddress, nil
}
