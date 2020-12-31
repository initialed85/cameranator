package page_renderers

import (
	"sync"
	"time"

	"github.com/initialed85/cameranator/pkg/front_end/legacy/page_renderer"
)

type Role struct {
	IsSegment bool
	Path      string
}

type PageRenderers struct {
	mu                 sync.Mutex
	roles              []Role
	pageRendererByRole map[Role]*page_renderer.PageRenderer
	url                string
	timeout            time.Duration
}

func NewPageRenderers(
	url string,
	timeout time.Duration,
	roles []Role,
) *PageRenderers {
	p := PageRenderers{
		roles:              roles,
		pageRendererByRole: make(map[Role]*page_renderer.PageRenderer),
		url:                url,
		timeout:            timeout,
	}

	return &p
}

func (p *PageRenderers) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var err error

	for _, role := range p.roles {
		p.pageRendererByRole[role], err = page_renderer.NewPageRenderer(
			p.url,
			p.timeout,
			role.IsSegment,
			role.Path,
		)
		if err != nil {
			return err
		}

		p.pageRendererByRole[role].Start()
	}

	return nil
}

func (p *PageRenderers) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pageRenderer := range p.pageRendererByRole {
		pageRenderer.Stop()
	}
}
