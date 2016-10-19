/**
 * Upload website to Google Cloud Bucket
 *
 * Author: @bobvanluijt
 * Info: https://cloud.google.com/go/home
 */

package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

/**
 * Set consts
 */
const (
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

/**
 * Define variables
 */
var (
	projectID   = flag.String("project", "", "Your cloud project ID.")
	bucketName  = flag.String("bucket", "", "The name of an existing bucket within your project.")
	fileName    = flag.String("file", "", "The file to upload.")
	dirName     = flag.String("dir", "", "The dir to upload.")
	public      = flag.String("public", "", "Make content public.")
	metaGzip    = flag.String("gzip", "", "Gzip the content and set the meta data to gzip")
	quite       = flag.String("quite", "", "Shows debug information")
	allowHidden = flag.String("allowHidden", "", "Allow hidden files to be uploaded")
)

/**
 * Function for fatal logs
 */
func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}

/**
 * Show debug info
 */
func showDebugInfo(info string) {
	if *quite != "true" {
		fmt.Println("DEBUG", info)
	}
}

/**
 * Find the contenttype of a file
 */
func contentTypeFinder(fileName string) string {
	// collect and set metadata of file
	fileMeta, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fileMeta.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = fileMeta.Read(buffer)
	if err != nil {
		panic(err)
	}

	// Reset the read pointer if necessary.
	fileMeta.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType
}

/**
 * Insert a file in gcloud storage bucket
 */
func insertFile(service *storage.Service, fileName string) {

	var fileToUpload io.Reader
	var object *storage.Object

	file, err := os.Open(fileName)
	if err != nil {
		fatalf(service, "Error opening %q: %v", fileName, err)
	}
	defer file.Close()

	// Check if gzip should be applied, off by default
	if *metaGzip == "true" {
		var byteBuffer = &bytes.Buffer{}
		w := gzip.NewWriter(byteBuffer)
		if _, err := io.Copy(w, file); err != nil {
			fatalf(service, "Error Copy %v", err)
		}
		if err := w.Flush(); err != nil {
			fatalf(service, "Error Flush %v", err)
		}
		if err := w.Close(); err != nil {
			fatalf(service, "Error Close %v", err)
		}
		fileToUpload = byteBuffer

		// set the object with gzip encoding
		object = &storage.Object{Name: fileName, ContentEncoding: "gzip", ContentType: contentTypeFinder(fileName)}

	} else {
		fileToUpload = file

		// set the object without encoding, the filetype will be set automatically
		object = &storage.Object{Name: fileName}
	}

	// Check if file should be public, private by default
	aclType := "private"
	if *public == "true" {
		aclType = "publicRead"
	}

	// Execute the upload
	if res, err := service.Objects.Insert(*bucketName, object).Media(fileToUpload).PredefinedAcl(aclType).Do(); err == nil {
		showDebugInfo("Created object " + res.Name)
	} else {
		fatalf(service, "Objects.Insert failed: %v", err)
	}
}

/**
 * Create a directory
 */
func createDir(path string) {
	showDebugInfo("directory: " + path)
}

/**
 * Main function
 */
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
		showDebugInfo("Bucket %s exists. Use it to add data" + *bucketName)
	} else {
		// Create a bucket.
		if res, err := service.Buckets.Insert(*projectID, &storage.Bucket{Name: *bucketName}).Do(); err == nil {
			showDebugInfo("Created bucket " + res.Name + " at location %v\n\n" + res.SelfLink)
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

			// search hidden files, if the path contains /. it contains an hidden entity
			isHidden := strings.Index(path, "/.")

			// if file is hidden, don't upload it
			if isHidden > -1 || string(path[0]) == "." && *allowHidden != "true" {

				showDebugInfo("Hidden files are not uploaded" + path)

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
				}

			}
			return nil
		})
	}
}
