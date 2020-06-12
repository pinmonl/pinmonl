package pkgrepo

const (
	Undefined  = Language("")
	Javascript = Language("javascript")
	PHP        = Language("php")
	Go         = Language("go")
)

var (
	// Readme.
	Readme = FilePath("readme.md")

	// Files of NPM.
	NpmPackage = FilePath("package.json")
)

type Language string

type Langer interface {
	Language() Language
}

type FilePath string

func (f FilePath) String() string {
	return string(f)
}
