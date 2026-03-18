package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	size := 4096 // 1 page
	shmFile, _ := os.OpenFile("/dev/shm/my_shared_mem", os.O_RDWR|os.O_CREATE, 0666)
	shmFile.Truncate(int64(size))
	defer shmFile.Close()

	data, _ := syscall.Mmap(int(shmFile.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	ptr := (*int64)(unsafe.Pointer(&data[0]))

	oldval := *ptr
	for {
		time.Sleep(time.Second)
		newval := *ptr
		if newval == oldval {
			break
		}
		oldval = newval
		fmt.Println("waiting...")
	}
	fmt.Println("engaged")
	exec.Command("wt.exe", "wsl", "go", "run", "main.go").Run()
	for {
		time.Sleep(time.Millisecond * 20)
		*ptr++
	}
}

func counting() {
	size := 4096 // 1 page
	shmFile, _ := os.OpenFile("/dev/shm/my_shared_mem", os.O_RDWR|os.O_CREATE, 0666)
	shmFile.Truncate(int64(size))
	defer shmFile.Close()

	data, _ := syscall.Mmap(int(shmFile.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	ptr := (*int64)(unsafe.Pointer(&data[0]))
	*ptr = *ptr + 1

	fmt.Println("newval", *ptr)
	time.Sleep(time.Second * 3)
	exec.Command("wt.exe", "wsl", "go", "run", "main.go").Run()
	fmt.Println("ending")
}
