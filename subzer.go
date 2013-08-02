package main

import (
    "bufio"
    "container/list"
    "fmt"
    flags "github.com/jessevdk/go-flags"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type SrtBlock struct {
    number uint64
    start  time.Duration
    end    time.Duration
    text   []string
}

func NewSrtBlock(number uint64, start, end time.Duration, text []string) *SrtBlock {
    b := SrtBlock{number, start, end, text}
    return &b
}

func (s *SrtBlock) String() string {
    return fmt.Sprintf("%#v", s)
}

func (s *SrtBlock) Text() []string {
    return s.text
}

func (s *SrtBlock) TextAsLine() string {
    return strings.Join(s.Text(), " ")
}

func (s *SrtBlock) Duration() uint64 {
    return ReduceDuration(s.start, s.end)
}

func SrtBlockParse(lines []string) *SrtBlock {
    if len(lines) < 3 {
        panic("Block should have at least 3 lines.")
    }

    number, _ := strconv.ParseUint(lines[0], 10, 64)

    durations := strings.Split(lines[1], " --> ")
    if len(durations) != 2 {
        panic("Block should have two timestamps.")
    }

    start := ParseTimestamp(durations[0])
    end := ParseTimestamp(durations[1])

    text := lines[2:]

    return NewSrtBlock(number, start, end, text)
}

func ReduceDuration(start time.Duration, end time.Duration) uint64 {
    return uint64(end.Seconds()) - uint64(start.Seconds())
}

func ParseTimestamp(timestamp string) time.Duration {
    pieces := regexp.MustCompile(":|,").Split(timestamp, 4)

    hours, _ := strconv.ParseUint(pieces[0], 10, 8)
    minutes, _ := strconv.ParseUint(pieces[1], 10, 8)
    seconds, _ := strconv.ParseUint(pieces[2], 10, 8)
    milliseconds, _ := strconv.ParseUint(pieces[3], 10, 16)

    return time.Duration(hours)*time.Hour +
        time.Duration(minutes)*time.Minute +
        time.Duration(seconds)*time.Second +
        time.Duration(milliseconds)*time.Millisecond
}

func ListToSlice(list *list.List) []string {
    slice := make([]string, list.Len())

    i := 0
    for e := list.Front(); e != nil; e = e.Next() {
        slice[i] = e.Value.(string)
        i++
    }

    return slice
}

func ParseSrtStream(source io.Reader) *list.List {
    reader := bufio.NewScanner(source)
    blocks := list.New()

    for {
        block_lines := list.New()

        for reader.Scan() {
            line := reader.Text()

            if line == "" {
                break
            }

            block_lines.PushBack(line)
        }

        if block_lines.Len() == 0 {
            break
        }

        srt_block := SrtBlockParse(ListToSlice(block_lines))
        blocks.PushBack(srt_block)
    }

    return blocks
}

func WriteFileLine(writer io.Writer, seconds uint64, text string) {
    line := fmt.Sprintf("%06d\t%s\n", seconds, text)
    io.WriteString(writer, line)
}

func ConvertSrtStream(filename string, source io.Reader, writer io.Writer) uint64 {
    blocks := ParseSrtStream(source)

    WriteFileLine(writer, 0, filename)

    var prev_seconds uint64 = 0
    var curr_seconds uint64 = 0

    for e := blocks.Front(); e != nil; e = e.Next() {
        block := e.Value.(*SrtBlock)

        WriteFileLine(writer, prev_seconds, block.TextAsLine())

        curr_seconds = block.Duration()
        prev_seconds += curr_seconds
    }

    return prev_seconds
}

func WriteTotalSeconds(seeker io.Seeker, writer io.Writer, seconds uint64) {
    io.WriteString(writer, fmt.Sprintf("%06d", seconds))
}

func StreamConversion(destination_name string, source io.Reader, writer io.Writer, seeker io.Seeker) {
    seconds := ConvertSrtStream(destination_name, source, writer)

    seeker.Seek(0, 0)
    WriteFileLine(writer, seconds, destination_name)
}

func ProcessFile(source_path string) {
    destination_path := regexp.MustCompile("srt$").ReplaceAllString(source_path, "txt")
    destination_name := filepath.Base(destination_path)

    source, err := os.Open(source_path)
    if err != nil {
        panic("Cannot open source file.")
    }
    defer source.Close()

    destination, err := os.Create(destination_path)
    if err != nil {
        panic("Cannot open source file.")
    }
    defer destination.Close()

    StreamConversion(destination_name, source, destination, destination)
}

func main() {
    var opts struct {
        File      string `short:"i" description:"File input"`
        Directory string `short:"r" description:"Folder input"`
    }

    parser := flags.NewParser(&opts, flags.Default)
    parser.ApplicationName = "subzer"
    parser.Usage = "[OPTIONS]"

    _, err := parser.Parse()

    if err != nil {
        os.Exit(1)
    }

    file := opts.File
    dir := opts.Directory

    if len(file) > 0 && len(dir) > 0 {
        fmt.Println("Error: input a file or a directory, not both.")
        os.Exit(1)
    }

    if len(file) == 0 && len(dir) == 0 {
        parser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

    files := list.New()

    if len(file) > 0 {
        if _, err := os.Stat(file); err != nil {
            fmt.Printf("Error: file \"%s\" does not exist.\n", file)
            os.Exit(1)
        }

        files.PushBack(file)
    }

    if len(dir) > 0 {
        if _, err := os.Stat(dir); err != nil {
            fmt.Printf("Error: directory %s does not exist.\n", dir)
            os.Exit(1)
        }

        err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
            if filepath.Ext(path) == ".srt" {
                files.PushBack(path)
            }

            return err
        })
    }

    for e := files.Front(); e != nil; e = e.Next() {
        ProcessFile(e.Value.(string))
    }
}
