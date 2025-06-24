package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// TODO: add for recursive watch
func watchAndRestart(path string, cmdPath string, args ...string) error {
    fd, err := unix.InotifyInit()

    if err != nil {
        return err
    }

    defer unix.Close(fd)

    wd, err := unix.InotifyAddWatch(
        fd,
        path,
        unix.IN_CREATE | unix.IN_MODIFY | unix.IN_MOVED_FROM | unix.IN_MOVED_TO,
    )

    if err != nil {
        return err
    }

    defer unix.InotifyRmWatch(fd, uint32(wd))

    var cmd *exec.Cmd

    restart := func() {
        buildCmd := exec.Command("go", "build", "-o", "app", "main.go")
        err := buildCmd.Run()

        if err != nil {
            fmt.Println("Error with build command")
        }

        if cmd != nil && cmd.Process != nil {
            fmt.Println("Killing existing process...")
            cmd.Process.Kill()
            cmd.Wait()
        }

        fmt.Println("Starting a new process...")
        cmd, err = startServer(cmdPath, args...)

        if err != nil {
            fmt.Println("Error starting process:", err)
        }
    }

    restart()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("Exiting...")
        if cmd != nil && cmd.Process != nil {
            cmd.Process.Kill()
            cmd.Wait()
        }

        os.Exit(0)
    }()

    buf := make([]byte, 4096)

    for {
        n, err := unix.Read(fd, buf)

        if err != nil {
            return err
        }

        time.Sleep(200 * time.Millisecond)

        var offset uint32 = 0

        for offset <= uint32(n - unix.SizeofInotifyEvent) {
            event := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

            if event.Len > 0 {
                restart()
                break
            }

            offset += unix.SizeofInotifyEvent + event.Len
        }
    }
}

func startServer(cmdPath string, args ...string) (*exec.Cmd, error) {
    cmd := exec.Command(cmdPath, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if err := cmd.Start(); err != nil {
        return nil, err
    }
    fmt.Println("Running process:", cmd.Process.Pid)
    return cmd, nil
}

// cheap method of reloading the server on changes
func _main() {
    err := watchAndRestart("./", "./app")
    if err != nil {
        log.Fatal(err)
    }
}
