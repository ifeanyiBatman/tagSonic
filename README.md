# tagSonic 🎵

tagSonic is a powerful, automated music tagging CLI tool built in Go. It scans your local audio directories, identifies songs using acoustic fingerprinting, and automatically embeds high-quality metadata and beautiful cover art directly into your MP3 files.

## Purpose

Tired of having a music library full of "Track 01" or missing album art? tagSonic eliminates the need to manually tag your downloaded music. 

Unlike simple text-based taggers, tagSonic uses **mathematical audio fingerprinting** to "listen" to your songs. It can identify audio tracks even if the current filename or internal ID3 tags are completely wrong or missing.

### Key Features
- 🔍 **Acoustic Fingerprinting:** Uses `fpcalc` and the AcoustID database to identify the actual audio signature of a track.
- 💿 **Multi-Source Metadata:** Fetches rich metadata primarily from the **iTunes Search API** (for commercial high-resolution data and 1000x1000px cover art), with an integrated fallback to the **MusicBrainz** database.
- 🧠 **Smart Confidence Scoring:** Evaluates the audio fingerprint score against the existing filename and previous ID3 tags. It implements strict active penalties to prevent false positives (tagging a track as the wrong song).
- 🖼️ **Cover Art Embedding:** Automatically downloads and injects high-resolution cover art into the MP3 via `id3v2` Attached Picture frames.
- 📝 **Detailed Logging:** Generates a clean `tagSonic_log.txt` report after every run, detailing exactly which songs were successfully updated and providing reasons for any tracks that were skipped or failed.

---

## Quick Start with Docker (Recommended)

The easiest way to use tagSonic. No need to install Go, Chromaprint, or anything else — just Docker.

```bash
docker run --rm -v /path/to/your/music:/music ifeanyibatman/tagsonic
```

Replace `/path/to/your/music` with the folder containing your audio files:

```bash
# Linux / macOS
docker run --rm -v ~/Music/untagged:/music ifeanyibatman/tagsonic

# Windows
docker run --rm -v "C:\Users\You\Music\untagged":/music ifeanyibatman/tagsonic
```

That's it! tagSonic will process every supported audio file and write the updated tags directly back to your files.

---

## Modifying and Building the Docker Image Yourself

If you want to modify tagSonic or build the image from source:

1. Clone the repo:
   ```bash
   git clone https://github.com/ifeanyiBatman/tagSonic.git
   cd tagSonic
   ```

2. Set up your API key — copy the example env file and fill in your key:
   ```bash
   cp .exampleenv .env
   ```
   Then edit `.env` and replace the placeholder with your [AcoustID API key](https://acoustid.org/login).

3. Build the binary and Docker image:
   ```bash
   # First compile the Linux binary
   CGO_ENABLED=0 GOOS=linux go build -o tagSonic main.go

   # Then build the Docker image
   docker build -t tagsonic .
   docker run --rm -v /path/to/your/music:/music tagsonic
   ```

---

## Manual Setup (Without Docker)

If you'd prefer to run tagSonic natively on your machine.

### Prerequisites

1. **Go:** Download and install from [golang.org](https://go.dev/dl/).
2. **fpcalc (Chromaprint):** The audio fingerprinting calculator required by AcoustID.
   - **Linux (Arch/EndeavourOS):** `sudo pacman -S chromaprint`
   - **Linux (Ubuntu/Debian):** `sudo apt-get install libchromaprint-tools`
   - **macOS:** `brew install chromaprint`
   - **Windows:** Download from [AcoustID.org](https://acoustid.org/chromaprint) and ensure it is added to your System PATH.
3. **AcoustID API Key:** Get a free developer API key by signing in at [acoustid.org](https://acoustid.org/login).

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ifeanyiBatman/tagSonic.git
   cd tagSonic
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Configure your API key:
   ```bash
   cp .exampleenv .env
   ```
   Then edit `.env` and add your AcoustID API key.

### Usage

```bash
# Run directly with Go
go run main.go /path/to/your/music/folder

# Or build a compiled binary
go build -o tagSonic ./...
./tagSonic ~/Music
```

If no directory is provided, tagSonic defaults to `./audios`.

### Reviewing the Results

Once the program finishes, open the generated `tagSonic_log.txt` file to see a detailed summary of successful tracks and explanations for any skipped files.
