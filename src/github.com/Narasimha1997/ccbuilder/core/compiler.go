package core

//Compiler : Abstract representation of the compiler functions
type Compiler interface {
	Compile(file *string, config *ConfigData, cacheHandler CacheHandler) bool
	Link(config *ConfigData, cacheHandler CacheHandler) bool
}

//CacheHandler : Abstract representation of the cache functions
type CacheHandler interface {
	GetCompiledObjects(config *ConfigData) []string
	DeleteCompiledObjects(prefixPath string, config *ConfigData)
}
