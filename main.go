package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"regexp"
	"os/exec"
	"path/filepath"
	"log"
	"fmt"
)

type Caches struct {
	Blocks map[string]map[string]string `yaml:"cache"`
}

const cacheDir = "/cache"

func main() {
	/*
	* Read yml file to get cache configurations
	* Zip file to current folder
	* Copy zipped file to "cache" foler
	*/

    yamlFile, err := ioutil.ReadFile("./framgia-ci.yml")
	if err != nil {
		panic(err)
	}

    var config Caches
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if ags := os.Args[1]; ags == "--create"{
		createCache(config.Blocks)
	} else if ags == "--restore" {
		restoreCache(config.Blocks)
	} else {
		fmt.Println("Invalid argument")
	}

}

func hash_file_md5(filePath string) (string, error) {
	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]

	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func createCache(cacheList map[string]map[string]string) {
	c := make(chan string)
	for  _, caches := range cacheList {
		go func(list map[string]string, c chan string){
			var md5 string
			fileName, hasFile := list["file"]
			folderName, _ := list["folder"]
			re := regexp.MustCompile(`[./]`)
			folder := re.ReplaceAllString(folderName, "_")
			//----------------------------------------------
			if hasFile{
				md5,_ = hash_file_md5(fileName)
			} else {
				md5 = folder
			}
			//----------------------------------------------
			zipFileName := folder + "_" + md5 + ".zip"

			// zip file not yet exist
			if _, err := os.Stat(cacheDir + zipFileName); os.IsNotExist(err) {
				executeCmd("zip --symlinks -r " + folder + "_" + md5 + " " + folderName)
				executeCmd("cp -rf " + zipFileName + " " + cacheDir)
				executeCmd("rm -rf " + zipFileName)
				c <- "Zip file " + zipFileName + " is created"
			}

			DeleteOldCacheFiles(folderName)

		}(caches, c)
	}

	for i := 0; i < len(cacheList); i++ {
		fmt.Println(<-c)
	}
}

func DeleteOldCacheFiles(patten string) {
	files, err := filepath.Glob(cacheDir + "/" + patten + "*.zip")
	if err != nil {
        log.Println(err)
    }
	
	if (len(files) > 4) {
		for num, file := range files {
			if num > (len(files) - 2) {
				removeOldCacheCmd := exec.Command("bash", "-c", "rm " + file)
				removeOldCacheErr := removeOldCacheCmd.Run()
				if removeOldCacheErr != nil {
					log.Println(removeOldCacheErr)
				}
			}
		}
	}
}

func executeCmd(command string) string{
	Cmd := exec.Command("bash", "-c", command)
	cmdErr := Cmd.Run()
	if cmdErr != nil {
		return "Execute " + command + " is down"
	}

	return "Execute " + command + " is ok"
}

func restoreCache(cacheList map[string]map[string]string) {
	c := make(chan string)
	for  _, caches := range cacheList {
		go func(list map[string]string, c chan string) {
			var md5 string
			fileName, hasFile := list["file"]
			folderName, _ := list["folder"]
			re := regexp.MustCompile(`[./]`)
			folder := re.ReplaceAllString(folderName, "_")
			//----------------------------------------------
			if hasFile{
				md5,_ = hash_file_md5(fileName)
			} else {
				md5 = folder
			}
			//----------------------------------------------
			zipFileName := folder + "_" + md5 + ".zip"
			// zip file not yet exist
			if _, err := os.Stat(cacheDir + "/" + zipFileName); err == nil {
				executeCmd("unzip -qo " + cacheDir + "/" + zipFileName + " -d .")
			} else {
				restoreIfZipFileNotExist(folderName)
			}

			c <- "Trying restore cache for " + folder + " !"
		}(caches, c)
	}

	for i := 0; i < len(cacheList); i++ {
		fmt.Println(<-c)
	}
}

func restoreIfZipFileNotExist(patten string) {
	files, err := filepath.Glob(cacheDir + "/" + patten + "*.zip")
	if err != nil {
		log.Println(err)
    }

	if len(files) > 1 {
		index := len(files) -1
		executeCmd("unzip -qo " + files[index] + " -d .")
	}
}