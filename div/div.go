// +build js,wasm

package div

import "github.com/dennwc/dom"

// Writer implements an io.Writer that appends content
// to the dom.Element. It should ideally be used on a <div> element.
type Writer dom.Element

// Write implements io.Writer.
func (d Writer) Write(p []byte) (n int, err error) {
	node := dom.GetDocument().CreateElement("div")
	node.SetInnerHTML(string(p))
	(*dom.Element)(&d).AppendChild(node)
	return len(p), nil
}
