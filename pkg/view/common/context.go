package common

// AgentContext contains agent-level resources shared across all displays
type AgentContext struct {
	ImageDir    string
	CacheDir    string
	FontsDir    string
	DataSources map[string]interface{}
}

type DisplayContext struct {
	MatrixRows        int
	MatrixCols        int
	DefaultImageSizeX int
	DefaultImageSizeY int
	DefaultFontSize   int
	DefaultFontColor  string
	DefaultFontStyle  string
	DefaultFontType   string
}
