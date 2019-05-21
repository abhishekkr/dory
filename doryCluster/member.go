package doryCluster

import (
	"fmt"
	"strings"

	golconv "github.com/abhishekkr/gol/golconv"
	golenv "github.com/abhishekkr/gol/golenv"
	gollog "github.com/abhishekkr/gol/gollog"
	"github.com/hashicorp/memberlist"
)

var (
	memberlistBindPort      = golenv.OverrideIfEnv("DORY_MEMBERS_BIND", "7946")
	memberlistAdvertisePort = golenv.OverrideIfEnv("DORY_MEMBERS_ADVERTISE", "7946")
	members                 *memberlist.Memberlist
)

func init() {
	var err error
	memberlistCfg := memberlist.DefaultLocalConfig()
	memberlistCfg.BindPort = golconv.StringToInt(memberlistBindPort, 7946)
	memberlistCfg.AdvertisePort = golconv.StringToInt(memberlistAdvertisePort,
		memberlistCfg.BindPort)
	members, err = memberlist.Create(memberlistCfg)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
}

func Join(leadersIP string) {
	if leadersIP == "" {
		gollog.Debug("no leaders to join")
		return
	}
	ipList := strings.Split(leadersIP, ",")
	n, err := members.Join(ipList)
	if err != nil {
		gollog.Err(fmt.Sprintf("Failed to join cluster using list '%s': %s",
			leadersIP, err.Error()))
	}
	gollog.Info(fmt.Sprintf("joined %d leaders from cluster: %s", n, leadersIP))
}

func Members() {
	for _, member := range members.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}
}
