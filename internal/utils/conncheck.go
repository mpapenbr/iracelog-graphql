package utils

import (
	"fmt"

	"net"
	"regexp"
	"time"
)

func WaitForTCP(addr string, timeout time.Duration) error {
	timeoutReached := time.Now().Add(timeout)
	start := time.Now()
	fmt.Printf("wait for tcp connection %s (timeout: %s)\n", addr, timeout.String())		
	for time.Now().Before(timeoutReached) {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()

			fmt.Printf("tcp connection to %s successful after %s\n",  addr, time.Since(start).String())
				
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("%s could not be reached after %v", addr, timeout)
}



func ExtractFromDBUrl(url string) string {
	param := resolveRegex(
		"^postgresql://(.*@)(?P<addr>(?P<host>.*?)(:(?P<port>\\d+))?)/.*", url)
	if len(param) == 0 {
		return ""
	}
	if port, ok := param["port"]; ok && len(port) > 0 {
		return param["addr"] // if port is found, the addr contains our wanted value
	} else {
		return fmt.Sprintf("%s:5432", param["addr"])
	}
}

func resolveRegex(regEx, url string) (paramsMap map[string]string) {
	compRegEx := regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
