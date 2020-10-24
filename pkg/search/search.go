package search

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

	root := context.Background()
	ctx, _ = context.WithTimeout(root, time.Second*10)

	for _, file := range files {
		read, err := os.Open(file)

		if err != nil {
			log.Print(err)
		}

		reads = append(reads, read)
	}
	wg := sync.WaitGroup{}

	goroutines := len(files)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		i := i
		go func(ch chan<- []Result, ctx context.Context) {
			defer wg.Done()

			<-ctx.Done()
			buf := make([]byte, 4094)
			content := make([]byte, 0)

			results := make([]Result, 0)

			for {
				read, err := reads[i].Read(buf)

				if err == io.EOF {
					content = append(content, buf[:read]...)
					break
				}
				if err != nil {
					log.Print(err)
				}

				content = append(content, buf[:read]...)
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
					results = append(results, result)
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
			return
		}(ch, ctx)

	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	<-ctx.Done()

	<-time.After(time.Second)

	return ch
}
