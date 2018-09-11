package main

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

type Caches struct {
	Blocks map[string]map[string]string `yaml:"cache"`
}


func main() {
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
		for key, value := range caches {
			fmt.Println("folder is " + key)
			fmt.Println("file is " + value)
		}	
	}

	md, err := hash_file_md5("composer.lock")

	if err !=  nil{
		fmt.Println("Oop!")
	} else {
		fmt.Println(md)
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