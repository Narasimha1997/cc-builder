### cc-builder
Live compilation and linking tool for C/C++ projects written in Go.

#### What this tool does?
This tool compiles your C/C++ projects and also acts like a live-compiler that detects changes in your project and compiles the changes to keep the build upto date. The tool requires a simple configuration file that allows users to include custom compiler and linker flags. The tool only compiles the changed files incrementally and runs the linker to produce the fresh build, thus unmodified files are ignored from compiling again. You can add your own compiler support as well by modifying the source code.

#### How it works?
The tool uses `fsevents` a cross-platform File-System notifications API to monitor the source changes. 
When the tool is initialized, it traverses the project tree and attaches a `watcher` object to each directory and also compiles the initial state of the project. It then spins-up a `Go-routine` which listens to File-system events, the routine then invokes respective functions as side-effect of these events. The project can work with any known C/C++ compiler like `gcc`, `g++`, `clang` etc. To keep things abstract, the actual compiler functions are defined by an interface called `Compiler`, anyone can ship their own logic/compiler by implementing the interface methods. The tool maintains and regularly updates a `build-cache` which is stores all the compiled objects, since the compiled objects are present, the tool only compiles the updated files again and links the updated objects with old ones to produce a fresh executable binary, just like the compiler, you can also write your own caching functionality by implementing methods of `CacheManager` interface.

#### Building the tool
Run the `build.sh` script to build the tool. You can also install the tool globally by adding it to the `$PATH` or copying it to `/usr/bin` or `/usr/local/bin`.

The binary will be installed at `./bin` of your `GOPATH`

#### Running the tool
The project requires a simple config json file as shown below :
```
{
    "sourceDir" : ".",
    "targetExts" : ["cpp", "c", "cc"],
    "ccopts" : "-I./include",
    "linkerOpts" : "",
    "compiler" : "g++",
    "targetBinaryName" : "test"
}
```
Parameters :
1. `sourceDir` : The root of the project which you want to monitor.
2. `targetExts` : An array of extensions to consider for compiling.
3. `ccopts` : Compiler options string, this will be used during compilation.
4. `linkerOpts` : Linker options string, this will be used during linking phase.
5. `compiler` : The compiler to use, pass executable name. Ex : `gcc`, `g++` etc.
6. `targetBinaryName` : Name of the final build along with platform specific extensions.

