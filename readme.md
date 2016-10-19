# Gcloud Storage Transfer Tool

Goal of this tool is to easily sync data (for data storage or website hosting) with the Google Cloud.
_Note: still working on this. Watch the repo to see when it is done_

## Install

1. Make sure you have [Google gcloud installed](https://cloud.google.com/sdk/gcloud/).
2. You can use the file in the `./bin` to run directly or (if you have Go installed) use `go run gcloud-st.go [arguments]`

## Example:

`$ ./gcloud --project=my-project --bucket=test-bucket-abc-123 --dir=./ --gzip=true --public=true --quite=true --watch=true`

## Usage with build version

You can use:

| Argument        | Description                                                             |
|-----------------|-------------------------------------------------------------------------|
| `--project`     | Project id (mandatory)                                                  |
| `--bucket`      | Bucket id (mandatory)                                                   |
| `--dir`         | Set the dir that needs to be uploaded                                   |
| `--file`        | Set the file that needs to be uploaded                                  |
| `--public`      | If set, content will be public (ie. for website)                        |
| `--gzip`        | If set, content will be gziped and content header will be set correctly |
| `--watch`       | If set, updated dir or file will re-upload                              |
| `--quite`       | If set, only errors will be shown                                       |
| `--allowHidden` | If set, hidden files will be uploaded too                               |

_Note: `--dir-` or `--file` is mandatory_

Example: `$ ./gcloud-st --project=projectId --bucket=bucketnNme --dir=./dirLocation --public --gzip`
