package cloud

import (
	encription "Rapid/src/transfer"
	"bytes"
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
	/*
		if !database.PerformTransaction(from_user_id, user_to, name_of_item, database.HashInfo(key)) {
			return nil, false
		}
	*/
	// makes name include the zip
	name := fmt.Sprintf("%s.zip", name_of_item)

	// Encrypts the file
	encription.ZipEncryptFolder(path, name, key)
	location := fmt.Sprintf("../temp/%s", name)

	// Sends that file to MEGA
	cmd := exec.Command("megacmd", "put", location, "mega:/")

	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err := cmd.Run()
	fmt.Println(cmd)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

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

func DownloadFromMega(file_name string, location string) {

	// Gets the current directory the user is in
	current_dir, _ := os.Getwd()

	// Destination the file will be downloaded to
	destination := filepath.Join(current_dir, location, file_name)
	fmt.Println(destination)
	// Formats it for the mega cloud
	cloud_dir := fmt.Sprintf("mega:/%s", file_name)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", "get", cloud_dir, destination)
	cmd.Dir = current_dir

	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err := cmd.Run()
	fmt.Println(cmd)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	// Ecncryption key (yeah ill do something about this at some point)
	key := "passphrasewhichneedstobe32bytes!"

	// Reverts to original name by splitting at zip
	original_name := strings.Split(file_name, ".zip")

	// Decripts folder
	err = encription.DecryptZipFolder(destination, original_name[0], key)
	if err != nil {
		fmt.Println(err)
	}

	// Deletes zip folder
	err = os.RemoveAll(destination)
	if err != nil {
		log.Println(err)
	}

}

func DeleteFromMega(file_name string) {

	// Formats it for the mega cloud
	cloud_dir := fmt.Sprintf("mega:/%s", file_name)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", "delete", cloud_dir)

	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err := cmd.Run()
	fmt.Println(cmd)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

}
