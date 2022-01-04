package config

import (
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitializeBoltConfig(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, dbErr := InitializeBoltConfig(fileName.Name())

	if dbErr != nil {
		t.Errorf("Unexpected error: %s", dbErr.Error())
	}

	_, fileErr := os.Stat(fileName.Name())
	if fileErr != nil {
		t.Errorf("Unexpected file error: %s", fileErr.Error())
	}

	conn := db.openDB()
	defer closeDB(conn)

	_ = conn.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)

		if projects == nil {
			t.Errorf("Projects bucket not found")
		}

		return nil
	})

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_AddDockerComposeFile(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	err := db.AddDockerComposeFile("aFileName", "aProjectName")
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	conn := db.openDB()
	defer closeDB(conn)

	conn.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte("aProjectName"))
		files := project.Bucket(bucketNameFiles)

		firstFile := files.Bucket([]byte("1"))

		if firstFile == nil {
			t.Errorf("No file bucket found")
		}

		name := firstFile.Get(bucketKeyFileName)

		if string(name) != "aFileName" {
			t.Errorf("Invalid file name. Expected %s, got %s", "aFileName", string(name))
		}

		return nil
	})

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetDockerComposeFilesByProject(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	_ = db.AddDockerComposeFile("aFileName2", "aProjectName")

	project, _ := db.GetDockerComposeFilesByProject("aProjectName")

	if project == nil {
		t.Errorf("Expected project files, got nil")
	}

	if len(project) != 2 {
		t.Errorf("Invalid files count. Expected %d, got %d", 2, len(project))
	}

	if project[0].GetFilename() != "aFileName" {
		t.Errorf("Invalid file name. Expected %s, got %s", "aFileName", project[0].GetFilename())
	}

	if project[1].GetFilename() != "aFileName2" {
		t.Errorf("Invalid file name. Expected %s, got %s", "aFileName2", project[1].GetFilename())
	}

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetDockerComposeProjectList(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	_ = db.AddDockerComposeFile("aFileName", "aProjectName2")

	projectList, _ := db.GetDockerComposeProjectList("")

	if len(projectList) != 2 {
		t.Errorf("Invalid project list count. Expected %d, got %d", 2, len(projectList))
	}

	if projectList[0] != "aProjectName" {
		t.Errorf("Invalid project name. Expected %s, got %s", "aProjectName", projectList[0])
	}

	if projectList[1] != "aProjectName2" {
		t.Errorf("Invalid project name. Expected %s, got %s", "aProjectName2", projectList[1])
	}

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetDockerComposeProjectList_duplicate(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	_ = db.AddDockerComposeFile("aFileName2", "aProjectName")

	projectList, _ := db.GetDockerComposeProjectList("")

	if len(projectList) != 1 {
		t.Errorf("Invalid project list count. Expected %d, got %d", 1, len(projectList))
	}

	if projectList[0] != "aProjectName" {
		t.Errorf("Invalid project name. Expected %s, got %s", "aProjectName", projectList[0])
	}

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetDockerComposeProjectList_filtered(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	_ = db.AddDockerComposeFile("aFileName2", "differentProjectName")

	projectList, _ := db.GetDockerComposeProjectList("different")

	if len(projectList) != 1 {
		t.Errorf("Invalid project list count. Expected %d, got %d", 1, len(projectList))
	}

	if projectList[0] != "differentProjectName" {
		t.Errorf("Invalid project name. Expected %s, got %s", "differentProjectName", projectList[0])
	}

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_DeleteProjectByName(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	_ = db.AddDockerComposeFile("aFileName2", "aProjectName")

	_ = db.DeleteProjectByName("aProjectName")

	projectList, _ := db.GetDockerComposeProjectList("")

	if len(projectList) != 0 {
		t.Errorf("Invalid project list count. Expected %d, got %d", 0, len(projectList))
	}

	_ = os.Remove(fileName.Name())
}
