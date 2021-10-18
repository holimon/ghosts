package ghosts

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/holimon/requests"
	"github.com/opesun/goquery"
)

var regFilterComments = regexp.MustCompile("(^|\n)[^#](.+)")
var regParseRecord = regexp.MustCompile(`[^\s]+`)

func hostsName() (string, error) {
	if runtime.GOOS == "windows" {
		hosts := filepath.Join(os.Getenv("windir"), "System32", "drivers", "etc", "hosts")
		return hosts, nil
	}
	if runtime.GOOS == "linux" {
		return "/etc/hosts", nil
	}
	return "", fmt.Errorf("%s not support", runtime.GOOS)
}

func hostsSep() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	if runtime.GOOS == "linux" {
		return "\n"
	}
	return "\n"
}

func readHosts() (hosts string, err error) {
	name, err := hostsName()
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(name)
	return string(b), err
}

func parseRecord(hosts string) (records []Records) {
	records = []Records{}
	recordLine := regFilterComments.FindAllString(hosts, -1)
	for _, line := range recordLine {
		parse := regParseRecord.FindAllString(line, -1)
		if len(parse) == 2 && net.ParseIP(parse[0]) != nil {
			hash := md5.New().Sum([]byte(parse[0] + parse[1]))
			records = append(records, Records{Domain: parse[1], IP: parse[0], md5: hex.EncodeToString(hash)})
		}
	}
	return records
}

func writeHosts(hosts string) error {
	if name, err := hostsName(); err == nil {
		return ioutil.WriteFile(name, []byte(hosts), os.ModePerm)
	} else {
		return err
	}
}

func defaultDebug(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func defaultResolve(domain string) (address []string, err error) {
	res, err := requests.Requests().Get(fmt.Sprintf("https://websites.ipaddress.com/%s", domain))
	if err != nil {
		return []string{}, err
	}
	etree, err := goquery.ParseString(string(res.Content()))
	if err != nil {
		return []string{}, err
	}
	nodes := etree.Find("strong")
	for i := 0; i < nodes.Length(); i++ {
		s := nodes.Eq(i).Text()
		if net.ParseIP(s) != nil {
			address = append(address, s)
		}
	}
	return address, nil
}
