package encription

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
)

// ZipEncryptFolder zips and encrypts a folder or path
func ZipEncryptFolder(sourceFolder, fileName, key string) error {

	// Items in temp folder will be deleted right after
	zipFilePath := filepath.Join("../../temp", fileName)

	// Create a zip file
	zipFileWriter, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFileWriter.Close()

	// Initialize AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFileWriter)
	defer zipWriter.Close()

	// Walk through the source folder
	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a new file header
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the file within the zip
		fileHeader.Name, err = filepath.Rel(sourceFolder, path)
		if err != nil {
			return err
		}

		// Create a new entry in the zip writer
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		// Open the source file
		inputFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		// Encrypt and write the file contents
		if err := encryptCopy(writer, inputFile, block); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// encryptCopy encrypts the file or folder selected
func encryptCopy(dst io.Writer, src io.Reader, block cipher.Block) error {

	// Creates a vector
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	// Write initialization vector to destination
	if _, err := dst.Write(iv); err != nil {
		return err
	}

	// Create cipher
	stream := cipher.NewCFBEncrypter(block, iv)

	// Create and use a buffer
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Encrypt and write data
		stream.XORKeyStream(buf[:n], buf[:n])
		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
