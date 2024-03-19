package main

import "Rapid/src/cloud"

/*
// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()
	user := database.GetCurrentId()
	fmt.Printf("Currently Logged in as %s\n", database.GetUserNameByID(database.GetCurrentId()))

	// Creates a splice from the command line input
	splice := os.Args[1:]
	arguments := strings.Join(splice, " ")
	for {
		// Retreives flags called and then runs the commands inside of them
		//fmt.Scanf(&arguments)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		arguments = scanner.Text()
		if strings.Compare(arguments, "quit") == 0 {
			fmt.Println("You have chosen to exit the program")
			break
		}
		// Retreives flags called and then runs the commands inside of them
		flags := retrieveFlags(arguments)
		command(flags, user)
		fmt.Print(">> ")
	}
}
*/

func main() {
	/*
		key := "passphrasewhichneedstobe32bytes!"
		encription.ZipEncryptFolder("../testing", "testing.zip", key)
		if err := encription.DecryptZipFolder("testing.zip", "../output", key); err != nil {
			fmt.Println("Error decrypting and unzipping folder:", err)
			return
		}
	*/
	//cloud.UploadToMega("../testing", 1, "adam")
	//cloud.DownloadFromMega("testing.zip", "")
	cloud.DeleteFromMega("testing.zip")
}
