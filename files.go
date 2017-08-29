package now

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryanuber/go-glob"
)

var defaultIgnorePaths = []string{
	".git",
	".gitmodules",
	".svn",
	".npmignore",
	".dockerignore",
	".gitignore",
	".*.swp",
	".DS_Store",
	".wafpicke-*",
	".lock-wscript",
	"npm-debug.log",
	"config.gypi",
	"node_modules",
	"CVS",
}

// FileInfo represents a deployment file to be uploaded
type FileInfo struct {
	Sha  string `json:"sha"`
	Size int64  `json:"size"`
	File string `json:"file"`
	Mode uint32 `json:"mode"`
}

// FileHashMap represents a map of files and their hashes
type FileHashMap map[string]FileHash

// FileHash represents a file's info, including duplicate names
type FileHash struct {
	Sha   string
	Names []FileInfo
	Path  string
}

// PackageType infers the project's type
func PackageType(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); !os.IsNotExist(err) {
		return "docker"
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); !os.IsNotExist(err) {
		return "npm"
	}
	return ""
}

// StaticFiles returns an array of paths for a given static project
func StaticFiles(dir string) (*[]string, error) {
	// TODO: make compliant with all config and ignore options
	ignoreFiles := append([]string{}, defaultIgnorePaths...)

	// Obey gitignore file if exists
	gitIgnoreFiles, err := readIgnore(dir, ".gitignore")
	if err != nil {
		return nil, err
	}
	ignoreFiles = append(ignoreFiles, gitIgnoreFiles...)

	// Walk the directory of files
	files, err := readDirFiles(dir, ignoreFiles)
	if err != nil {
		return nil, err
	}
	return &files, nil
}

// DockerFiles returns an array of paths for a given Docker project
func DockerFiles(dir string) (*[]string, error) {
	// TODO: make compliant with all config and ignore options
	ignoreFiles := append([]string{}, defaultIgnorePaths...)

	// Obey dockerignore file if exists
	dockerIgnoreFiles, err := readIgnore(dir, ".dockerignore")
	if err != nil {
		return nil, err
	}
	ignoreFiles = append(ignoreFiles, dockerIgnoreFiles...)

	// Obey gitignore file if exists
	gitIgnoreFiles, err := readIgnore(dir, ".gitignore")
	if err != nil {
		return nil, err
	}
	ignoreFiles = append(ignoreFiles, gitIgnoreFiles...)

	// Walk the directory of files
	files, err := readDirFiles(dir, ignoreFiles)
	if err != nil {
		return nil, err
	}
	return &files, nil
}

// NpmFiles returns an array of paths for a given npm package
func NpmFiles(dir string) (*[]string, error) {
	// TODO: make compliant with all config and ignore options
	ignoreFiles := append([]string{}, defaultIgnorePaths...)

	// Obey npmignore file if exists
	npmIgnoreFiles, err := readIgnore(dir, ".npmignore")
	if err != nil {
		return nil, err
	}
	ignoreFiles = append(ignoreFiles, npmIgnoreFiles...)

	// Obey gitignore file if exists
	gitIgnoreFiles, err := readIgnore(dir, ".gitignore")
	if err != nil {
		return nil, err
	}
	ignoreFiles = append(ignoreFiles, gitIgnoreFiles...)

	// Walk the directory of files
	files, err := readDirFiles(dir, ignoreFiles)
	if err != nil {
		return nil, err
	}

	return &files, nil
}

func readDirFiles(dir string, ignoreFiles []string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		// TODO: error if package.json doesn't have a start or now-start script, or server.js
		for _, wildcard := range ignoreFiles {
			if glob.Glob(wildcard, filepath.Base(path)) {
				return nil
			}
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func readIgnore(dir, typ string) ([]string, error) {
	ignorePath := filepath.Join(dir, typ)
	if _, err := os.Stat(ignorePath); !os.IsNotExist(err) {
		ignoreFile, err := ioutil.ReadFile(ignorePath)
		if err != nil {
			return []string{}, err
		}
		return strings.Split(string(ignoreFile), "\n"), nil
	}
	return []string{}, nil
}

// NewFilesList returns an array of FileInfo arrays for the given list of paths
func NewFilesList(wd string, paths []string) (*[]FileInfo, *FileHashMap, error) {
	fhm := make(FileHashMap)
	for _, p := range paths {
		hasher := sha1.New()
		file, err := os.Open(p)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, nil, err
		}
		bs := hasher.Sum(nil)
		sha := hex.EncodeToString(bs)
		fileInfo, err := newFileInfoForSha(wd, sha, p, file)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := fhm[sha]; ok {
			m := fhm[sha]
			m.Names = append(m.Names, *fileInfo)
			fhm[sha] = m
		} else {
			fhm[sha] = FileHash{
				Sha:   sha,
				Names: []FileInfo{*fileInfo},
				Path:  p,
			}
		}
	}
	var files []FileInfo
	for _, f := range fhm {
		files = append(files, f.Names...)
	}
	return &files, &fhm, nil
}

func newFileInfoForSha(wd, sha, filepath string, file *os.File) (*FileInfo, error) {
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Sha:  sha,
		Size: stats.Size(),
		File: strings.Replace(filepath, wd+"/", "", 1),
		Mode: uint32(stats.Mode()) ^ 0100000,
	}, nil
}
