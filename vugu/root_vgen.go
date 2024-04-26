package main

// Code generated by vugu via vugugen. Please regenerate instead of editing or add additional code in a separate file. DO NOT EDIT.

import "fmt"
import "reflect"
import "github.com/vugu/vjson"
import "github.com/vugu/vugu"
import js "github.com/vugu/vugu/js"

import "encoding/json"
import "net/http"
import "log"

type Root struct {
	bpi		bpi
	isLoading	bool
}

type bpi struct {
	Time	struct {
		Updated string `json:"updated"`
	}	`json:"time"`
	BPI	map[string]struct {
		Code		string	`json:"code"`
		Symbol		string	`json:"symbol"`
		RateFloat	float64	`json:"rate_float"`
	}	`json:"bpi"`
}

var c Root

func (c *Root) HandleClick(event vugu.DOMEvent) {

	c.bpi = bpi{}

	go func(ee vugu.EventEnv) {

		ee.Lock()
		c.isLoading = true
		ee.UnlockRender()

		res, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
		if err != nil {
			log.Printf("Error fetch()ing: %v", err)
			return
		}
		defer res.Body.Close()

		var newb bpi
		err = json.NewDecoder(res.Body).Decode(&newb)
		if err != nil {
			log.Printf("Error JSON decoding: %v", err)
			return
		}

		ee.Lock()
		defer ee.UnlockRender()
		c.bpi = newb
		c.isLoading = false

	}(event.EventEnv())
}

func (c *Root) Build(vgin *vugu.BuildIn) (vgout *vugu.BuildOut) {

	vgout = &vugu.BuildOut{}

	var vgiterkey interface{}
	_ = vgiterkey
	var vgn *vugu.VGNode
	vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "div", Attr: []vugu.VGAttribute{vugu.VGAttribute{Namespace: "", Key: "class", Val: "demo-comp"}}}
	vgout.Out = append(vgout.Out, vgn)	// root for output
	{
		vgparent := vgn
		_ = vgparent
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n    "}
		vgparent.AppendChild(vgn)
		if c.isLoading {
			vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "div", Attr: []vugu.VGAttribute(nil)}
			vgparent.AppendChild(vgn)
			{
				vgparent := vgn
				_ = vgparent
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "Loading..."}
				vgparent.AppendChild(vgn)
			}
		}
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n    "}
		vgparent.AppendChild(vgn)
		if len(c.bpi.BPI) > 0 {
			vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "div", Attr: []vugu.VGAttribute(nil)}
			vgparent.AppendChild(vgn)
			{
				vgparent := vgn
				_ = vgparent
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n        "}
				vgparent.AppendChild(vgn)
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "div", Attr: []vugu.VGAttribute(nil)}
				vgparent.AppendChild(vgn)
				{
					vgparent := vgn
					_ = vgparent
					vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "Updated: "}
					vgparent.AppendChild(vgn)
					vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "span", Attr: []vugu.VGAttribute(nil)}
					vgparent.AppendChild(vgn)
					vgn.SetInnerHTML(c.bpi.Time.Updated)
				}
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n        "}
				vgparent.AppendChild(vgn)
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "ul", Attr: []vugu.VGAttribute(nil)}
				vgparent.AppendChild(vgn)
				{
					vgparent := vgn
					_ = vgparent
					vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n            "}
					vgparent.AppendChild(vgn)
					for key, value := range c.bpi.BPI {
						var vgiterkey interface{} = key
						_ = vgiterkey
						key := key
						_ = key
						value := value
						_ = value
						vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "li", Attr: []vugu.VGAttribute(nil)}
						vgparent.AppendChild(vgn)
						{
							vgparent := vgn
							_ = vgparent
							vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n                "}
							vgparent.AppendChild(vgn)
							vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "span", Attr: []vugu.VGAttribute(nil)}
							vgparent.AppendChild(vgn)
							vgn.SetInnerHTML(key)
							vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: " "}
							vgparent.AppendChild(vgn)
							vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "span", Attr: []vugu.VGAttribute(nil)}
							vgparent.AppendChild(vgn)
							vgn.SetInnerHTML(fmt.Sprint(value.Symbol, value.RateFloat))
							vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n            "}
							vgparent.AppendChild(vgn)
						}
					}
					vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n        "}
					vgparent.AppendChild(vgn)
				}
				vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n    "}
				vgparent.AppendChild(vgn)
			}
		}
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n    "}
		vgparent.AppendChild(vgn)
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(3), Namespace: "", Data: "button", Attr: []vugu.VGAttribute(nil)}
		vgparent.AppendChild(vgn)
		vgn.DOMEventHandlerSpecList = append(vgn.DOMEventHandlerSpecList, vugu.DOMEventHandlerSpec{
			EventType:	"click",
			Func:		func(event vugu.DOMEvent) { c.HandleClick(event) },
			// TODO: implement capture, etc. mostly need to decide syntax
		})
		{
			vgparent := vgn
			_ = vgparent
			vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "Fetch Bitcoin Price Index"}
			vgparent.AppendChild(vgn)
		}
		vgn = &vugu.VGNode{Type: vugu.VGNodeType(1), Data: "\n"}
		vgparent.AppendChild(vgn)
	}
	return vgout
}

// 'fix' unused imports
var _ fmt.Stringer
var _ reflect.Type
var _ vjson.RawMessage
var _ js.Value
