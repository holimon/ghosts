package ghosts

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type Records struct{ IP, Domain, md5 string }

type RecordsSort struct {
	records  []Records
	lessRule func(i, j Records) bool
}

type SortedType int

const (
	SortedByDomain SortedType = 0
	SortedByIP     SortedType = 1
)

func (r RecordsSort) Len() int { return len(r.records) }

func (r RecordsSort) Swap(i, j int) { r.records[i], r.records[j] = r.records[j], r.records[i] }

func (r RecordsSort) Less(i, j int) bool {
	return r.lessRule(r.records[i], r.records[j])
}

type Option struct {
	DebugInfo     func(format string, a ...interface{}) (n int, err error)
	CustomResolve func(domain string) (address []string, err error)
	DNSFlushName  string
	DNSFlushArgs  []string
}

type ghosts struct {
	opt     Option
	hosts   string
	records []Records
	lock    sync.RWMutex
}

type Ghosts = *ghosts

func (h *ghosts) recordsSorted(sorteType SortedType) []Records {
	sorted := RecordsSort{records: h.records, lessRule: func(i, j Records) bool {
		if j.Domain == i.Domain {
			return j.IP > i.IP
		}
		return j.Domain > i.Domain
	}}
	if sorteType == SortedByIP {
		sorted.lessRule = func(i, j Records) bool {
			if j.IP == i.IP {
				return j.Domain > i.Domain
			}
			return j.IP > i.IP
		}
	}
	sort.Stable(sorted)
	return sorted.records
}

func (h *ghosts) recordsAll() []Records {
	return h.records
}

func (h *ghosts) recordDelbyDomain(param string) error {
	records := []Records{}
	for _, v := range h.records {
		if param != v.Domain {
			records = append(records, v)
		}
	}
	h.records = records
	return nil
}

func (h *ghosts) recordDelbyIP(param string) error {
	records := []Records{}
	for _, v := range h.records {
		if param != v.IP {
			records = append(records, v)
		}
	}
	h.records = records
	return nil
}

func (h *ghosts) recordDelbyMD5(param string) error {
	records := []Records{}
	for _, v := range h.records {
		if param != v.md5 {
			records = append(records, v)
		}
	}
	h.records = records
	return nil
}

func (h *ghosts) recordAdd(domain, address string) error {
	if net.ParseIP(address) != nil {
		h.records = append(h.records, Records{Domain: domain, IP: address})
		return nil
	}
	return fmt.Errorf("IP address is not valid")
}

func (h *ghosts) recordReplace(domain, address string) error {
	if err := h.recordDelbyDomain(domain); err != nil {
		return err
	}
	return h.recordAdd(domain, address)
}

func (h *ghosts) recordsWrite() error {
	sep := hostsSep()
	content := regFilterComments.ReplaceAllString(h.hosts, "")
	tips := fmt.Sprintf("# The following is appended by the ghost editor.%s", sep)
	content = strings.ReplaceAll(content, tips, "")
	content += tips
	for _, item := range h.recordsSorted(SortedByDomain) {
		content += fmt.Sprintf("%-48s%s%s", item.IP, item.Domain, sep)
	}
	return writeHosts(content)
}

func newGHosts(opt Option) (h *ghosts, err error) {
	h = &ghosts{opt: opt, records: []Records{}}
	if h.opt.DebugInfo == nil {
		h.opt.DebugInfo = defaultDebug
	}
	if h.opt.CustomResolve == nil {
		h.opt.CustomResolve = defaultResolve
	}
	h.hosts, err = readHosts()
	h.records = parseRecord(h.hosts)
	return h, err
}

func New(opt Option) (h Ghosts, err error) {
	return newGHosts(opt)
}

func (g Ghosts) Add(domain, address string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.recordAdd(domain, address)
}

func (g Ghosts) Delete(param string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if net.ParseIP(param) != nil {
		return g.recordDelbyIP(param)
	}
	return g.recordDelbyDomain(param)
}

func (g Ghosts) DeleteRecord(record Records) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.recordDelbyMD5(record.md5)
}

func (g Ghosts) Replace(domain, address string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.recordReplace(domain, address)
}

func (g Ghosts) Records() []Records {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.recordsAll()
}

func (g Ghosts) OrdedRecords(sorteType SortedType) []Records {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.recordsSorted(sorteType)
}

func (g Ghosts) Resolve(domain string) ([]string, error) {
	addrs, err := g.opt.CustomResolve(domain)
	if err == nil {
		g.opt.DebugInfo("resolve domain %s to %v\n", domain, addrs)
	}
	return addrs, err
}

func (g Ghosts) WriteBack() error {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.recordsWrite()
}

func (g Ghosts) DNSFlush() error {
	if g.opt.DNSFlushName == "" || len(g.opt.DNSFlushArgs) == 0 {
		if runtime.GOOS == "linux" {
			return exec.Command("systemctl", "restart", "NetworkManager.service").Run()
		}
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "ipconfig", "/flushdns").Run()
		}
		return fmt.Errorf("%s not support", runtime.GOOS)
	}
	return exec.Command(g.opt.DNSFlushName, g.opt.DNSFlushArgs...).Run()
}
