package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"os/exec"
	"path/filepath"
	"log"
)

type Caches struct {
	Blocks map[string]map[string]string `yaml:"cache"`
}


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
  
	for  _, caches := range config.Blocks {
		var md5 string
		fileName, hasFile := caches["file"]
		folderName, _ := caches["folder"]
		folder := strings.Replace(folderName, "/", "_", -1)
		//----------------------------------------------
		if hasFile{
			md5,_ = hash_file_md5(fileName)
		} else {
			md5 = folder
		}
		//----------------------------------------------
		zipFileName := folder + "_" + md5 + ".zip"

		// zip file not yet exist
		if _, err := os.Stat("./cache" + zipFileName); os.IsNotExist(err) {
			executeCmd("zip --symlinks -r " + folder + "_" + md5 + " " + folderName)
			executeCmd("cp -rf " + zipFileName + " ./cache")
			executeCmd("rm -rf " + zipFileName)
		}
		DeleteOldCacheFiles(folderName)
		//----------------------------------------------
		/* for key, value := range caches {
			fmt.Println("folder is " + key)
			fmt.Println("file is " + value)
		} */
	}

	//md, err := hash_file_md5("composer.lock")

	/* if err !=  nil{
		fmt.Println("Oop!")
	} else {
		fmt.Println(md)
	} */
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


func DeleteOldCacheFiles(patten string) {
	files, err := filepath.Glob("./cache/" + patten + "*.zip")
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

func executeCmd(command string){
	Cmd := exec.Command("bash", "-c", command)
	cmdErr := Cmd.Run()
	if cmdErr != nil {
		log.Println(cmdErr)
	}
}