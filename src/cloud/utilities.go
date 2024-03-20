package cloud

import (
	"Rapid/src/database"
	encription "Rapid/src/transfer"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Generates a random 32 character string for encryption purposes
func GenerateKey() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

/*
Uploads a zip to the cloud
*/
func UploadToMega(path string, from_user_id int, user_to string) (error, bool) {

	// Formats the zipped file
	split := strings.Split(path, "/")
	name_of_item := split[len(split)-1]

	// Ecncryption key (randomly generated)
	key, _ := GenerateKey()

	// Makes sure user is allowed to send the file before procceding

	if !database.PerformTransaction(from_user_id, user_to, name_of_item, key) {
		return nil, false
	}

	hashed_key := database.HashInfo(key)
	fmt.Println(hashed_key)
	// makes name include the zip
	name := fmt.Sprintf("%s_%s.zip", hashed_key, database.HashInfo(name_of_item))
	fmt.Println(path)

	// Encrypts the file
	encription.ZipEncryptFolder(path, name, key)
	location := fmt.Sprintf("../temp/%s", name)
	fmt.Println(location)
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

func DownloadFromMega(user int, file_name string, location string) (error, bool) {

	if !database.UserCanViewTransaction(user, file_name) {
		return nil, false
	}
	fmt.Println("Made it past here")
	// Ecncryption key
	key := database.RetrieveKey(file_name, user)

	// Gets the current directory the user is in
	current_dir, _ := os.Getwd()

	// Destination the file will be downloaded to
	destination := filepath.Join(current_dir, location)
	fmt.Println(destination)

	// Formats it for the mega cloud (readjusts the name to fit the hashing)
	cloud_dir := fmt.Sprintf("mega:/%s_%s.zip", database.HashInfo(key), database.HashInfo(file_name))

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
	encryped_name := fmt.Sprintf("%s_%s.zip", database.HashInfo(key), database.HashInfo(file_name))
	file_location := filepath.joinpath(destination, encryped_name)
	// Decripts folder
	err = encription.DecryptZipFolder(file_location, file_name, key)
	if err != nil {
		fmt.Println(err)
	}

	// Deletes zip folder
	err = os.RemoveAll(destination)
	if err != nil {
		log.Println(err)
	}

	// Removes the copy from the cloud so that no users can access it
	//DeleteFromMega(user, file_name)
	return nil, true
}

// Removes the file from the cloud
func DeleteFromMega(user int, file_name string) {
	// Ecncryption key
	key := database.RetrieveKey(file_name, user)

	hashed_key := database.HashInfo(key)

	// Formats it for the mega cloud
	cloud_dir := fmt.Sprintf("mega:/%s_%s.zip", hashed_key, database.HashInfo(file_name))

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
