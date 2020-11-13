package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//CCompiler : struct impelementing C compiler
type CCompiler struct {
}

//ObjectCacheHandler : a struct impelementing cache manager
type ObjectCacheHandler struct {
}

const cacheDir string = "./builds"

func getFileExtension(file *string) (string, string) {
	extIndex := strings.LastIndex(*file, ".")
	if extIndex == -1 {
		return "", ""
	}

	charArray := []byte(*file)
	return string(charArray[:extIndex]), string(charArray[extIndex+1:])
}

func getFilePathAndName(file *string) (string, string) {
	fileNameIndex := strings.LastIndex(*file, "/")
	if fileNameIndex == -1 {
		return "", *file
	}
	charArray := []byte(*file)
	return string(charArray[:fileNameIndex]), string(charArray[fileNameIndex+1:])
}

//Compile : implementation of C compiler caller
func (cCompiler CCompiler) Compile(file *string, config *ConfigData, cacheHandler CacheHandler) bool {

	var fileName, ext string
	var extFound bool = false

	for _, suffix := range config.TargetSuffixes {
		fileName, ext = getFileExtension(file)
		if ext == "" {
			extFound = false
			break
		}

		if ext == suffix {
			extFound = true
			break
		}
	}

	if !extFound {
		return false
	}

	info, err := os.Stat(cacheDir)

	if os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0700)
	}

	if info.Size() == 0 {
		//ignore a zero-sized file
		return false
	}

	filePath, fileNameWithExt := getFilePathAndName(&fileName)

	objectName := strings.ReplaceAll(filePath, "/", "")
	fileNameWithExt = fileNameWithExt + ".o"

	objectName = objectName + "_" + fileNameWithExt

	objectPath := filepath.Join(cacheDir, objectName)

	objectPath = strings.Trim(objectPath, " ")

	ccOpts := strings.Split(config.Ccopts, " ")

	//execute the compiler
	commandArgs := []string{"-c", "-o", objectPath, *file}
	fmt.Println(commandArgs)

	commandArgs = append(commandArgs, ccOpts...)
	compilerExec := exec.Command(config.Compiler, commandArgs...)

	//execute the command with stdin, stdout and stderr
	compilerExec.Stdout = os.Stdout
	compilerExec.Stderr = os.Stderr

	fmt.Printf("Executing %s\n", config.Compiler+" "+strings.Join(commandArgs, " "))

	err = compilerExec.Run()

	if err != nil {
		return false
	}

	return true
}

//Link : Link binaries and generate an executable file
func (cCompiler CCompiler) Link(config *ConfigData, cacheHandler CacheHandler) bool {

	//get cache files
	objectFiles := cacheHandler.GetCompiledObjects(config)

	outputPath := config.TatgetBinaryName

	linkerArgs := []string{}

	if config.LinkOpts == "" {
		linkerArgs = objectFiles
	} else {
		linkerFlags := strings.Split(config.LinkOpts, " ")
		linkerArgs = append(objectFiles, linkerFlags...)
	}

	linkerArgs = append(linkerArgs, []string{"-o", outputPath}...)
	linkerExec := exec.Command(config.Compiler, linkerArgs...)

	linkerExec.Stdout = os.Stdout
	linkerExec.Stderr = os.Stderr

	fmt.Printf("Executing %s\n", config.Compiler+" "+strings.Join(linkerArgs, " "))

	err := linkerExec.Run()
	if err != nil {
		return false
	}

	return true
}

//GetCompiledObjects : Returns a list of objects present in the cache for linking
func (cache ObjectCacheHandler) GetCompiledObjects(config *ConfigData) []string {
	dir, err := os.Open(cacheDir)
	if err != nil {
		fmt.Println("Linker failed to open the cache dir")
		os.Exit(0)
	}

	objectFiles, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Println("Error getting directory list from cache dir")
		os.Exit(0)
	}

	for idx := 0; idx < len(objectFiles); idx++ {
		objectFiles[idx] = strings.Trim(filepath.Join(cacheDir, objectFiles[idx]), " ")
	}

	return objectFiles
}

//DeleteCompiledObjects Delete one or more object files provided the prefix path
func (cache ObjectCacheHandler) DeleteCompiledObjects(prefixPath string, config *ConfigData) {
	prefixPathString := strings.ReplaceAll(prefixPath, "/", "")
	dir, err := os.Open(cacheDir)
	if err != nil {
		fmt.Println("Linker failed to delete files from cacheDir.")
		os.Exit(0)
	}

	objectFiles, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Println("Linker failed to get cache objects.")
		os.Exit(0)
	}

	for _, objectFile := range objectFiles {
		splits := strings.Split(objectFile, "_")
		if len(splits) > 0 {
			objectPrefix := splits[0]
			if objectPrefix == prefixPathString {
				path := filepath.Join(cacheDir, objectFile)
				err = os.Remove(path)
				if err != nil {
					fmt.Println("Error removing path " + path)
					os.Exit(0)
				}
			}
		}
	}
}
