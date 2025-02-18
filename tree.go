package onebotClient

import (
	"github.com/For-December/onebotClient/msg"
	"strings"
)

// RouteTrieNode 路由基数树节点
type RouteTrieNode struct {
	children map[string]*RouteTrieNode
	//isParam   bool
	paramName string
	handlers  []PluginHandler
}

// RouteTrie 路由基数树
type RouteTrie struct {
	callbackFunc CallbackFunc
	root         *RouteTrieNode

	depth int
}

// NewRouteTrie 创建一个新的路由基数树
func NewRouteTrie(callbackFunc CallbackFunc) *RouteTrie {

	// 如果没有设置回调函数，使用默认的回调函数，可穿透
	if callbackFunc.OnNotFound == nil {
		callbackFunc.OnNotFound = append(callbackFunc.OnNotFound,
			func(api BotActionAPIInterface,
				ctx *msg.GroupMessageContext) bool {
				return false
			})
	}

	if callbackFunc.AfterEach == nil {
		callbackFunc.AfterEach = append(callbackFunc.AfterEach,
			func(api BotActionAPIInterface,
				ctx *msg.GroupMessageContext) bool {
				return false
			})
	}

	return &RouteTrie{
		root: &RouteTrieNode{
			children: make(map[string]*RouteTrieNode),
		},
		callbackFunc: callbackFunc,
	}
}

// Insert 在路由基数树中插入路由和对应的处理函数
// 优先级 !! > 精确匹配 > **, **只能放到最后
// !! 代表无论什么情况都会执行的回调函数，按插入顺序先后排优先级
func (t *RouteTrie) Insert(path string, handler PluginHandler) {
	if path == "!!" {
		// 无论什么情况都会执行的回调函数，按插入顺序先后排优先级
		t.callbackFunc.AfterEach = append(t.callbackFunc.AfterEach, handler)
		t.callbackFunc.OnNotFound = append(t.callbackFunc.OnNotFound, handler)
		return
	}

	parts := strings.Split(path, " ")
	node := t.root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if part[0] == '$' { // 参数节点，用空字符串表示该特殊节点
			if _, ok := node.children[""]; !ok {
				node.children[""] = &RouteTrieNode{
					children: make(map[string]*RouteTrieNode),
					//isParam:   true,
					paramName: part[1:],
				}
			}

			// 将子节点作为当前节点
			node = node.children[""]
		} else { // 包含精确匹配和通配符匹配**
			if _, ok := node.children[part]; !ok {
				node.children[part] = &RouteTrieNode{
					children: make(map[string]*RouteTrieNode),
				}
			}
			node = node.children[part]
		}

		t.depth++

	}

	node.handlers = append(node.handlers, handler)
}

// Search 在路由基数树中查找路由对应的处理函数
func (t *RouteTrie) Search(path string) []PluginHandler {
	parts := strings.Split(path, " ")
	node := t.root
	params := make(map[string]string)

	for i, part := range parts {
		if part == "" {
			continue
		}

		// 如果当前是路径的一部分，继续查找
		if _, ok := node.children[part]; ok {
			node = node.children[part]

			// 否则当前可能是参数，查找此路径下面的参数节点
		} else if _, ok = node.children[""]; ok {
			params[node.children[""].paramName] = part
			node = node.children[""]

			// 否则可能是通配符节点
		} else if _, ok = node.children["**"]; ok {
			// 将指令后面的所有内容作为**匹配到的值（所以通配符**只能放到最后）
			node = node.children["**"]
			params["**"] = strings.Join(parts[i:], " ")
			break
		} else {
			// 不是参数节点，也找不到这部分对应的精确节点或通配符节点，返回404回调函数
			return t.callbackFunc.OnNotFound
		}
	}

	if node.handlers == nil {
		return t.callbackFunc.OnNotFound
	}

	return []PluginHandler{
		func(api BotActionAPIInterface,
			ctx *msg.GroupMessageContext) bool {
			// 代理函数
			if ctx != nil {
				ctx.Params = params
			}

			for _, handler := range node.handlers {
				if handler(api, ctx) {
					return true
				}
			}

			return false
		},
	}
}

// SearchAndExec 在路由基数树中查找路由对应的处理函数并执行
func (t *RouteTrie) SearchAndExec(api BotActionAPIInterface, ctx *msg.GroupMessageContext) {
	if ctx == nil {
		panic("MessageContext is nil")
		return
	}

	parts := strings.Split(ctx.GroupMessageChain.ToPath(), " ")
	node := t.root
	params := make(map[string]string)

	for i, part := range parts {
		if part == "" {
			continue
		}

		// 如果当前是路径的一部分，继续查找
		if _, ok := node.children[part]; ok {
			node = node.children[part]

			// 否则当前可能是参数，查找此路径下面的参数节点
		} else if _, ok = node.children[""]; ok {
			params[node.children[""].paramName] = part
			node = node.children[""]
		} else if _, ok = node.children["**"]; ok {
			// 将指令后面的所有内容作为**匹配到的值（所以通配符**只能放到最后）
			node = node.children["**"]
			params["**"] = strings.Join(parts[i:], " ")
			break
		} else {
			//不是参数节点，也找不到这部分对应的精确节点或通配符节点，执行所有404回调函数
			for _, handler := range t.callbackFunc.OnNotFound {
				handler(api, ctx)
			}
			return
		}
	}

	if node.handlers == nil {
		for _, handler := range t.callbackFunc.OnNotFound {
			handler(api, ctx)
		}
		return
	}

	ctx.Params = params
	for _, handler := range node.handlers {
		// 终止
		if handler(api, ctx) {
			break
		}
	}

	// 执行所有AfterEach回调函数
	for _, handler := range t.callbackFunc.AfterEach {
		if handler(api, ctx) {
			break
		}
	}

}
