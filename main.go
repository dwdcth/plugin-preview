package preview

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
)

//go:embed ui
var f embed.FS

type PreviewConfig struct {
}

func (p *PreviewConfig) OnEvent(event any) {

}

var _ = InstallPlugin(&PreviewConfig{})

func (p *PreviewConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/preview/" {
		var s string
		Streams.Range(func(streamPath string, _ *Stream) {
			s += fmt.Sprintf("<a href='%s'>%s</a><br>", streamPath, streamPath)
		})
		if s != "" {
			s = "<b>Live Streams</b><br>" + s
		}
		for name, p := range Plugins {
			if pullcfg, ok := p.Config.(config.PullConfig); ok {
				if pullonsub := pullcfg.GetPullConfig().PullOnSub; pullonsub != nil {
					s += fmt.Sprintf("<b>%s pull stream on subscribe</b><br>", name)
					for streamPath, url := range pullonsub {
						s += fmt.Sprintf("<a href='%s'>%s</a> <-- %s<br>", streamPath, streamPath, url)
					}
				}
			}
		}
		w.Write([]byte(s))
		return
	}
	ss := strings.Split(r.URL.Path, "/")
	if b, err := f.ReadFile("ui/" + ss[len(ss)-1]); err == nil {
		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(ss[len(ss)-1])))
		w.Write(b)
	} else {
		//w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		//w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		b, err = f.ReadFile("ui/index.html")
		w.Write(b)
	}
}
