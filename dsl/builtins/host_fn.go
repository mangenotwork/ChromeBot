package builtins

import (
	"ChromeBot/dsl/interpreter"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	gt "github.com/mangenotwork/gathertool"
)

var hostFn = map[string]interpreter.Function{
	"NsLookUp":           hostNsLookUp,           // NsLookUp(host) DNS查询方法
	"Whois":              hostWhois,              // Whois(host) Whois查询方法
	"SearchPort":         hostSearchPort,         // SearchPort(ip) 端口扫描方法
	"WebSiteScanBadLink": hostWebSiteScanBadLink, // WebSiteScanBadLink(domain, depth) 网站死链检查, depth是遍历网站的深度
	"WebCertificateInfo": hostWebCertificateInfo, // WebCertificateInfo(domain) 网站证书信息
	"WebScanUrl":         hostWebScanUrl,         // WebScanUrl(domain, depth) NewHostScanUrl 创建扫描站点
	"WebScanExtLinks":    hostWebScanExtLinks,    // WebScanExtLinks(domain, depth) 创建站点链接采集，只支持get请求
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

	ip, ipOK := args[0].(string)
	if !ipOK {
		return nil, fmt.Errorf("SearchPort(ip)  参数要求是字符串 ")
	}

	timeOut := 4 * time.Second
	queue := gt.NewQueue()
	for i := 0; i < 65536; i++ {
		buf := &bytes.Buffer{}
		buf.WriteString(ip)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(i))
		_ = queue.Add(&gt.Task{
			Url: buf.String(),
		})
	}

	var wg sync.WaitGroup
	for job := 0; job < 65536; job++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				if task == nil {
					continue
				}
				conn, err := net.DialTimeout("tcp", task.Url, timeOut)
				if err == nil {
					fmt.Println(task.Url, "开放")
					_ = conn.Close()
				}
			}
		}(job)
	}
	wg.Wait()
	fmt.Println("端口扫描完成！！！")

	return nil, nil
}

func hostWebSiteScanBadLink(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("WebSiteScanBadLink(domain, depth) 需要一个参数")
	}

	fmt.Println("执行死链检查 ....")

	domain, domainOK := args[0].(string)
	if !domainOK {
		return nil, fmt.Errorf("WebSiteScanBadLink(domain, depth)  参数要求是字符串 ")
	}

	depth, depthOK := args[1].(int64)
	if !depthOK {
		return nil, fmt.Errorf("WebSiteScanBadLink(domain, depth)  depth参数要求是整数 ")
	}

	obj := gt.NewHostScanBadLink(domain, int(depth))
	res, count := obj.Run()
	fmt.Println("死链数量: ", count)
	fmt.Println("死链: ", strings.Join(res, ";"))
	return nil, nil
}

func hostWebCertificateInfo(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WebCertificateInfo(domain) 需要一个参数")
	}

	domain, domainOK := args[0].(string)
	if !domainOK {
		return nil, fmt.Errorf("WebSiteScanBadLink(domain, depth)  参数要求是字符串 ")
	}

	certificateInfo, _ := gt.GetCertificateInfo(domain)

	fmt.Println("证书 url : ", certificateInfo.Url)
	fmt.Println("证书 有效时间 : ", certificateInfo.EffectiveTime)
	fmt.Println("证书 起始 : ", certificateInfo.NotBefore)
	fmt.Println("证书 结束 : ", certificateInfo.NotAfter)
	fmt.Println("证书 DNSName : ", certificateInfo.DNSName)
	fmt.Println("证书 OCSPServer : ", certificateInfo.OCSPServer)
	fmt.Println("证书 CRL分发点 : ", certificateInfo.CRLDistributionPoints)
	fmt.Println("证书 颁发者 : ", certificateInfo.Issuer)
	fmt.Println("证书 颁发证书URL : ", certificateInfo.IssuingCertificateURL)
	fmt.Println("证书 公钥算法 : ", certificateInfo.PublicKeyAlgorithm)
	fmt.Println("证书 颁发对象 : ", certificateInfo.Subject)
	fmt.Println("证书 版本 : ", certificateInfo.Version)
	fmt.Println("证书 证书算法 : ", certificateInfo.SignatureAlgorithm)

	res := interpreter.DictType{
		"Url":                   certificateInfo.Url,
		"EffectiveTime":         certificateInfo.EffectiveTime,
		"NotBefore":             certificateInfo.NotBefore,
		"NotAfter":              certificateInfo.NotAfter,
		"DNSName":               certificateInfo.DNSName,
		"OCSPServer":            certificateInfo.OCSPServer,
		"CRLDistributionPoints": certificateInfo.CRLDistributionPoints,
		"Issuer":                certificateInfo.Issuer,
		"IssuingCertificateURL": certificateInfo.IssuingCertificateURL,
		"PublicKeyAlgorithm":    certificateInfo.PublicKeyAlgorithm,
		"Subject":               certificateInfo.Subject,
		"Version":               certificateInfo.Version,
		"SignatureAlgorithm":    certificateInfo.SignatureAlgorithm,
	}

	return res, nil
}

func hostWebScanUrl(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("WebScanUrl(domain, depth) 需要一个参数")
	}

	fmt.Println("WebScanUrl ....")

	domain, domainOK := args[0].(string)
	if !domainOK {
		return nil, fmt.Errorf("WebScanUrl(domain, depth)  参数要求是字符串 ")
	}

	depth, depthOK := args[1].(int64)
	if !depthOK {
		return nil, fmt.Errorf("WebScanUrl(domain, depth)  depth参数要求是整数 ")
	}

	obj := gt.NewHostScanUrl(domain, int(depth))
	res, count := obj.Run()
	fmt.Println("数量: ", count)
	fmt.Println("扫描站点的url: ", strings.Join(res, ";"))
	return nil, nil
}

func hostWebScanExtLinks(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WebScanExtLinks(domain) 需要一个参数")
	}

	fmt.Println("WebScanExtLinks ....")

	domain, domainOK := args[0].(string)
	if !domainOK {
		return nil, fmt.Errorf("WebScanExtLinks(domain)  参数要求是字符串 ")
	}

	obj := gt.NewHostScanExtLinks(domain)
	res, count := obj.Run()
	fmt.Println("数量: ", count)
	fmt.Println("站点链接采集: ", strings.Join(res, ";"))
	return nil, nil
}

func hostWebPageSpeedCheck(args []interpreter.Value) (interpreter.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("WebPageSpeedCheck(domain, depth) 需要一个参数")
	}

	fmt.Println("WebPageSpeedCheck ....")

	domain, domainOK := args[0].(string)
	if !domainOK {
		return nil, fmt.Errorf("WebPageSpeedCheck(domain, depth)  参数要求是字符串 ")
	}

	depth, depthOK := args[1].(int64)
	if !depthOK {
		return nil, fmt.Errorf("WebPageSpeedCheck(domain, depth)  depth参数要求是整数 ")
	}

	obj := gt.NewHostPageSpeedCheck(domain, int(depth))
	res, count := obj.Run()
	fmt.Println("数量: ", count)
	fmt.Println("站点所有url测速: ", strings.Join(res, ";"))
	return nil, nil
}
