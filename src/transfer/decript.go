package encription

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
	"path/filepath"
)

// DecryptZipFolder decrypts and unzips the contents of a zip
func DecryptZipFolder(Name, destFolder, key string) error {

	// Open the zip file
	zipFileReader, err := zip.OpenReader(Name)
	if err != nil {
		return err
	}
	defer zipFileReader.Close()

	// AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	// Walk through each file in the zip
	for _, file := range zipFileReader.File {
		// Open the file from the zip
		inputFile, err := file.Open()
		if err != nil {
			return err
		}
		defer inputFile.Close()

		// Create the destination file
		outputPath := filepath.Join(destFolder, file.Name)
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Decrypt and copy file contents
		if err := decryptCopy(outputFile, inputFile, block); err != nil {
			return err
		}
	}

	return nil
}

// decrypts the zip and writes it to the destination.
func decryptCopy(dst io.Writer, src io.Reader, block cipher.Block) error {
	// Read vector
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(src, iv); err != nil {
		return err
	}

	// cipher stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Create buffer
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Decrypt and write data
		stream.XORKeyStream(buf[:n], buf[:n])
		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
