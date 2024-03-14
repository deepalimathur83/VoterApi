package tests

import (
	"log"
	"os"
	"testing"

	"drexel.edu/voter-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var projectWorkingDirectory string

func init() {
	log.Println("initializing file utility tests...")
	var err error
	projectWorkingDirectory, err = os.Getwd()
	if err != nil {
		panic("Could not open the working directory in the file utility test")
	}
}

func TestCanCheckIfDirectoryExists(t *testing.T) {

	projectWorkingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	exists, err := storage.CheckIfFileOrFolderExistsAndNotEmpty(projectWorkingDirectory)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCheckIfFileOrFolderExistsAndNotEmpty(t *testing.T) {

	exists, err := storage.CheckIfFileOrFolderExistsAndNotEmpty("nonexistentfile.txt")
	if exists || err != nil {
		t.Errorf("Expected false and nil, got %v and %v", exists, err)
	}

	emptyFilePath := "testfile.txt"
	emptyFile, _ := os.Create(emptyFilePath)
	emptyFile.Close()
	exists, err = storage.CheckIfFileOrFolderExistsAndNotEmpty(emptyFilePath)
	if exists || err != nil {
		t.Errorf("Expected false and nil, got %v and %v", exists, err)
	}
	os.Remove(emptyFilePath)

	nonEmptyFilePath := "testfile.txt"
	nonEmptyFile, _ := os.Create(nonEmptyFilePath)
	nonEmptyFile.WriteString("test content")
	nonEmptyFile.Close()
	exists, err = storage.CheckIfFileOrFolderExistsAndNotEmpty(nonEmptyFilePath)
	if !exists || err != nil {
		t.Errorf("Expected true and nil, got %v and %v", exists, err)
	}
	os.Remove(nonEmptyFilePath)

}

func TestCreateFile(t *testing.T) {

	path := "nonexistentdir/testfile.txt"
	file, err := storage.CreateFile(path)
	if file == nil || err != nil {
		t.Errorf("Expected file and nil, got %v and %v", file, err)
	}
	file.Close()
	os.RemoveAll("nonexistentdir")
}
