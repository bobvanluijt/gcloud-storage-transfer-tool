/**
 * Upload website to Google Cloud Bucket
 *
 * Author: @bobvanluijt
 * Info: https://cloud.google.com/go/home
 */

package main

import (
	"flag"
	"fmt"
	"strings"
	"path/filepath"
	"log"
	"os"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

const (
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

var (
	projectID  = flag.String("project", "", "Your cloud project ID.")
	bucketName = flag.String("bucket", "", "The name of an existing bucket within your project.")
	fileName   = flag.String("file", "", "The file to upload.")
	dirName    = flag.String("dir", "", "The dir to upload.")
)

func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}

func insertFile(service *storage.Service, fileName string) {
	// Insert an object into a bucket.
	object := &storage.Object{Name: fileName}

	file, err := os.Open(fileName)
	if err != nil {
		fatalf(service, "Error opening %q: %v", fileName, err)
	}
	if res, err := service.Objects.Insert(*bucketName, object).Media(file).Do(); err == nil {
		fmt.Println("Created object ", res.Name)
	} else {
		fatalf(service, "Objects.Insert failed: %v", err)
	}
}

func createDir(path string){
	fmt.Println("directory: ", path)
}

func main() {
	flag.Parse()
	if *bucketName == "" {
		log.Fatalf("Bucket argument is required. See --help.")
	}
	if *projectID == "" {
		log.Fatalf("Project argument is required. See --help.")
	}
	if *fileName == "" && *dirName == "" {
		log.Fatalf("You need to add to file or dir argument. See --help.")
	}

	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}

	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	// If the bucket already exists and the user has access, warn the user, but don't try to create it.
	if _, err := service.Buckets.Get(*bucketName).Do(); err == nil {
		fmt.Printf("Bucket %s exists. Use it to add data", *bucketName)
	} else {
		// Create a bucket.
		if res, err := service.Buckets.Insert(*projectID, &storage.Bucket{Name: *bucketName}).Do(); err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Failed creating bucket %s: %v", *bucketName, err)
		}
	}

	// Check if a single dir or file needs to be uploaded
	if *fileName != "" {
		// upload single file
		insertFile(service, *fileName)
	} else if *dirName != "" {
		// upload dir

		filepath.Walk(*dirName, func(path string, fileInfo os.FileInfo, err error) error {
	        //fmt.Printf("Visited: %s\n", path)
			
			// search hidden files, if the path contains /. it contains an hidden entity
			isHidden := strings.Index(path, "/.")
		    if isHidden > -1 || string(path[0]) == "." {

		    	fmt.Println("Hidden files are not uploaded", path);

		    } else {

				f, err := os.Open(path)
				if err != nil {
			        log.Fatalf("Something went wrong looping through dirs %v", err)
			    }

			    defer f.Close()
			    fi, err := f.Stat()
			    if err != nil {
			        log.Fatalf("Something went wrong looping through dirs %v", err)
			    }

			    switch mode := fi.Mode(); {
			    case mode.IsDir():
			        // do directory stuff
			        createDir(path)
			    case mode.IsRegular():
			        // do file stuff
			        insertFile(service, path)
			        //fmt.Println("file: %v", path)
			    }

			}

			return nil
	    })
	}

}