package main

import (
	"os"
	"os/exec"
)

func setUnbuffered() {
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
}

func setBuffered() {
    exec.Command("stty", "-f", "/dev/tty", "echo", "-cbreak")
}

func getByte() byte {
    var b []byte = make([]byte, 1)
    os.Stdin.Read(b)
    return b[0]
}
