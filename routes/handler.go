package routes

import (
	"github.com/alexedwards/flow"
	"icyphox.sh/legit/config"
)

func Handlers(c *config.Config) *flow.Mux {
	mux := flow.New()
	d := deps{c}
	mux.HandleFunc("/:name", d.RepoIndex, "GET")
	mux.HandleFunc("/:name/tree/:ref/...", d.RepoTree, "GET")
	mux.HandleFunc("/:name/blob/:ref/...", d.FileContent, "GET")
	return mux
}
