package main

import (
	"flag"
	"strings"

	"github.com/kiwicom/sheriff"
)

func parseMembers(members string) []string {
	addrs := strings.Split(members, ",")
	if members == "" {
		return nil
	}
	result := []string{}
	for _, addr := range addrs {
		result = append(result, addr)
	}
	return result
}

func main() {
	//flag.String("service", "", "service name")
	//flag.Int("port", 0, "sheriff port")
	members := flag.String("members", "", "comma separeted list of known members(ip:port)")
	flag.Parse()

	existing := parseMembers(*members)
	localMember, err := sheriff.NewMember(existing)
	if err != nil {
		panic(err)
	}
	localMember.Run()
}
