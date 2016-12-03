/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : node.go

* Purpose :

* Creation Date : 02-23-2015

* Last Modified : Tue 10 Mar 2015 07:32:22 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Node struct {
	Id          int64  `json:"id"`
	Hostname    string `json:"hostname"`
	ManageIp    string `json:"manage_ip"`
	InterfaceIp string `json:"interface_ip"`
	Group       string `json:"group"`
}

type NodeWriteOptions struct {
	NoCommit  bool
	NoPush    bool
	CommitMsg string
}

func (n *Node) Write(options *NodeWriteOptions) error {
	var channels []Channel
	err := x.DB.Where("groups LIKE ?", "%\""+n.Group+"\"%").Find(&channels)
	if err != nil {
		return err
	}
	for _, c := range channels {
		t0 := reListenInterface443.ReplaceAllString(c.Config, "listen              "+n.InterfaceIp+":443 ssl")
		t1 := reListenInterface80.ReplaceAllString(t0, "listen              "+n.InterfaceIp+":80")
		dir := fmt.Sprintf("nginxrepo/nginx_%s/%s", n.Group, n.Hostname)
		if _, err := os.Stat(dir); err != nil {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return err
			}
		}
		filename := fmt.Sprintf("%s%s/%s", myDir, dir, c.Domain)
		err = ioutil.WriteFile(filename, []byte(t1), 0644)
		if err != nil {
			return err
		}
		err := Groups[n.Group].Add(filename)
		if err != nil {
			return err
		}
	}

	if options.CommitMsg == "" {
		options.CommitMsg = "commit all"
	}

	if !options.NoCommit {
		err = Groups[n.Group].CommitAll(options.CommitMsg)
		if err != nil {
			return err
		}
	}

	if !options.NoPush {
		Groups[n.Group].Push()
	}

	return nil
}
