package main

import (
	"fmt"
	"log"
	"time"
)

var x = 0.0
var y = 0.0
var z = 0.0

type Vector2f struct {
    X, Y float64
}

func logListen(ch <-chan string) {
	for message := range ch {
		log.Println("<<<", string(message))
	}
}

func getPoint(cs *GCodeCommandSender) Vector2f {
    setUnbuffered()
    defer setBuffered()

    key := 0
    for key != 10 { // enter
        send := true
        key = int(getByte())

        switch key {
            case 68: // left
                x -= 1
            case 67: // right
                x += 1
            case 66: // down
                y += 1
            case 65: // up
                y -= 1
            case 10: break 
            default: send = false
        }

        if send {
            cs.SendCommand(MoveXY(x, y))
        }
    }

    return Vector2f{X: float64(x), Y: float64(y)}
}

func getCorrectZ(cs *GCodeCommandSender) float64 {
    setUnbuffered()
    defer setBuffered()

    key := 0
    for key != 10 { // enter
        send := true
        key = int(getByte())

        switch key {
            case 66: // down
                z -= 0.05 
            case 65: // up
                z += 0.05 
            case 10: break 
            default: send = false
        }

        if send {
            cs.SendCommand(MoveZ(z))
        }
    }

    return z
}

func main() {
    ports := getSerialPorts()
    points := 3 
    
    fmt.Println("Please enter your bed size (i.e 200x200)")
    bed := Vector2f{}
    fmt.Scan(&bed.X)
    fmt.Scan(&bed.Y)

    fmt.Println("List of serial ports:")
    for i, port := range ports {
        fmt.Println(i, ":", port)
    }

    port := 0
    fmt.Scanf("%d", &port)

    baudrate := 0
    fmt.Println("Which baudrate?")
    fmt.Scanf("%d", &baudrate)

    log.Println("Attempting to connect to", ports[port], "at", baudrate, "baud/s")
    connection, err := NewConnection(ports[port], baudrate)
    if err != nil {
        log.Fatal(err)
    }

    listen := make(chan string)
    connection.AddListener(listen)
    go logListen(listen)

    cs := NewGCodeCommandSender(connection)
    time.Sleep(5 * time.Second)
    cs.SendCommand(SetCurrentLineNumber)
    cs.SendCommand(MoveZ(3))
    cs.SendCommand(HomeAllAxes)
    cs.WG.Wait()


    fmt.Println("Please move to initial XY position")
    initial := getPoint(cs)
    fmt.Println("x_0 =", initial)

    position := initial
    y_steps := 0
    for position.Y <= bed.Y + initial.Y {
        for i := 0; i < points; i++ {
            position.X = initial.X + (bed.X / float64(points - 1)) * float64(i)
            fmt.Println("Moving to position")
            cs.SendCommand(MoveXY(position.X, position.Y))

            fmt.Println("Move z axis until paper test is ok")
            height := getCorrectZ(cs)
            fmt.Println(position, "->", height)
        }

        y_steps += 1
        position.X = initial.X
        position.Y = initial.Y + (bed.Y / float64(points - 1)) * float64(y_steps)
    }

    cs.WG.Wait()
}
