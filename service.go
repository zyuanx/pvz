package main

import (
	"encoding/binary"
	"log"
	"time"

	"golang.org/x/sys/windows"
)

var (
	UnlimitedSunshineChan chan struct{}
	NoCDChan              chan struct{}
	AllZombieChan         chan struct{}
)

func setUnlimitedSunshineFunc(value bool) {
	if !value {
		if UnlimitedSunshineChan != nil {
			close(UnlimitedSunshineChan)
		}
		return
	}
	UnlimitedSunshineChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-UnlimitedSunshineChan:
				return
			default:
				if IsHandleValid(pHandler) {
					setUnlimitedSunshine(pHandler, baseAddress)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

func setUnlimitedSunshine(pHandler windows.Handle, baseAddress uint64) {
	offset := []uint64{0x331C50, 0x868, 0x5578}
	sunshineNum := []byte{0x06, 0x27, 0x00, 0x00}
	curAddress := baseAddress
	for i, v := range offset {
		curAddress += v
		buffer := make([]byte, 4)
		if i < len(offset)-1 {
			if err := windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil); err != nil {
				log.Printf("[setUnlimitedSunshine] ReadProcessMemory, err: %v, address: 0x%X, buffer: %v\n", err, curAddress, buffer)
			}
			curAddress = uint64(binary.LittleEndian.Uint32(buffer))
		} else {
			if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &sunshineNum[0], uintptr(len(sunshineNum)), nil); err != nil {
				log.Printf("[setUnlimitedSunshine] WriteProcessMemory, err: %v, address: 0x%X, buffer: %v\n", err, curAddress, sunshineNum)
			}
		}
	}
}

func setNoCoolingFunc(value bool) {
	if !value {
		if NoCDChan != nil {
			close(NoCDChan)
		}
		return
	}
	NoCDChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-NoCDChan:
				return
			default:
				if IsHandleValid(pHandler) {
					setNoCooling(pHandler, baseAddress)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

func setNoCooling(pHandler windows.Handle, baseAddress uint64) {
	addrs := make([]uint64, 10)
	offset := []uint64{0x331C50, 0x868, 0x15C, 0x70}
	cdLock := []byte{0x01, 0x00, 0x00, 0x00}

	p := 0 // 每个格子的偏移量
	for i := 0; i < 10; i++ {
		go func(curAddress uint64, idx int, p int) {
			for i, v := range offset {
				curAddress += v
				buffer := make([]byte, 4)
				if i < len(offset)-1 {
					windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil)
					curAddress = uint64(binary.LittleEndian.Uint32(buffer))
				} else {
					addrs[idx] = curAddress + uint64(p)
					windows.WriteProcessMemory(pHandler, uintptr(addrs[idx]), &cdLock[0], uintptr(len(cdLock)), nil)
				}
			}
		}(baseAddress, i, p)
		p += 0x50
	}

}

func setAllZombieComingFunc(value bool) {
	if !value {
		if AllZombieChan != nil {
			close(AllZombieChan)
		}
		return
	}
	AllZombieChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-AllZombieChan:
				return
			default:
				if IsHandleValid(pHandler) {
					setAllZombieComing(pHandler, baseAddress)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func setAllZombieComing(pHandler windows.Handle, baseAddress uint64) {
	offset := []uint64{0x331C50, 0x868, 0x55B4}
	sunshineNum := []byte{0x01, 0x00, 0x00, 0x00}

	curAddress := baseAddress
	for i, v := range offset {
		curAddress += v
		buffer := make([]byte, 4)
		if i < len(offset)-1 {
			if err := windows.ReadProcessMemory(pHandler, uintptr(curAddress), &buffer[0], uintptr(len(buffer)), nil); err != nil {
				log.Println(err, curAddress, buffer)
			}
			curAddress = uint64(binary.LittleEndian.Uint32(buffer))
		} else {
			if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &sunshineNum[0], uintptr(len(sunshineNum)), nil); err != nil {
				log.Println(err, curAddress, sunshineNum)
			}
		}
	}
}

func setKillInstantlyFunc(value bool) {
	if IsHandleValid(pHandler) {
		setKillInstantly(pHandler, baseAddress, value)
	}
}

func setKillInstantly(pHandler windows.Handle, baseAddress uint64, flag bool) {
	// 秒杀
	offset := []uint64{0x145DFA}
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
	offset = []uint64{0x145B14}
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
	offset = []uint64{0x145771}
	killVal = []byte{0x89, 0x96, 0xDC}
	if !flag {
		killVal = []byte{0x29, 0x86, 0xDC} // origin
	}
	curAddress += offset[0]
	if err := windows.WriteProcessMemory(pHandler, uintptr(curAddress), &killVal[0], uintptr(len(killVal)), nil); err != nil {
		log.Println(err, curAddress, killVal)
	}
}
