# Gcloud Storage Transfer Tool

Goal of this tool is to easily sync data (for data storage or website hosting) with the Google Cloud.
_Note: still working on this. Watch the repo to see when it is done_

## Install

### Install the compiled version
...

### Use Go directly
...

## Usage with build version

You can use:

|-----------------|-------------------------------------------------------------------------|
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
|-----------------|-------------------------------------------------------------------------|

_Note: `--dir-` or `--file` is mandatory_

Example: `$ ./gcloud-st --project=projectId --bucket=bucketnNme --dir=./dirLocation --public --gzip`