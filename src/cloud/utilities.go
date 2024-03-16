package cloud

import (
	"Rapid/src/database"
	encription "Rapid/src/transfer"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
Uploads a zip to the cloud
*/
func UploadToMega(path string, from_user_id int, user_to string) (error, bool) {

	// Formats the zipped file
	split := strings.Split(path, "/")
	name_of_item := split[len(split)-1]

	// Ecncryption key (yeah ill do something about this at some point)
	key := "passphrasewhichneedstobe32bytes!"

	// Makes sure user is allowed to send the file before procceding
	if !database.PerformTransaction(from_user_id, user_to, name_of_item, database.HashInfo(key)) {
		return nil, false
	}

	// makes name include the zip
	name := fmt.Sprintf("%s.zip", name_of_item)

	// Encrypts the file
	encription.ZipEncryptFolder(path, name, key)
	location := fmt.Sprintf("../../temp/%s", name)
	mega_login := "../../.megacmd.json"
	// Sends that file to MEGA
	cmd := exec.Command("megacmd", mega_login, "put", location, "mega:/")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// Print the output
	fmt.Println(string(stdout))

	// Remove file from temp
	working_dir, _ := os.Getwd()
	temp_location := filepath.Join(working_dir, "../temp")
	dir, err := ioutil.ReadDir(temp_location)
	if err != nil {
		log.Fatal(err)
	}

	// Loops through directory and removes everything inside
	for _, d := range dir {
		err := os.RemoveAll(filepath.Join(temp_location, d.Name()))
		if err != nil {
			log.Println(err)
		}
	}
	return nil, true
}

func DownloadFromMega(mega_login, file_name string) {

	// Gets the current directory the user is in
	current_dir, _ := os.Getwd()

	// Formats it for the mega cloud
	cloud_dir := filepath.Join("mega:/", file_name)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", mega_login, "get", cloud_dir, current_dir)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// Print the output
	fmt.Println(string(stdout))
}

func DeleteFromMega(mega_login, file_name string) {

	// Formats it for the mega cloud
	cloud_dir_location := filepath.Join("mega:/", file_name)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", mega_login, "delete", cloud_dir_location)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// Print the output
	fmt.Println(string(stdout))
}
