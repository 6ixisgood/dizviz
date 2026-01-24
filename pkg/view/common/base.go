package common

import (
	"image"
	"sync"
	"time"

	compCommon "github.com/6ixisgood/matrix-ticker/pkg/component/common"
)

type BaseView struct {
	// Context fields - populated from agent config via context tags
	MatrixRows        int                    `context:"matrix_rows"`
	MatrixCols        int                    `context:"matrix_cols"`
	DefaultFontSize   int                    `context:"default_font_size"`
	DefaultFontColor  string                 `context:"default_font_color"`
	DefaultFontStyle  string                 `context:"default_font_style"`
	DefaultFontType   string                 `context:"default_font_type"`
	DefaultImageSizeX int                    `context:"default_image_size_x"`
	DefaultImageSizeY int                    `context:"default_image_size_y"`
	ImageDir          string                 `context:"image_dir"`
	CacheDir          string                 `context:"cache_dir"`
	FontsDir          string                 `context:"fonts_dir"`
	DataSources       map[string]interface{} `context:"data_sources"`

	// Internal fields
	template        *compCommon.Template
	dataRefresh     *time.Ticker
	templateRefresh *time.Ticker
	stopChan        chan struct{}
	mu              sync.RWMutex
}

func (v *BaseView) Init(configJSON string, ctx ViewContext) error {
	return InitViewFromJSON(v, configJSON, ctx)
}

func (v *BaseView) Template() *compCommon.Template {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.template
}

func (v *BaseView) SetTemplateValue(t compCommon.Template) {
	v.mu.Lock()
	defer v.mu.Unlock()
	*v.template = t
}

func (v *BaseView) SetTemplate(t *compCommon.Template) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.template = t
}

func (v *BaseView) RenderTemplate() image.Image {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.template.Render()
}

func (v *BaseView) TemplateData() map[string]interface{} {
	return map[string]interface{}{}
}

func (v *BaseView) TemplateString() string {
	return ""
}

func (v *BaseView) Stop() {}

func (v *BaseView) getBaseView() *BaseView {
	return v
}
