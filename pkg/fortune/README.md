# Fortune Package

The `fortune` package provides functionality for managing and retrieving random fortune messages from text files. It supports weighted random selection to ensure fair representation of fortunes across files of different sizes.

## Overview

This package allows you to:
- Store fortune messages in simple text files
- Retrieve random fortunes with weighted selection
- Format fortunes as Telegram-compatible markdown
- Support different fortune types (plain text and quotes)

## How It Works

### Fortune File Structure

Fortune files are plain text files stored in the `files/fortune/` directory. Each file should:
- Have a `.txt` extension
- Contain one or more fortune entries
- Separate entries with the delimiter `\n%\n` (newline, percent sign, newline)

Example fortune file (`wisdom.txt`):
```
The journey of a thousand miles begins with a single step.
%
To be yourself in a world that is constantly trying to make you something else is the greatest accomplishment.
%
In the middle of difficulty lies opportunity.
```

### Fortune Types

The package supports two types of fortunes:

#### 1. Text Fortunes
Simple text entries that can span multiple lines:
```
This is a simple fortune.
It can have multiple lines.
%
Another fortune here.
```

#### 2. Quote Fortunes
Quotes with source attribution, formatted as:
```
Quote text goes here.
Multiple lines are supported.

-- Author Name or Source
```

The package automatically detects the type based on the pattern `\n\n-- ` for quotes.

### Dialog Format
Within any fortune, you can use dialog format for conversations:
```
Alice: Hello there!
Bob: Hi, how are you?
Alice: Doing great!
```

Lines matching the pattern `Speaker: Message` will be formatted with emphasis on the speaker name.

## Weighted Random Selection

The `GetRandomFortune()` function implements weighted random selection based on the number of entries in each file. This ensures that:
- Files with more entries have proportionally higher representation
- Small files don't get over-represented
- Large files don't get under-represented

For example, if you have:
- `small.txt` with 2 fortunes
- `large.txt` with 8 fortunes

When calling `GetRandomFortune()`, entries from `large.txt` will be selected 4 times more often than entries from `small.txt`, ensuring fair distribution.

## Usage

### Basic Usage

```go
import "github.com/br0-space/bot/pkg/fortune"

// Create a fortune service
service := fortune.MakeService()

// Get a random fortune from all files
fortune, err := service.GetRandomFortune()
if err != nil {
    log.Fatal(err)
}

// Get the fortune as markdown text
markdown := fortune.ToMarkdown()

// Get the source file name
fileName := fortune.File()
```

### Get Fortune from Specific File

```go
// Get a random fortune from a specific file
fortune, err := service.GetFortune("wisdom")
if err != nil {
    log.Fatal(err)
}
```

### List Available Files

```go
// Get list of all fortune files
files := service.GetList()
for _, file := range files {
    fmt.Println(file)
}
```

### Check if File Exists

```go
// Check if a fortune file exists
exists := service.Exists("wisdom")
if exists {
    fmt.Println("File exists!")
}
```

## Deploying Your Own Fortune Files

To deploy your own fortune files for the bot:

1. **Create Your Fortune File**
   - Create a new `.txt` file in the `files/fortune/` directory
   - Name it descriptively (e.g., `programming.txt`, `movies.txt`)
   - Add your fortunes, separated by `\n%\n`

2. **File Naming**
   - Use lowercase names
   - Use hyphens for multi-word names (e.g., `star-trek.txt`)
   - Avoid special characters

3. **Content Guidelines**
   - Keep individual fortunes concise (1-10 lines typically)
   - Use the quote format (`\n\n-- Author`) for attributed quotes
   - Test that your fortunes display correctly with `service.GetFortune("yourfile")`

4. **Example Structure**
   ```
   files/fortune/
   ├── wisdom.txt
   ├── programming.txt
   ├── movies.txt
   └── jokes.txt
   ```

## Markdown Formatting

Fortunes are automatically formatted for Telegram markdown v2:
- Special characters are escaped
- Dialog lines (`Speaker: Text`) are formatted with bold speaker names
- Quotes include italicized source attribution
- Multi-line fortunes preserve line breaks

## Error Handling

The package handles several error conditions:
- Returns error if no fortune files are found
- Returns error if specified file doesn't exist
- Returns error if file cannot be read
- Skips unreadable files during weighted random selection

## Testing

The package includes comprehensive tests covering:
- Service creation and file listing
- Fortune file reading and parsing
- Random selection (both weighted and per-file)
- Fortune type detection (text vs. quote)
- Markdown formatting and escaping
- Dialog format detection

Run tests with:
```bash
go test ./pkg/fortune/...
```

## Dependencies

- `github.com/br0-space/bot/interfaces` - For the FortuneInterface
- `github.com/br0-space/bot-telegramclient` - For markdown escaping

## Package Structure

- `service.go` - Main service with file management and selection logic
- `fortune.go` - Fortune type and markdown formatting
- `type.go` - Fortune type detection and parsing
- `*_test.go` - Comprehensive test suites
