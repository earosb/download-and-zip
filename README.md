# Download and compress
Download multiple files concurrently and get all compress zip file, written in go lang

```
curl http://127.0.0.1:8080?files=https://picsum.photos/200/300 --output downloads.zip
```

## TODO

- Clean downloaded fields after response
- Validate request input
- Download from private bucket (S3, Cloud Storage)
- Compress in other formats 