package config

import (
	docker_compose_manager "docker-compose-manager/src/docker-compose-manager"
	"docker-compose-manager/src/tests"
	"io/ioutil"
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
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

func TestBoltConfigStorage_SaveExecConfig_ForNonExistingProject(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	config := docker_compose_manager.InitProjectExecConfig("container", "command")
	db.SaveExecConfig(config, "aProjectName")

	conn := db.openDB()
	defer closeDB(conn)

	conn.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte("aProjectName"))
		config := project.Bucket([]byte(bucketNameProjectConfig))

		tests.AssertStringEquals(t, "container", string(config.Get(bucketKeyContainerName)), "Container name string")
		tests.AssertStringEquals(t, "command", string(config.Get(bucketKeyCommand)), "Command string")

		return nil
	})

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_SaveExecConfig_ForExistingProject(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())
	db.AddDockerComposeFile("fileName", "aProjectName")

	config := docker_compose_manager.InitProjectExecConfig("container", "command")
	db.SaveExecConfig(config, "aProjectName")

	conn := db.openDB()
	defer closeDB(conn)

	conn.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte("aProjectName"))
		config := project.Bucket([]byte(bucketNameProjectConfig))

		tests.AssertStringEquals(t, "container", string(config.Get(bucketKeyContainerName)), "Container name string")
		tests.AssertStringEquals(t, "command", string(config.Get(bucketKeyCommand)), "Command string")

		return nil
	})

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetExecConfigByProjectForEmptyExistingProject(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())
	_ = db.AddDockerComposeFile("aFileName", "aProjectName")

	_, err := db.GetExecConfigByProject("aProjectName")

	tests.AssertErrorEquals(t, "no config found", err)

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetExecConfigByProject(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())
	_ = db.AddDockerComposeFile("aFileName", "aProjectName")
	config := docker_compose_manager.InitProjectExecConfig("container", "command")
	db.SaveExecConfig(config, "aProjectName")

	configEntry, err := db.GetExecConfigByProject("aProjectName")

	tests.AssertNil(t, err, "config retrieval error")

	tests.AssertStringEquals(t, "container", configEntry.GetContainerName(), "container name")
	tests.AssertStringEquals(t, "command", configEntry.GetCommand(), "command")

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_StoreSettingsEntry(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	db.StoreSettingsEntry("anyKey", "anyValue")

	conn := db.openDB()
	defer closeDB(conn)

	var result []byte
	conn.View(func(tx *bolt.Tx) error {
		settings := tx.Bucket(bucketNameSettings)
		result = settings.Get([]byte("anyKey"))

		tests.AssertStringEquals(t, "anyValue", string(result), "stored settings value")

		return nil
	})

	_ = os.Remove(fileName.Name())
}
func TestBoltConfigStorage_StoreSettingsEntry_changeValue(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	db.StoreSettingsEntry("anyKey", "anyValue")
	db.StoreSettingsEntry("anyKey", "changedValue")

	conn := db.openDB()
	defer closeDB(conn)

	var result []byte
	conn.View(func(tx *bolt.Tx) error {
		settings := tx.Bucket(bucketNameSettings)
		result = settings.Get([]byte("anyKey"))

		tests.AssertStringEquals(t, "changedValue", string(result), "stored settings value")

		return nil
	})

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetSettingsEntry(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	conn := db.openDB()
	conn.Update(func(tx *bolt.Tx) error {
		settings, _ := tx.CreateBucketIfNotExists(bucketNameSettings)
		return settings.Put([]byte("someKey"), []byte("someValue"))
	})
	closeDB(conn)

	result, err := db.GetSettingsEntry("someKey")

	tests.AssertStringEquals(t, "someValue", string(result), "retrieved settings value")
	tests.AssertNil(t, err, "setting retrieval error")

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetSettingsEntry_noSettingsTable(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	result, err := db.GetSettingsEntry("someKey")

	tests.AssertStringEquals(t, "", result, "setting retrieval result on wrong key")
	tests.AssertErrorEquals(t, "key does not exists", err)

	_ = os.Remove(fileName.Name())
}

func TestBoltConfigStorage_GetSettingsEntry_invalidKey(t *testing.T) {
	fileName, _ := ioutil.TempFile(os.TempDir(), "db-")
	db, _ := InitializeBoltConfig(fileName.Name())

	db.StoreSettingsEntry("anyKey", "anyValue")
	result, err := db.GetSettingsEntry("someKey")

	tests.AssertStringEquals(t, "", result, "setting retrieval result on wrong key")
	tests.AssertErrorEquals(t, "key does not exists", err)

	_ = os.Remove(fileName.Name())
}
