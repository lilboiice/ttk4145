package networkmod

import (
	"fmt"
	"../config"
	//"reflect"
	//"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.



func RecieveData(id string, ch config.NetworkChannels, reciever chan<- []string, recievedNewStateUpdateFromNetwork chan<- map[string]config.ElevatorState) {


	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			reciever <- p.New


		case a := <-ch.RecieveStateCh:
			recievedNewStateUpdateFromNetwork <- a
		}
	}
}
