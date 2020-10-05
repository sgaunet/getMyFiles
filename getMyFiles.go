package main

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Calc Sha256 of the content of a file
func calcSha256(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), err
}

func creeFolder(dir string) error {
	var err error
	dirSplit := strings.Split(dir, "/")
	completeDir := ""

	for _, value := range dirSplit {
		if len(completeDir) == 0 {
			completeDir = value
		} else {
			completeDir = fmt.Sprintf("%s/%s", completeDir, value)
		}

		fmt.Println(completeDir)

		if _, err := os.Stat(completeDir); os.IsNotExist(err) {
			errDir := os.MkdirAll(completeDir, 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}
	}
	return err
}

func ensureDirExists(dir1 string) error {
	var err error

	if _, err := os.Stat(dir1); os.IsNotExist(err) {
		err = creeFolder(dir1)
	}
	return err
}

func copyMyFile(pathitem string, origin string, dest string, sha256 string) error {
	var completeDest = pathitem
	completeDest = strings.ReplaceAll(completeDest, origin, dest)
	fmt.Println("Copy from ", pathitem, "to", completeDest)
	err := ensureDirExists(path.Dir(completeDest))

	from, err := os.Open(pathitem)
	if err != nil {
		fmt.Println("Failed to open :", pathitem)
		return err
	}
	defer from.Close()

	// Check if dest exists
	if _, err := os.Stat(completeDest); err == nil {
		// dest exists, check if checksum are correct
		newSha256, err := calcSha256(completeDest)
		if sha256 == newSha256 {
			fmt.Printf("%s already copied\n", pathitem)
			return err
		}
	}

	to, err := os.OpenFile(completeDest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Failed to create :", completeDest)
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		fmt.Println("Failed to copy :", completeDest)
		return err
	}

	newSha256, err := calcSha256(completeDest)
	if sha256 != newSha256 {
		err = errors.New("checksum error : the checksum of the copied file is not the same as original")
	}
	return err
}

func main() {
	var origin, dest string

	flag.StringVar(&origin, "o", ".", "Origin folder")
	flag.StringVar(&dest, "d", ".", "Destination folder")
	flag.Parse()
	// fmt.Println(origin, dest)

	_, err := ioutil.ReadDir(origin)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(origin,
		func(pathitem string, info os.FileInfo, err error) error {
			var sha256 string
			var cpt int

			//fmt.Println(pathitem, info.Size(), info.IsDir(), path.Base(pathitem), path.Dir(pathitem))
			if !info.IsDir() {
				for {
					sha256, _ = calcSha256(pathitem)
					err = copyMyFile(pathitem, origin, dest, sha256)
					if err != nil {
						cpt++
						fmt.Println("Failed to copy, wait 30s, Try #", cpt)
						time.Sleep(30 * time.Second)
					}
					if cpt > 100 {
						break
					}
					if err == nil {
						break
					}
				}
			}
			return err
		})
	if err != nil {
		log.Println(err)
	}
}
