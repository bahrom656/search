package search

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
)

type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

func All(ctx context.Context, phrase string, files []string) <-chan []Result {

	ch := make(chan []Result)
	reads := make([]*os.File, 0)

	var read *os.File
	for _, file := range files {
		read, err := os.Open(file)
		if err != nil {
			log.Print(err)
		}

		reads = append(reads, read)

	}

	goroutines := len(files)

	for i := 0; i < goroutines; i++ {
		i := i
		go func(ch chan<- []Result, ctx context.Context) {
			buf := make([]byte, 4094)
			content := make([]byte, 0)

			results := make([]Result, 0)

			for {
				rd, err := reads[i].Read(buf)

				if err == io.EOF {
					content = append(content, buf[:rd]...)
					break
				}
				if err != nil {
					log.Print(err)
				}

				content = append(content, buf[:rd]...)
			}
			data := string(content)

			datas := strings.Split(data, string(rune(13))+string(rune(10)))
			for line, str := range datas {
				pos := strings.Index(str, phrase)
				if pos != -1 {
					result := Result{}
					result = Result{
						Phrase:  phrase,
						Line:    str,
						LineNum: int64(line + 1),
						ColNum:  int64(pos + 1),
					}
					if	result.ColNum != -1{
						results = append(results, result)
					}
				} //} else {
			}
			//	result = Result{
			//		Phrase:  phrase,
			//		ColNum:  0,
			//		LineNum: 0,
			//		Line:    ""}
			//	results = append(results, result)
			//}
			ch <- results

		}(ch, ctx)

	}
	go func() {
		defer read.Close()
		return
	}()
	return ch
}
