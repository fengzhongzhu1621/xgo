package uploader

import (
	"fmt"
	"sync"
)

type Progress struct {
	TotalChunks int
	Uploaded    int
	mu          sync.Mutex
}

func (p *Progress) Increment() {
	p.mu.Lock()
	p.Uploaded++
	p.mu.Unlock()
}

func (p *Progress) PrintProgress() {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Printf(
		"\rUploaded: %d/%d (%.2f%%)",
		p.Uploaded,
		p.TotalChunks,
		float64(p.Uploaded)/float64(p.TotalChunks)*100,
	)
}
