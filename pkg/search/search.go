package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

func All(ctx context.Context, phrase string, files []string) <-chan []Result {

	ch := make(chan []Result)
	//reads := make([]*os.File, 0)
	//
	wg := sync.WaitGroup{}
	//for _, file := range files {
	//	read, err := os.Open(file)
	//	if err != nil {
	//		log.Print(err)
	//	}
	//
	//	reads = append(reads, read)
	//
	//}

	goroutines := len(files)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(ch chan<- []Result, name string, i int, ctx context.Context) {
			defer wg.Done()
			//buf := make([]byte, 4094)
			//content := make([]byte, 0)
			//
			var results []Result
			//
			//for {
			//	rd, err := reads[i].Read(buf)
			//
			//	if err == io.EOF {
			//		content = append(content, buf[:rd]...)
			//		break
			//	}
			//	if err != nil {
			//		log.Print(err)
			//	}
			//
			//	content = append(content, buf[:rd]...)
			//}
			data, err := ioutil.ReadFile(name)
			if err != nil {
				log.Println(err)
			}

			datas := strings.Split(string(data), "\n")
			for line, str := range datas {
				pos := strings.Index(str, phrase)
				if pos != -1 {
					result := Result{
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
			if len(results) > 0 {
				ch <- results
			}
		}(ch, files[i], i, ctx)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}
