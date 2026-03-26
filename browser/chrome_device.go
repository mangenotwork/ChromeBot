package browser

type deviceArg struct {
	userAgent  string
	windowSize string
}

var chromeDevice = map[string]deviceArg{
	"iphone": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=430,932`,
	},
	"iphone15": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=393,852`,
	},
	"iphone15P": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=430,932`,
	},
	"iphone14": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=390,844`,
	},
	"iphone13": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=390,844`,
	},
	"iphone12": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=390,844`,
	},
	"iphoneES": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=375,667`,
	},
	"iphone7": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"`,
		windowSize: `--window-size=375,667`,
	},
	"ipad11": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"`,
		windowSize: `--window-size=834,1194`,
	},
	"ipad12": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"`,
		windowSize: `--window-size=1024,1366`,
	},
	"ipadAir": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15"`,
		windowSize: `--window-size=820,1180`,
	},
	"ipadMini": {
		userAgent:  `--user-agent="Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15"`,
		windowSize: `--window-size=768,1024`,
	},
	"android": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; Pixel 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=412,868`,
	},
	"galaxy": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=384,854`,
	},
	"galaxyS24": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=384,854`,
	},
	"galaxyS23": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=384,854`,
	},
	"galaxyS22": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 12; SM-S908B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=360,800`,
	},
	"galaxyZFold5": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; SM-F946B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=384,854`,
	},
	"huawei": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 10; ALN-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.225 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"huaweiMate60": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 10; ALN-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.225 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"huaweiPura70": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 10; NAM-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.225 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"huaweiMagic6": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; BVL-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"xiaomi": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; 24031PN0DC) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"xiaomi14": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; 24031PN0DC) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"xiaomi13": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; 2211133C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=360,800`,
	},
	"redmi": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; 23090RA98C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=360,800`,
	},
	"oppo": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; PHY110) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"vivo": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; V2324A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=396,888`,
	},
	"pixel": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; Pixel 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=412,868`,
	},
	"pixel8": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 14; Pixel 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=412,868`,
	},
	"pixel7": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; Pixel 7a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36"`,
		windowSize: `--window-size=412,915`,
	},
	"androidPad": {
		userAgent:  `--user-agent="Mozilla/5.0 (Linux; Android 13; 23043RP34G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"`,
		windowSize: `--window-size=1080,720`,
	},
}
