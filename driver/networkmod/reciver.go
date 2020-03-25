package networkmod

import (
	"fmt"
	"../config"
	//"reflect"
	//"time"
	//"sync"
	"strconv"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.



func RecieveData(id int, ch config.NetworkChannels, elevatorList *[config.NumElevators]config.ElevatorState,
	activeElevators *[config.NumElevators]bool) {
	iter := 1
	
	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			
			//If recievig from a new peer, upadate avtive elevator map
			ch.TransmittStateCh <- *elevatorList
			
			for _, peer := range p.New{
				peerId, _ := strconv.Atoi(peer)
				activeElevators[peerId] = true
				ch.TransmittStateCh <- *elevatorList
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					activeElevators[peerId] = false
				}
			}
			

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			fmt.Println("Iteration %v", iter)
			for i, elevatorState := range newState{
				if (i != id && elevatorState.Id != id){
					elevatorList[i] = elevatorState
				}
				fmt.Println("Id: ", i)
				fmt.Println(elevatorList[i])
			}
			iter++

		case newOrder := <-ch.RecieveOrderCh:
			fmt.Println(newOrder)
			id := newOrder.ExecutingElevator
			elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)

		}

	}
}