package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"

	"github.com/holimon/ghosts"
	"github.com/urfave/cli"
)

func recordAdd(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	if len(c.Args()) != 2 {
		return fmt.Errorf("the number of arguments for this instruction must be equal to 2")
	}
	if net.ParseIP(c.Args()[0]) != nil {
		h.Add(c.Args()[1], c.Args()[0])
	}
	if net.ParseIP(c.Args()[1]) != nil {
		h.Add(c.Args()[0], c.Args()[1])
	}
	return h.WriteBack()
}

func deleteBycontent(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	if len(c.Args()) < 1 {
		return fmt.Errorf("the number of arguments for this instruction must be greater than 1")
	}
	for _, arg := range c.Args() {
		h.Delete(arg)
	}
	return h.WriteBack()
}

func deleteByindex(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	records := h.OrdedRecords(ghosts.SortedByDomain)
	if len(c.Args()) < 1 {
		return fmt.Errorf("the number of arguments for this instruction must be greater than 1")
	}
	for _, arg := range c.Args() {
		if k, e := strconv.ParseInt(arg, 10, 32); e == nil && int(k) < len(records) && int(k) >= 0 {
			h.DeleteRecord(records[k])
		}
	}
	return h.WriteBack()
}

func showRecord(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	for k, item := range h.OrdedRecords(ghosts.SortedByDomain) {
		fmt.Printf("Index: %-4d IP: %-48s Domain: %s\n", k, item.IP, item.Domain)
	}
	return nil
}

func replaceRecord(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	if len(c.Args()) != 2 {
		return fmt.Errorf("the number of arguments for this instruction must be equal to 2")
	}
	if net.ParseIP(c.Args()[0]) != nil {
		if e := h.Replace(c.Args()[1], c.Args()[0]); e != nil {
			return e
		}
		return h.WriteBack()
	}
	if net.ParseIP(c.Args()[1]) != nil {
		if e := h.Replace(c.Args()[0], c.Args()[1]); e != nil {
			return e
		}
		return h.WriteBack()
	}
	return fmt.Errorf("arguments error")
}

func resolveDomain(c *cli.Context) error {
	h, err := ghosts.New(ghosts.Option{})
	if err != nil {
		return err
	}
	if len(c.Args()) < 1 {
		return fmt.Errorf("the number of arguments for this instruction must be greater than 1")
	}
	for _, arg := range c.Args() {
		if addrs, err := h.Resolve(arg); err == nil {
			h.Delete(arg)
			for _, ip := range addrs {
				h.Add(arg, ip)
			}
		} else {
			fmt.Println(err)
		}
	}
	return h.WriteBack()
}

func main() {
	app := cli.NewApp()
	app.Usage = "A tool for adding, deleting, and replacing hosts records."
	app.Version = "v0.0.1"
	app.Commands = []cli.Command{{Name: "add", Usage: "Add a record to the hosts file. Receive 2 arguments.", Action: recordAdd},
		{Name: "del", Usage: "Remove records from the hosts file.", Subcommands: []cli.Command{
			{Name: "field", Usage: "Remove records by Domain or IP. Receive at least 1 argument.", Action: deleteBycontent},
			{Name: "index", Usage: "Remove records by index. Receive at least 1 argument.", Action: deleteByindex},
		}},
		{Name: "replace", Usage: "Replace a record from the hosts file. Receive 2 arguments.", Action: replaceRecord},
		{Name: "resolve", Usage: "Resolve domains and add records to the hosts file. Receive at least 1 argument.", Action: resolveDomain},
		{Name: "show", Usage: "Show all records in the hosts file. No arguments.", Action: showRecord}}
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
