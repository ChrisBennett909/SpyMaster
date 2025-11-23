package main


import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/fsnotify/fsnotify"
)

func NewOperation(dirToWatch string){
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

	done := make(chan bool)
	
	go func(){
		for {
		    select{
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
	}()

	<-done
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

func main(){
    NewOperation("./testFolder")
}
