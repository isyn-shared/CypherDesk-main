package main

import (
	initPkg "CypherDesk-main/init"
	//	"CypherDesk-main/model"
	"CypherDesk-main/router"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var port string = ":8080"

func rec() {
	if err := recover(); err != nil {
		fmt.Println(initPkg.TextColorFail)
		fmt.Println(err.(string))
		fmt.Print(initPkg.TextColorEnd)
	}
}

func main() {
	defer rec()
	// model.Debug()

	fmt.Print(initPkg.TextColorTitle)
	fmt.Println(" ___ ______   ___   _   ____         __ _                          ")
	fmt.Println("|_ _/ ___\\ \\ / / \\ | | / ___|  ___  / _| |___      ____ _ _ __ ___ ")
	fmt.Println(" | |\\___ \\\\ V /|  \\| | \\___ \\ / _ \\| |_| __\\ \\ /\\ / / _` | '__/ _ \\")
	fmt.Println(" | | ___) || | | |\\  |  ___) | (_) |  _| |_ \\ V  V / (_| | | |  __/")
	fmt.Println("|___|____/ |_| |_| \\_| |____/ \\___/|_|  \\__| \\_/\\_/ \\__,_|_|  \\___|")
	fmt.Println("                              CypherDesk                          ")
	fmt.Println(initPkg.TextColorEnd)

	runPtr := flag.Bool("start", false, "use to run server")
	setMysqlUserPtr := flag.Bool("set-mysql-credentials", false, "use to configure db connection")
	genAESKeysPtr := flag.Bool("genKeys", false, "use to generate keys for DB encryption")
	portPtr := flag.Int("port", 8080, "specify the port value")
	modePtr := flag.Bool("debug", true, "????")

	flag.Parse()

	if !*genAESKeysPtr && !*runPtr && !*setMysqlUserPtr {
		//panic("[!!!] Error: You may set all reqired parameters. Use -h to view usage")
		*runPtr = true
	}

	initPkg.DEBUG = *modePtr
	port = ":" + string(strconv.Itoa(*portPtr))

	consoleReader := bufio.NewReader(os.Stdin)
	if *setMysqlUserPtr {
		fmt.Print(initPkg.TextColorBold + "Username: " + initPkg.TextColorEnd)
		mysqlLogin, _ := consoleReader.ReadString('\n')
		mysqlLogin = strings.Replace(mysqlLogin, "\n", "", -1)
		fmt.Print(initPkg.TextColorBold + "Password: " + initPkg.TextColorEnd)
		bytePass, _ := terminal.ReadPassword(int(syscall.Stdin))
		mysqlPassword := strings.Replace(string(bytePass), "\n", "", -1)
		fmt.Print("\nDB name: ")
		dbName, _ := consoleReader.ReadString('\n')
		dbName = strings.Replace(dbName, "\n", "", -1)
		initPkg.SetMysqlCredentials(mysqlLogin, mysqlPassword, dbName)
	}
	if *genAESKeysPtr {
		fmt.Print(initPkg.TextColorWarning + "Be careful! After this all your data will bi lost! [Yes, I am sure]: " + initPkg.TextColorEnd)
		answer, _ := consoleReader.ReadString('\n')
		if answer == "Yes, I am sure\n" {
			fmt.Print("[i] Generating.")
			initPkg.GenerateAESKeys()
			fmt.Println("\n" + initPkg.TextColorOKGreen + "Done!" + initPkg.TextColorEnd)
		} else {
			fmt.Println("[i] Aborting" + initPkg.TextColorEnd)
		}
	}
	if *runPtr {
		router := router.New()
		initPkg.ProjectInit()
		router.Run(port)
	}
}
