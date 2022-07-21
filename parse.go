package parse

import (
	"errors"
	"strings"

	"github.com/chyroc/lark"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

/*
	ErrorNode NodeType = iota
	TextNode // TextNode is a text node.
	DocumentNode // DocumentNode is the root of the document.
	ElementNode // ElementNode is an element.
	CommentNode // CommentNode is a comment node.
	DoctypeNode // DoctypeNode is a document type node.
	RawNode
	scopeMarkerNode
*/

/*
一般而言，html.Parse之后，根节点是DocumentNode，无兄弟节点，其子节点是ElementNode。
html 标签也会作为一个node节点存在，这个node节点的Type为ElementNode，Data为标签名，Attr为属性列表。
其中 html、header、body标签都会自动补全
*/

func ParseHtmlToLarkPostMessage(htmlStr string, cli *lark.Lark) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	data := &lark.MessageContentPostAll{
		ZhCn: &lark.MessageContentPost{},
		JaJp: nil,
		EnUs: nil,
	}

	var f func(*html.Node) error
	f = func(n *html.Node) error {
		switch n.Type {
		case html.ErrorNode:
			return errors.New("html: ErrorNode")
		case html.TextNode:
			// text node don't have children
			if n.Data != "" {
				data.ZhCn.Content = append(data.ZhCn.Content, []lark.MessageContentPostItem{
					lark.MessageContentPostText{Text: n.Data},
				})
			}
		case html.ElementNode:
			switch n.DataAtom {
			case atom.Img:
				// img node don't have children
				data.ZhCn.Content = append(data.ZhCn.Content, []lark.MessageContentPostItem{
					// lark.MessageContentPostImage{URL: n.Attr[0].Val},
				})
			case atom.A:
				// a node have children,but lark just support one and text child
				href := ""
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}
				text := ""
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						text = c.Data
						break
					}
				}
				if href != "" && text != "" {
					data.ZhCn.Content = append(data.ZhCn.Content, []lark.MessageContentPostItem{
						lark.MessageContentPostLink{
							Text:     text,
							UnEscape: false,
							Href:     href,
						},
					})
				}
				// don't parse children
				return nil
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			err = f(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err = f(doc)
	if err != nil {
		return "", err
	}

	return data.String(), nil
}
