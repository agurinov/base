package tools

import (
	"fmt"
	"testing"
)

func TestIsBot(t *testing.T) {
	tableTests := []struct {
		ua  string // url for parsing
		bot bool   // is it bot?
	}{
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) Safari/600.2.5 (Applebot/0.1)", true},       // apple applebot
		{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", true},            // google
		{"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)", true},             // bingbot
		{"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", true}, // yahoo
		{"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)", true},                          // duckduck
		{"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)", true}, // baiduspider
		{"facebookexternalhit/1.0 (+http://www.facebook.com/externalhit_uatext.php)", true},           //facebook
		{"Mozilla/5.0 (compatible; Twitterbot/1.0)", true},                                            // twitter
		{"Mozilla/5.0 (Windows NT 6.1; WOW64) Safari/537.36 MicroMessenger/6.5.2.501", true},          // wechat

		{"", false}, // empty
		{"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) Version/5.0 Mobile/11A465 Safari/9537.53", false}, // usual
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if isBot := IsBot(tt.ua); isBot != tt.bot {
				t.Errorf("Expected \"%t\", got \"%t\"", tt.bot, isBot)
			}
		})
	}
}
