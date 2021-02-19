package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func handleFsError(err *error, message string) {
	if *err != nil {
		fmt.Printf("[Error] %s\n", message)
		os.Exit(0)
	}
}

//FsWatchObject object type to store watcher
type FsWatchObject struct {
	rootWatcher *fsnotify.Watcher
	nCounter    int
	dirMap      map[string][]string
}

//CompilerMetadata compiler metadata that will be passed on to each function
type CompilerMetadata struct {
	compiler     Compiler
	cacheHandler CacheHandler
	configMap    ConfigData
	rootPath     string
}

//Event functions

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		handleFsError(&err, "File "+path+" not exist.")
	}

	return stat.IsDir()
}

func onDirectoryCreate(fswatcher *FsWatchObject, path string, cm *CompilerMetadata) {
	//traverse the tree and with current directory as root and add all the subdirectories
	err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			fmt.Printf("Adding %s to watcher tree.\n", path)
			fswatcher.rootWatcher.Add(subPath)
			fswatcher.nCounter++
		} else {
			path := filepath.Join(cm.rootPath, path, subPath)
			cm.compiler.Compile(&path, &cm.configMap, cm.cacheHandler)
		}

		return err
	})

	handleFsError(&err, "Error adding "+path+" to watcher")

	//call linker to link during directory change
	cm.compiler.Link(&cm.configMap, cm.cacheHandler)

	fmt.Printf("Total watchers %d\n", fswatcher.nCounter)
}

func onFileCreate(path string, cm *CompilerMetadata) {
	fmt.Printf("Add %s detected.\n", path)
	path = filepath.Join(cm.rootPath, path)

	result := cm.compiler.Compile(&path, &cm.configMap, cm.cacheHandler)

	if result {
		cm.compiler.Link(&cm.configMap, cm.cacheHandler)
	}

}

func onFileChanged(path string, cm *CompilerMetadata) {
	path = filepath.Join(cm.rootPath, path)
	fmt.Printf("File %s modification detected.\n", path)
	result := cm.compiler.Compile(&path, &cm.configMap, cm.cacheHandler)

	if result {
		//call the linker
		cm.compiler.Link(&cm.configMap, cm.cacheHandler)
	}
}

//initFsWatcher Initial method to add initial state-directory of the project
func initFsWatcher(rootDir string, compilerMetadata *CompilerMetadata) *FsWatchObject {
	_, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		handleFsError(&err, "File "+rootDir+" not found")
	}

	//traverse the rootdirectory and watch
	watcher := FsWatchObject{}
	watcher.rootWatcher, _ = fsnotify.NewWatcher()

	//walk directory and add ech directory to the watcher list
	err = filepath.Walk(rootDir, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			fmt.Printf("Adding %s to watcher tree\n", subPath)
			watcher.rootWatcher.Add(subPath)
			watcher.nCounter++
		} else {
			path := filepath.Join(rootDir, subPath)
			compilerMetadata.compiler.Compile(&path, &compilerMetadata.configMap, compilerMetadata.cacheHandler)
		}

		return nil
	})

	handleFsError(&err, "Failed to fs watcher event")
	compilerMetadata.compiler.Link(&compilerMetadata.configMap, compilerMetadata.cacheHandler)

	return &watcher
}

func handleFsChanges(watcher *FsWatchObject, compilerMetadata *CompilerMetadata) {
	for {
		select {
		case event := <-watcher.rootWatcher.Events:

			switch eventType := event.Op; eventType {
			case fsnotify.Create:
				if isDir(event.Name) {
					onDirectoryCreate(watcher, event.Name, compilerMetadata)
				} else {
					onFileCreate(event.Name, compilerMetadata)
				}
			case fsnotify.Remove:
				fmt.Printf("Delete %s detected\n", event.Name)
			case fsnotify.Write:
				onFileChanged(event.Name, compilerMetadata)
			}

		case err := <-watcher.rootWatcher.Errors:
			fmt.Println(err)
		}
	}
}

func readConfigFromJSON(jsonPath string) ConfigData {
	_, err := os.Stat(jsonPath)
	if os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist.\n", jsonPath)
		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Faile to read %s.\n", jsonPath)
		os.Exit(0)
	}

	configMap := ConfigData{}

	err = json.Unmarshal(bytes, &configMap)
	if err != nil {
		fmt.Printf("Invalid json config %s provided.", jsonPath)
		os.Exit(0)
	}

	return configMap
}

//InitWatcher : Initialize the watcher and its root function that handles fs changes
func InitWatcher(jsonFile string) {

	configMap := readConfigFromJSON(jsonFile)

	rootPath := configMap.SourceDir

	//init the compiler
	//use your own compiler and object-cache-handler
	compiler := CCompiler{}
	cacheHandler := ObjectCacheHandler{}

	compilerMetadata := CompilerMetadata{}
	compilerMetadata.configMap = configMap
	compilerMetadata.compiler = compiler
	compilerMetadata.cacheHandler = cacheHandler
	compilerMetadata.rootPath = rootPath

	watcher := initFsWatcher(rootPath, &compilerMetadata)
	fmt.Printf("Initialized %d watchers\n", watcher.nCounter)

	defer watcher.rootWatcher.Close()

	done := make(chan bool)

	go handleFsChanges(watcher, &compilerMetadata)

	<-done
}
