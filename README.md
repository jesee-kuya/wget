# Wget Clone in Go

## Overview

This project is a Go implementation of a subset of the `wget` utility, designed to download files and mirror websites from the web. It supports HTTP, HTTPS, and asynchronous downloading, with features like rate limiting, background downloads, and offline link conversion. The program provides detailed feedback, including progress bars, timestamps, and status updates, making it user-friendly and informative.

## Objectives

The goal is to recreate core functionalities of the `wget` command-line tool using Go, including:

- Downloading a file from a given URL (e.g., `wget https://example.com/file.zip`).
- Saving files with custom names.
- Saving files to specific directories.
- Limiting download speed.
- Downloading in the background with logs redirected to a file.
- Asynchronously downloading multiple files listed in a file.
- Mirroring an entire website, with options to reject file types, exclude directories, and convert links for offline viewing.

## Introduction

`Wget` is a free, non-interactive utility for downloading files from the web, supporting HTTP, HTTPS, and FTP protocols. This project mimics its core features in Go, providing a lightweight, compiled alternative with modern concurrency patterns. For more details on `wget`, refer to the [official manual](https://www.gnu.org/software/wget/manual/wget.html) or run `man wget`.

## Features

### Basic Usage
Download a file from a URL:
```bash
go run . https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
```
**Output Example:**
```
start at 2017-10-14 03:46:06
sending request, awaiting response... status 200 OK
content size: 56370 [~0.06MB]
saving file to: ./EMtmPFLWkAA8CIS.jpg
55.05 KiB / 55.05 KiB [============================] 100.00% 1.24 MiB/s 0s
Downloaded [https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg]
finished at 2017-10-14 03:46:07
```

### Flags
- **`-B`**: Download in the background, redirecting output to `wget-log`.
  ```bash
  go run . -B https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
  ```
  **Output:**
  ```
  Output will be written to "wget-log".
  ```
  **Log File (`wget-log`):**
  ```
  start at 2017-10-14 03:46:06
  sending request, awaiting response... status 200 OK
  content size: 56370 [~0.06MB]
  saving file to: ./EMtmPFLWkAA8CIS.jpg
  Downloaded [https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg]
  finished at 2017-10-14 03:46:07
  ```

- **`-O <filename>`**: Save the file with a custom name.
  ```bash
  go run . -O meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
  ```
  **Output:**
  ```
  saving file to: ./ m eme.jpg
  ...
  ```

- **`-P <directory>`**: Save the file to a specific directory.
  ```bash
  go run . -P ~/Downloads -O meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
  ```
  **Output:**
  ```
  saving file to: ~/Downloads/meme.jpg
  ...
  ```

- **`--rate-limit=<value>`**: Limit download speed (e.g., `400k` for 400 KiB/s, `2M` for 2 MiB/s).
  ```bash
  go run . --rate-limit=400k https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
  ```

- **`-i <file>`**: Download multiple files asynchronously from a file containing URLs.
  ```bash
  cat download.txt
  https://assets.01-edu.org/wgetDataSamples/20MB.zip
  https://assets.01-edu.org/wgetDataSamples/Image_10MB.zip
  go run . -i download.txt
  ```
  **Output:**
  ```
  start at 2025-06-09 13:59:05
  sending request, awaiting response... status 200 OK
  content size: 20971520 [~20.00MB]
  saving file to: ./20MB.zip
  ...
  start at 2025-06-09 13:59:05
  sending request, awaiting response... status 200 OK
  content size: 10485760 [~10.00MB]
  saving file to: ./Image_10MB.zip
  ...
  Download finished: [https://assets.01-edu.org/wgetDataSamples/20MB.zip https://assets.01-edu.org/wgetDataSamples/Image_10MB.zip]
  ```

- **`--mirror`**: Mirror an entire website, saving files in a directory named after the domain.
  ```bash
  go run . --mirror https://example.com
  ```
  **Directory Structure:**
  ```
  example.com/
  ├── index.html
  ├── static/
  │   ├── style.css
  │   └── image.jpg
  ...
  ```

- **Mirror-Specific Flags**:
  - **`-R, --reject=<suffixes>`**: Skip files with specified suffixes.
    ```bash
    go run . --mirror -R=jpg,gif https://example.com
    ```
  - **`-X, --exclude=<paths>`**: Skip specified directories.
    ```bash
    go run . --mirror -X=/assets,/css https://example.com
    ```
  - **`--convert-links`**: Rewrite links in HTML/CSS for offline viewing.
    ```bash
    go run . --mirror --convert-links https://example.com
    ```

## Installation

1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd wget
   ```

2. **Initialize Go Module** (if not already done):
   ```bash
   go mod init github.com/jesee-kuya/wget
   ```

3. **Install Dependencies**:
   ```bash
   go get golang.org/x/net/html
   ```

4. **Build and Run**:
   ```bash
   go build
   ./wget [flags] <URL>
   ```
   Or run directly:
   ```bash
   go run . [flags] <URL>
   ```

## Project Structure

```
wget/
├── downloader/
│   ├── downloader.go      # Core download logic with rate limiting
│   ├── inputDownloader.go # Handles multiple URLs from a file
│   ├── mirror.go          # Website mirroring functionality
│   ├── options.go         # Configuration struct for flags
├── logger/
│   └── logger.go          # Logging with progress bars and status updates
├── parser/
│   ├── parser.go          # HTML/CSS link extraction
│   └── reference.go       # Link rewriting for offline viewing
├── util/
│   └── util.go            # Utility functions (e.g., ContentSize, FormatSpeed)
├── worker/
│   └── worker.go          # Flag parsing and execution logic
└── main.go                # Entry point
```

## Usage Examples

### Download a Single File
```bash
go run . https://example.com/file.zip
```

### Download with Custom Name and Directory
```bash
go run . -P downloads -O custom.zip https://example.com/file.zip
```

### Background Download
```bash
go run . -B https://example.com/file.zip
cat wget-log
```

### Rate-Limited Download
```bash
go run . --rate-limit=200k https://example.com/largefile.zip
```

### Download Multiple Files
```bash
echo -e "https://example.com/file1.zip\nhttps://example.com/file2.zip" > urls.txt
go run . -i urls.txt
```

### Mirror a Website
```bash
go run . --mirror https://example.com
```

### Mirror with Filters
```bash
go run . --mirror -R=jpg,png -X=/assets,/js --convert-links https://example.com
```

## Implementation Details

- **Concurrency**: Uses Go goroutines for asynchronous downloads (`-i` flag), with `sync.WaitGroup` and `sync.Mutex` for safe coordination.
- **Progress Bar**: Displays KiB/MiB downloaded, percentage, speed, and ETA, updating every 500ms. For unknown content lengths, a simplified bar is shown.
- **Rate Limiting**: Implements byte-by-byte throttling in `downloader.go` using a ticker to enforce speed limits.
- **Mirroring**: Recursively crawls websites using `parser.ExtractLinks`, downloading HTML, CSS, and assets. Supports filtering (`-R`, `-X`) and offline link conversion.
- **Logging**: Centralized in `logger.go`, outputs to `os.Stdout` or `wget-log` for background mode.
- **Error Handling**: Logs errors without halting other downloads, ensuring robustness.

## Contributing

Contributions are welcome! Please:
1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/fooBar`).
3. Commit changes (`git commit -am 'Add fooBar feature'`).
4. Push to the branch (`git push origin feature/fooBar`).
5. Open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.