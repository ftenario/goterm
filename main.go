/*
  goterm - is a terminal application to manage your serial devices.
  author: Ferdinand Enario ftenario@gmail.com
*/

package main
import (
        "log"
        "time"
        "github.com/tarm/serial"
        "bufio"
        "fmt"
        "os"
)
const version = "0.0.0.1"

/*
  function name : getInput
  parameters: channel string
  purpose: get the input from the user
  and send it to a channel
*/
func getInput(input chan string) {
  for {
    reader := bufio.NewReader(os.Stdin)
    d,_ := reader.ReadString('\n')
    input <- d
  }
}

/*
  function name: main
  - the main function of the application
*/
func main() {

  //create a serial object
  c := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600, ReadTimeout: time.Millisecond * 25}
  s, err := serial.OpenPort(c)
  if err != nil {
          log.Fatal(err)
          fmt.Println("Error opeing serial port...")
  }
  defer s.Close()

    rxChan := make(chan string)
    txChan := make(chan string)

    //create go routine for transmitting data
    var m string
    go func() {
      for {
        time.Sleep(time.Millisecond)
        m = <- txChan
        s.Write([]byte(m))
      }
    }()

    //go routine to read data from serial port
    go func() {
      serial := bufio.NewReader(s)
      for {
          time.Sleep(50 * time.Millisecond)

          //read until newline
          recv,_ := serial.ReadBytes('\x0a')
          rxChan <- string(recv)

        }
    }()

    //This is your main loop
    input := make(chan string)
    go getInput(input)
    fmt.Println("\n->")

    for {
      time.Sleep(50 * time.Millisecond)

      select {

        case tx := <-input:
          txChan <- tx

        case r := <- rxChan:
          if len(r) > 0 {
            fmt.Printf("%s", r)
          }

        case <- time.After(4000 * time.Millisecond):
          fmt.Println("timeout\n")
      }
    }
}
