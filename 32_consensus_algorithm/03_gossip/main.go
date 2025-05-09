package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"time"
)

func main() {
	/* Create the initial memberlist from a safe configuration.
	   Please reference the godoc for other default config types.
	   http://godoc.org/github.com/hashicorp/memberlist#Config
	*/
	list, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		fmt.Printf("Failed to create memberlist: " + err.Error())
		return
	}

	t := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-t.C:
			// Join an existing cluster by specifying at least one known member.
			n, err := list.Join([]string{"127.0.0.1"})
			if err != nil {
				fmt.Println("Failed to join cluster: " + err.Error())
				continue
			}
			fmt.Println("member number is:", n)
			goto END
		}
	}
END:
	for {
		select {
		case <-t.C:
			// Ask for members of the cluster
			for _, member := range list.Members() {
				fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
			}
		}
	}

	// Continue doing whatever you need, memberlist will maintain membership
	// information in the background. Delegates can be used for receiving
	// events when members join or leave.
}
