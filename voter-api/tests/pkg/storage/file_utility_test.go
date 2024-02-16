package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"drexel.edu/voter-api/pkg/storage"
	"github.com/google/uuid"
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

func TestCanCheckIfFileExists(t *testing.T) {

	exists, err := storage.CheckIfFileOrFolderExistsAndNotEmpty(projectWorkingDirectory + "/file_utility_test.go")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCanCreateFile(t *testing.T) {
	guid := uuid.New()
	path := "../../../data/" + guid.String()
	_, err := storage.CreateFile(path)
	assert.NoError(t, err)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	file.WriteString("file Can't Be Blank")
	exists, err := storage.CheckIfFileOrFolderExistsAndNotEmpty(path)
	assert.NoError(t, err)
	assert.True(t, exists)
	os.Remove(path)
}
