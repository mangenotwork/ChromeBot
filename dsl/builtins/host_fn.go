package builtins

import (
	"ChromeBot/dsl/interpreter"
	"fmt"
	"strings"

	gt "github.com/mangenotwork/gathertool"
)

var hostFn = map[string]interpreter.Function{
	"NsLookUp":           hostNsLookUp,           // NsLookUp(host) DNS查询方法
	"Whois":              hostWhois,              // Whois(host) Whois查询方法
	"SearchPort":         hostSearchPort,         // SearchPort(ip) 端口扫描方法
	"WebSiteScanBadLink": hostWebSiteScanBadLink, // WebSiteScanBadLink(domain, depth) 网站死链检查, depth是遍历网站的深度
	"WebCertificateInfo": hostWebCertificateInfo, // WebCertificateInfo(domain) 网站证书信息
	"WebScanUrl":         hostWebScanUrl,         // WebScanUrl(domain) NewHostScanUrl 创建扫描站点
	"WebScanExtLinks":    hostWebScanExtLinks,    // WebScanExtLinks(domain) 创建站点链接采集，只支持get请求
	"WebPageSpeedCheck":  hostWebPageSpeedCheck,  // WebPageSpeedCheck(domain, depth) NewHostPageSpeedCheck 创建站点所有url测速，只支持get请求
}

func hostNsLookUp(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("NsLookUp(host) 需要一个参数")
	}

	host, hostOK := args[0].(string)
	if !hostOK {
		return nil, fmt.Errorf("NsLookUp(host)  参数要求是字符串 ")
	}

	nsRes := gt.NsLookUp(host)

	fmt.Println("NsLookUp DnsServerName : ", nsRes.DnsServerName)
	fmt.Println("NsLookUp DnsServerIP : ", nsRes.DnsServerIP)
	fmt.Println("NsLookUp IPs : ", strings.Join(nsRes.IPs, ","))
	fmt.Println("NsLookUp IsCDN : ", nsRes.IsCDN)
	fmt.Println("NsLookUp LookupCNAME : ", nsRes.LookupCNAME)
	fmt.Println("NsLookUp Ms : ", nsRes.Ms)

	return nil, nil
}

func hostWhois(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Whois(host) 需要一个参数")
	}

	host, hostOK := args[0].(string)
	if !hostOK {
		return nil, fmt.Errorf("NsLookUp(host)  参数要求是字符串 ")
	}

	whoisInfo := gt.Whois(host)

	fmt.Println("Whois Root : ", whoisInfo.Root)
	fmt.Println("Whois Rse : ", whoisInfo.Rse)

	return nil, nil
}

func hostSearchPort(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SearchPort(ip) 需要一个参数")
	}

	fmt.Println("SearchPort ....")

	// todo....

	return nil, nil
}

func hostWebSiteScanBadLink(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("WebSiteScanBadLink(domain, depth) 需要一个参数")
	}

	fmt.Println("WebSiteScanBadLink ....")

	// todo....

	return nil, nil
}

func hostWebCertificateInfo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WebCertificateInfo(domain) 需要一个参数")
	}

	fmt.Println("WebCertificateInfo ....")

	// todo....

	return nil, nil
}

func hostWebScanUrl(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WebScanUrl(domain) 需要一个参数")
	}

	fmt.Println("WebScanUr ....")

	// todo....

	return nil, nil
}

func hostWebScanExtLinks(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WebScanExtLinks(domain) 需要一个参数")
	}

	fmt.Println("WebScanExtLinks ....")

	// todo....

	return nil, nil
}

func hostWebPageSpeedCheck(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("WebPageSpeedCheck(domain, depth) 需要一个参数")
	}

	fmt.Println("WebPageSpeedCheck ....")

	// todo....

	return nil, nil
}
