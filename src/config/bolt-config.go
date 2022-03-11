package config

import (
	"bytes"
	dcf "docker-compose-manager/src/docker-compose-manager"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type BoltConfigStorage struct {
	filePath string
}

func InitializeBoltConfig(path string) (BoltConfigStorage, error) {
	db := BoltConfigStorage{
		filePath: path,
	}

	connection := db.openDB()
	defer closeDB(connection)

	err := connection.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketNameProjects)
		return err
	})

	return db, err
}

var bucketNameProjects = []byte("Projects")
var bucketNameFiles = []byte("Files")
var bucketKeyFileName = []byte("fileName")
var bucketNameProjectConfig = []byte("projectConfig")
var bucketKeyContainerName = []byte("containerName")
var bucketKeyCommand = []byte("command")

func (c *BoltConfigStorage) AddDockerComposeFile(file, projectName string) error {
	db := c.openDB()
	defer closeDB(db)

	return db.Update(func(tx *bolt.Tx) error {
		projects, err := tx.CreateBucketIfNotExists(bucketNameProjects)
		if err != nil {
			return err
		}

		projectBucket, err2 := projects.CreateBucketIfNotExists([]byte(projectName))
		if err2 != nil {
			return err2
		}

		projectFilesBucket, err3 := projectBucket.CreateBucketIfNotExists(bucketNameFiles)
		if err3 != nil {
			return err3
		}

		nextID, _ := projectFilesBucket.NextSequence()
		fileBucket, subErr := projectFilesBucket.CreateBucket([]byte(strconv.Itoa(int(nextID))))
		if subErr != nil {
			return subErr
		}

		subPutErr := fileBucket.Put(bucketKeyFileName, []byte(file))
		if err != nil {
			return subPutErr
		}

		return nil
	})
}

func (c *BoltConfigStorage) GetDockerComposeFilesByProject(projectName string) (dcf.DockerComposeProject, error) {
	result := dcf.DockerComposeProject{}

	db := c.openDB()
	defer closeDB(db)

	err := db.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte(projectName))
		files := project.Bucket(bucketNameFiles)
		c := files.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			if value == nil { // a nil value means it's a bucket
				filesBucket := files.Bucket(key)
				result = append(result, dcf.InitDockerComposeFile(string(filesBucket.Get(bucketKeyFileName))))
			}
		}

		return nil
	})
	if err != nil {
		return dcf.DockerComposeProject{}, err
	}

	return result, nil
}

func (c *BoltConfigStorage) GetExecConfigByProject(projectName string) (dcf.ProjectExecConfig, error) {
	var result dcf.ProjectExecConfig

	db := c.openDB()
	defer closeDB(db)

	err := db.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte(projectName))
		config := project.Bucket([]byte(bucketNameProjectConfig))

		if config == nil {
			return errors.New("no config found")
		}
		name := string(config.Get(bucketKeyContainerName))
		command := string(config.Get(bucketKeyCommand))

		result = dcf.InitProjectExecConfig(name, command)

		return nil
	})

	return result, err
}

func (c *BoltConfigStorage) SaveExecConfig(execConfig dcf.ProjectExecConfigInterface, projectName string) error {
	db := c.openDB()
	defer closeDB(db)

	return db.Update(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		project := projects.Bucket([]byte(projectName))
		config, _ := project.CreateBucketIfNotExists([]byte(bucketNameProjectConfig))

		config.Put(bucketKeyContainerName, []byte(execConfig.GetContainerName()))
		config.Put(bucketKeyCommand, []byte(execConfig.GetCommand()))

		return nil
	})
}

func (c *BoltConfigStorage) GetDockerComposeProjectList(projectNamePrefix string) ([]string, error) {
	db := c.openDB()
	defer closeDB(db)

	var result []string

	err := db.View(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)

		cursor := projects.Cursor()

		prefix := []byte(projectNamePrefix)
		for key, _ := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, _ = cursor.Next() {
			result = append(result, string(key))
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return result, nil
}

func (c *BoltConfigStorage) DeleteProjectByName(name string) error {
	db := c.openDB()
	defer closeDB(db)

	return db.Update(func(tx *bolt.Tx) error {
		projects := tx.Bucket(bucketNameProjects)
		return projects.DeleteBucket([]byte(name))
	})
}

func (c *BoltConfigStorage) openDB() *bolt.DB {
	db, err := bolt.Open(c.filePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func closeDB(db *bolt.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
