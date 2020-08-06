package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/scylladb/go-set/strset"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	main := NewMain()
	exitStatus := main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}

type metaHistogram struct {
	extensions map[string]uint
}

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		ignore := strset.New(args...)

		var wg sync.WaitGroup
		var mu sync.Mutex
		result := make(map[string]*metaHistogram, 1000000)
		metaPathCh := make(chan typedpath.RawPath, 4)
		var err error

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(metaPathCh)
			if err2 := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if !strings.HasSuffix(path, unity.MetaSuffix) {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(strings.TrimSuffix(path, unity.MetaSuffix)))
				if !ignore.Has(ext) {
					metaPathCh <- typedpath.RawPath(path)
				}
				return nil
			}); err2 != nil {
				err = err2
			}
		}()

		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				buf := &bytes.Buffer{}
				for metaPath := range metaPathCh {
					file, err := os.Open(string(metaPath))
					if err != nil {
						_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
						return
					}

					buf.Reset()
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := scanner.Text()
						if strings.HasPrefix(line, "guid: ") {
							buf.WriteString("guid: _\n")
							continue
						}
						if strings.HasPrefix(line, "timeCreated: ") {
							buf.WriteString("timeCreated: _\n")
							continue
						}
						buf.WriteString(line)
						buf.WriteByte('\n')
					}
					_ = file.Close()
					masked := buf.String()
					ext := filepath.Ext(string(unity.TrimMetaFromRaw(metaPath)))

					mu.Lock()
					hist, ok := result[masked]
					if ok {
						hist.extensions[ext] += 1
					} else {
						result[masked] = &metaHistogram{
							extensions: map[string]uint{ext: 1},
						}
					}
					mu.Unlock()
				}
			}()
		}

		wg.Wait()
		if err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		for masked, hist := range result {
			_, _ = fmt.Fprintf(procInout.Stdout, "%v ==================\n%s\n", hist.extensions, masked)
		}

		return cli.ExitNormal
	}
}
