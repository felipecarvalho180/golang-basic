package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoringCount = 3
const delay = 5

func main() {
	intro()

	for {
		showMenu()
		command := getCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			showLogs()
		case 0:
			fmt.Println("See you later")
			os.Exit(0)
		default:
			fmt.Println("I don't know this command")
			os.Exit(-1)
		}
	}
}

func intro() {
	name := "Felipe"
	version := 1.1

	fmt.Println("Hello World!", name)
	fmt.Println("The is in the version:", version)
}

func showMenu() {
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show Logs")
	fmt.Println("0 - Logout")
}

func getCommand() int {
	var command int
	fmt.Scan(&command)
	fmt.Println("")
	return command
}

func startMonitoring() {
	fmt.Println("Monitoring...")

	urls := readUrlsFile()

	for i := 0; i < monitoringCount; i++ {
		for _, url := range urls {
			testSite(url)
		}

		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func testSite(url string) {
	fmt.Println("Testing Site:", url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error testing site:", url, err.Error())
		registerLog(url, false)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Println("The site", url, "is online")
		registerLog(url, true)
	} else {
		fmt.Println("The site", url, "is offline. Error:", response.StatusCode)
		registerLog(url, false)
	}
}

func readUrlsFile() []string {
	var urls []string

	file, err := os.Open("urls.txt")

	if err != nil {
		fmt.Println("Error reading urls.txt:", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		urls = append(urls, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return urls
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error reading logs.txt:", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func showLogs() {
	fmt.Println("Showing Logs...")
	fmt.Println()

	file, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(file))
}
