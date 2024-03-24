package cloud

import (
	database "Rapid/src/api"
	custom "Rapid/src/handling"
	encription "Rapid/src/transfer"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

	if from_user_id == 0 {
		return custom.NewError("User must be logged in to use this method"), false
	}

	// Formats the zipped file
	split := strings.Split(path, "/")
	name_of_item := split[len(split)-1]

	// Ecncryption key (randomly generated)
	key, _ := GenerateKey()

	// Makes sure user is allowed to send the file before procceding
	result, err := database.PerformTransaction(from_user_id, user_to, name_of_item, key)
	if err != nil {
		return err, false
	}
	if !result {
		return nil, false
	}

	hashed_key := database.HashInfo(key)
	// makes name include the zip
	name := fmt.Sprintf("%s_%s.zip", hashed_key, database.HashInfo(name_of_item))

	// Encrypts the file
	encription.ZipEncryptFolder(path, name, key)

	working_dir, _ := os.Getwd()

	// Handles megacmd config
	home, _ := os.UserHomeDir()
	directory := filepath.Join(home, "Rapid/.megacmd.json")
	config := fmt.Sprintf(`-conf=%s`, directory)

	// Sends that file to MEGA
	cmd := exec.Command("megacmd", config, "put", name, "mega:/")
	cmd.Dir = working_dir
	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err = cmd.Run()
	if err != nil {
		return err, false
	}

	// Remove file from temp
	temp_zip := filepath.Join(working_dir, name)

	// Deletes zip folder
	err = os.RemoveAll(temp_zip)
	if err != nil {
		log.Println(err)
	}

	return nil, true
}

func DownloadFromMega(user int, file_name string, location string) (error, bool) {

	if user == 0 {
		return custom.NewError("User must be logged in to use this method"), false
	}

	if !database.UserCanViewTransaction(user, file_name) {
		return nil, false
	}
	// Ecncryption key
	key, err := database.RetrieveKey(file_name, user)
	if err != nil {
		return err, false
	}
	// Gets the current directory the user is in
	current_dir, _ := os.Getwd()
	encryped_name := fmt.Sprintf("%s_%s.zip", database.HashInfo(key), database.HashInfo(file_name))

	// Destination the file will be downloaded to
	destination := filepath.Join(current_dir, location, encryped_name)

	// Formats it for the mega cloud (readjusts the name to fit the hashing)
	cloud_dir := fmt.Sprintf("mega:/%s_%s.zip", database.HashInfo(key), database.HashInfo(file_name))

	// Handles megacmd config
	home, _ := os.UserHomeDir()
	directory := filepath.Join(home, "Rapid/.megacmd.json")
	config := fmt.Sprintf(`-conf=%s`, directory)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", config, "get", cloud_dir, destination)
	cmd.Dir = current_dir

	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	// Decripts folder
	err = encription.DecryptZipFolder(destination, file_name, key)
	if err != nil {
		fmt.Println(err)
	}

	// Deletes zip folder
	err = os.RemoveAll(destination)
	if err != nil {
		log.Println(err)
	}

	// Removes the copy from the cloud so that no users can access it
	_, err = DeleteFromMega(user, file_name)
	if err != nil {
		return err, false
	}
	return nil, true
}

// Removes the file from the cloud
func DeleteFromMega(user int, file_name string) (bool, error) {

	if user == 0 {
		return false, custom.NewError("User must be logged in to use this method")
	}

	// Ecncryption key
	key, err := database.RetrieveKey(file_name, user)
	if err != nil {
		return false, err
	}
	hashed_key := database.HashInfo(key)

	// Formats it for the mega cloud
	cloud_dir := fmt.Sprintf("mega:/%s_%s.zip", hashed_key, database.HashInfo(file_name))

	// Handles megacmd config
	home, _ := os.UserHomeDir()
	directory := filepath.Join(home, "Rapid/.megacmd.json")
	config := fmt.Sprintf(`-conf=%s`, directory)

	// Calls cmd command to retrieve the file
	cmd := exec.Command("megacmd", config, "delete", cloud_dir)

	// Error handing
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Runs cmd command
	err = cmd.Run()
	if err != nil {
		return false, err
	}
	err = database.DeleteTransaction(key)

	if err != nil {
		return false, nil
	}
	return true, nil

}
