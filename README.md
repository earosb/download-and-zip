# Download and compress
Download multiple files concurrently from different clouds (aws s3, gcs, do spaces, etc) and get all compress zip file, written in go lang

How to test

```
docker run -it --rm -p 8080:8080 --name download-and-compress earosb/download-and-compress:v1.0.1

curl http://127.0.0.1:8080?files=https://picsum.photos/200/300,https://picsum.photos/200/300 --output downloads.zip
```

## TODO

- Validate request input
- Validate env vars
- Compress in other formats 