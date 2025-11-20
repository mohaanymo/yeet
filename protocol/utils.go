package protocol

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func IsPathExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err!=nil {
		if err == os.ErrNotExist {
			return false, nil
		}

		return false, err
	}

	return true, nil

}

func IsDir(thing string) (bool, error) {
	file, err := os.Stat(thing)
	if err!=nil {
		return false, err
	}
	return file.IsDir(), nil
}

func ExtractFilesFromDirs(inputs []string) []string {
	var results []string
	for _, thing := range inputs {
		if ok, err := IsDir(thing); ok==true && err==nil {
			results = append(results, GetDirFiles(thing)...) 
		}else if ok==false && err==nil{
			results = append(results, thing)
		}else{
			panic(err)
		}
	}
	return results
}

func GetDirFiles(dir string) []string {

	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip directories
		if d.IsDir() {
			return nil
		}

		// append full file path
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}
