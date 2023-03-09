package utils

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"time"
	"xc7/log"
	"xc7/utils/celpro"

	"github.com/imroc/req/v3"
)

func SendRequest(Method string, url string, body io.Reader, hdrs map[string]string, proxy string, redirects bool) (error, *celpro.Response_C) {

	var err error
	var res *req.Response
	var resp celpro.Response_C

	user_agent := make(map[string]string)
	user_agent["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

	client := req.C().EnableInsecureSkipVerify()

	if redirects == false {
		client.SetRedirectPolicy(req.NoRedirectPolicy())
	}
	if proxy != "" {
		client.SetProxyURL(proxy)
	}
	switch Method {
	case "GET":
		res, err = client.R().SetBody(body).SetHeaders(user_agent).SetHeaders(hdrs).Get(url)
	case "POST":
		res, err = client.R().SetBody(body).SetHeaders(user_agent).SetHeaders(hdrs).Post(url)
	case "PUT":
		res, err = client.R().SetBody(body).SetHeaders(user_agent).SetHeaders(hdrs).Put(url)
	case "DELETE":
		res, err = client.R().SetBody(body).SetHeaders(user_agent).SetHeaders(hdrs).Delete(url)
	case "PATCH":
		res, err = client.R().SetBody(body).SetHeaders(user_agent).SetHeaders(hdrs).Patch(url)
	case "OPTIONS":
		log.GreenOut("OPTIONS方法还未实现")
	}
	if err != nil {
		return err, nil
	}
	/*
	 * 这是处理前后的对比
	 * map[Content-Length:[31] Content-Type:[text/html; charset=UTF-8] Date:[Tue, 20 Dec 2022 05:00:58 GMT] Server:[Apache/2.4.38 (Debian)] X-Powered-By:[PHP/7.2.31]]
	 * map[Content-Length:31 Content-Type:text/html; charset=UTF-8 Date:Tue, 20 Dec 2022 05:00:58 GMT Server:Apache/2.4.38 (Debian) X-Powered-By:PHP/7.2.31]
	 */
	header := make(map[string]string)
	for k := range res.Header {
		header[k] = res.Header.Get(k)
	}

	resp.Headers = header
	resp.Status = int32(res.StatusCode)
	resp.Url = ParseUrl(res.Request.URL)

	resp.Body = res.Bytes()
	resp.ContentType = res.GetContentType()

	return err, &resp
}

/*
	*tips
	反连平台
*/

// http://ceye.io/api

func ReverseCheck(r *celpro.Reverse_C, timeout int64) bool {
	_, ceyeapi := LoadConfig()
	if ceyeapi == "" {
		log.ErrorLog("未配置ceye-api")
		return false
	}
	log.InfoLog("等待平台是否响应...")

	time.Sleep(time.Second * time.Duration(timeout))
	flag := r.Flag
	apiurl := fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s", ceyeapi, flag)
	client := req.C()
	res, err := client.R().Get(apiurl)

	if err != nil {
		log.ErrorLog("发生网络错误有，无法获取ceye api")
	}
	// 如果api返回有结果那么就不存在 "data": []，正常[]内会有内容。
	if bytes.Contains(res.Bytes(), []byte(`"data": []`)) == false && bytes.Contains(res.Bytes(), []byte(`"message": "OK"`)) == true {
		return true
	}
	log.InfoLog("未检测到返回结果！")
	return false
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomLowercase(n int) string {
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	chararr := strings.Split(alphabet, "")
	charlen := len(chararr)
	var flag string = ""
	for i := 0; i <= n; i++ {
		flag = flag + chararr[ran.Intn(charlen)]
	}
	return flag
}

// newReverse解析出来返回 *celpro.Reverse_C （还有UrlType_C）

func NewReverse() *celpro.Reverse_C {
	flag := RandomLowercase(9)
	ceye, _ := LoadConfig()
	if ceye == "" {
		log.ErrorLog("未配置反连平台地址")
		return nil
	}

	ceyeurl := fmt.Sprintf("http://%s.%s/", flag, ceye)
	log.InfoLog("Use:" + ceyeurl)
	// 解析ceyeurl,
	u, _ := url.Parse(ceyeurl)
	return &celpro.Reverse_C{
		Url: ParseUrl(u),

		Domain:             u.Hostname(),
		Ip:                 "",
		IsDomainNameServer: false,
		Flag:               flag,
	}
}

/*
 * 解析url出来给proto那些
 */

func ParseUrl(u *url.URL) *celpro.UrlType_C {

	nu := &celpro.UrlType_C{}
	nu.Domain = u.Hostname()  // Hostname返回u.Host
	nu.Fragment = u.Fragment  // 匹配#之后的字段
	nu.Host = u.Host          // host or host:port
	nu.Path = u.EscapedPath() // 返回网站目录
	nu.Port = u.Port()        // 端口
	nu.Query = u.RawQuery     // 匹配?后的字段
	nu.Scheme = u.Scheme      // http 或https
	return nu
}

func UrltypeCTostring(u *celpro.UrlType_C) string {

	var buf strings.Builder

	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}

	if u.Scheme != "" || u.Host != "" {
		if u.Host != "" || u.Path != "" {
			buf.WriteString("//")
		}
		if h := u.Host; h != "" {
			buf.WriteString(u.Host)
		}
	}
	path := u.Path
	if path != "" && path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)

	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}
