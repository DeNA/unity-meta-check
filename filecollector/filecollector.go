package filecollector

import (
	"bufio"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
	"sync"
)

type FileCollector func(projRootAbs typedpath.RawPath, targetRel typedpath.RawPath, opts *Options, writer chan<- typedpath.SlashPath) error

func New(gitLsFiles git.LsFiles, logger logging.Logger) FileCollector {
	return func(projRootAbs typedpath.RawPath, targetRel typedpath.RawPath, opts *Options, writer chan<- typedpath.SlashPath) error {
		targetAbs := projRootAbs.JoinRawPath(targetRel)
		targetRelSlash := targetRel.ToSlash()

		logger.Debug(fmt.Sprintf("searching: %q", targetAbs))
		dirSet := pathutil.NewPathSetWithSize(opts.IgnoreCase, 50000)

		stdoutReader, stdoutWriter, err := os.Pipe()
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			dirSet.Add(targetRelSlash) // XXX: To prevent endless loop.

			scanner := bufio.NewScanner(stdoutReader)
			for scanner.Scan() {
				var filePath typedpath.SlashPath
				if targetRelSlash == "." {
					filePath = typedpath.SlashPath(scanner.Text())
				} else {
					filePath = targetRelSlash.JoinSlashPath(typedpath.SlashPath(scanner.Text()))
				}
				writer <- filePath

				// NOTE: It should be path.Dir instead of filepath.Dir because it will be used as file collections
				//       (file collections are slash separated path for all OSes).
				dirname := filePath.Dir()
				for {
					if dirSet.Has(dirname) {
						break
					}

					writer <- dirname
					dirSet.Add(dirname)
					dirname = dirname.Dir()
				}
			}
		}()

		if err := gitLsFiles(targetAbs, []string{}, stdoutWriter); err != nil {
			return err
		}
		wg.Wait()
		logger.Debug(fmt.Sprintf("length of set for collected path footprints: %d", dirSet.Len()))
		return nil
	}
}
