package common

import "github.com/6ixisgood/matrix-ticker/pkg/store"

type AgentContext struct {
	ImageDir    string
	CacheDir    string
	FontsDir    string
	Store       *store.Store
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

type ViewContext struct {
	Agent   *AgentContext
	Display *DisplayContext
}
