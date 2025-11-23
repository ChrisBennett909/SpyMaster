package main


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"github.com/fsnotify/fsnotify"
)

func NewOperation(dirToWatch string, done chan bool){
	spyMaster, err := fsnotify.NewWatcher()
	if err != nil{
		log.Fatal("Failed to make new watcher")
		return
	}
	defer spyMaster.Close()

	err = spyMaster.Add(dirToWatch)
	if err != nil{
		log.Fatal("Failed to add directory")
		return
	}

	fmt.Println("Spy Master in position")

	for {
	    select{
		case <-done:
		    fmt.Println("Stopping Watcher...")
		    return 

		case event, okay := <-spyMaster.Events:
		    if !okay{
			return
		    }

		    if event.Op&fsnotify.Create == fsnotify.Create{
			info, err := os.Stat(event.Name)
			if err != nil{
			    log.Println("Failed to read event data", err)
				continue
			    }

			    if info.IsDir(){
				fmt.Println("Folder Created: ", event.Name)
				logReport("Folder Created: " + event.Name)

				dirErr := spyMaster.Add(event.Name)
				if dirErr != nil{
					log.Println("Failed to add new folder to  watchlist")
				}
			    }else{
				fmt.Println("File Created: ", event.Name)
				logReport("File Created: " + event.Name)
			    }
		    }

		    if event.Op&fsnotify.Write == fsnotify.Write{
			fmt.Println("File modified: ", event.Name)
			logReport("File Modified: " + event.Name)
		    }

		    if event.Op&fsnotify.Rename == fsnotify.Rename{
			fmt.Println("File/Folder renamed or moved", event.Name)
			logReport("File/Folder Renamed or Moved: " + event.Name)
		    }

		case err, okay := <-spyMaster.Errors:
		    if !okay{
			return
		    }
		    log.Println("Watcher Error: ", err)
	    }
	}
}

func logReport(note string){
	timeNow := time.Now()
	timeStamp := timeNow.Format("2006-01-02 15:04:05")

	report := note + " " + timeStamp + "\n"

	file, err := os.OpenFile("spyReport.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	    log.Println("Error Opening report file")
	    return
	}
	defer file.Close()

	_, err = file.WriteString(report)
	if err != nil{
	    log.Println("Failed to write report")
	    return 
	}
}

func readTerminal(ch chan string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
	    fmt.Println(">>")
	    if scanner.Scan(){
		input := strings.TrimSpace(scanner.Text())
		ch <- input
	    }

	    err := scanner.Err(); 
	    if err != nil{
		log.Println("Error reading input: ", err)
	    }
	}
}

func main(){
    done := make(chan bool)
    inputChan := make(chan string)

    go readTerminal(inputChan)

    go NewOperation("./Storage", done)

    for input := range inputChan {
	if input =="exit"{
	    fmt.Println("Exiting ...")
	    close(done)
	    close(inputChan)
	    break
	}else {
		fmt.Println("Echo: ", input)
	}
    }
}
