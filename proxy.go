package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// 创建一个反向代理处理程序
	proxy := &ProxyHandler{}

	// 启动代理服务器
	log.Println("Proxy server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

// ProxyHandler 是一个反向代理处理程序
type ProxyHandler struct{}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建一个新的请求，将客户端的请求转发到百度
	req, err := http.NewRequest(r.Method, "https://api.openai.com"+r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 复制客户端请求的 Header 到新的请求中
	req.Header = r.Header

	// 创建一个 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 将百度的响应复制到客户端的响应中
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// 复制 Header
// copyHeader copies the header from src to dst.
func copyHeader(dst, src http.Header) {
	// Iterate through each key/value pair in src
	for key, values := range src {
		// Iterate through each value in the current key/value pair
		for _, value := range values {
			// Add the key/value pair to dst
			dst.Add(key, value)
		}
	}
}
