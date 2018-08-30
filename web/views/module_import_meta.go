package views

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"govanityimport/web/controllers"
	"govanityimport/zaplog"
	"govanityimport/errorcode"
	"strings"
	"net/http"
	"html/template"
	"github.com/gin-contrib/location"
	"fmt"
)

const (
	defaultCacheSeconds = 5 * 60
)

const moduleMetaTemplateName = "moduleMeta"
const moduleMetaTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{.ImportModulePath}} {{.Vcs}} {{.RepoURL}}">
        <meta name="go-source" content="{{.ImportModulePath}} {{.HomeURL}} {{.DirPattern}} {{.FilePattern}}">
        <meta http-equiv="refresh" content="0; url={{.DocHost}}/{{.ImportModulePath}}">
    </head>
    <body>
        Nothing to see here. Please <a href="{{.DocHost}}/{{.ImportModulePath}}">move along</a>.
    </body>
</html>
`

var htmlRender render.HTMLRender

type moduleMetaRenderData struct {
	ImportModulePath string
	Vcs              string
	RepoURL          string
	HomeURL          string
	DirPattern       string
	FilePattern      string
	DocHost          string
}

func InitRender() error {
	tmpl, err := template.New(moduleMetaTemplateName).Parse(moduleMetaTemplate)
	if err != nil {
		return err
	}
	htmlRender = render.HTMLProduction{
		Template: tmpl,
	}
	return nil
}

func ModuleImportMetaView(c *gin.Context) {
	log := zaplog.GetSugarLogger()

	controller := controllers.GetController()
	importPath := strings.TrimRight(c.Request.URL.Path, "/")
	u := location.Get(c)
	importPath = u.Host + importPath

	metaInfo, err := controller.GetModuleMetaInfo(c.Request.Context(), importPath)
	if err != nil {
		log.Errorw("query module fail", "error", err,
			"module", importPath)
		c.JSON(http.StatusOK, errorcode.ErrServerError)
		return
	}
	data := &moduleMetaRenderData{
		ImportModulePath: metaInfo.ImportInfo.ModuleImportPath,
		Vcs:              metaInfo.ImportInfo.Vcs,
		RepoURL:          metaInfo.ImportInfo.RepoUrl,
		HomeURL:          metaInfo.SourceInfo.HomeUrl,
		DirPattern:       metaInfo.SourceInfo.DirPattern,
		FilePattern:      metaInfo.SourceInfo.FilePattern,
		DocHost:          metaInfo.SourceInfo.DocHost,
	}
	r := htmlRender.Instance(moduleMetaTemplateName, data)
	c.Header("Cache-Control", fmt.Sprintf("max-age=%d", defaultCacheSeconds))
	c.Render(http.StatusOK, r)
}
