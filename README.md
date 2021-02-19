### cc-builder
Live compilation and linking tool for C/C++ projects written in Go.

#### What the tool does?
This tool compiles your C/C++ projects and also acts like a live-compiler that detects changes in your project and compiles the changes to keep the build upto date. The tool requires a simple configuration file that allows users to include custom compiler and linker flags. The tool only compiles the changed files incrementally and runs the linker to produce the fresh build, thus unmodified files are ignored from compiling again. You can add your own compiler support as well by modifying the source code.

### Features:
1. Live compilation of C/C++ projects with less configuration.
2. Automatically takes care of compiling and linking, you don't have to write compile commands.
3. Tracks changes in source tree and compiles them automatically.
4. Supports multiple compilers - `gcc`, `g++`, `musl-c` and `llvm-clang`.
5. Ability to include custom Compiler and Linker options.
6. Simple build cache - Caches most of the build objects, so you spend compute on compiling only changed objects.
7. Support for adding any custom linker (suitable for Operating System Kernel development)

### Requirements:
1. Linux Distribution
2. Go programming language (latest version which supports go-modules)
3. 

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
7. `customLinkerCmd` : Linker which has to be used, by default uses the compiler toolchain specified in `compiler`.

Then you can run the tool as follows :
```
ccbuilder ./myconfig.json
```
The tool will then start and monitors your project to provide hot code reload functionality.
