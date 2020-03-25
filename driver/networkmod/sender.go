package networkmod

import (
	//"fmt"
	//"time"
	"../config"
)


func SendData(ch config.NetworkChannels, newOrder <-chan config.ElevatorOrder, newState <-chan [config.NumElevators]config.ElevatorState) {
	for {
		select{
		case orderMsg := <- newOrder:
			ch.TransmittOrderCh <- orderMsg
			//time.Sleep(1 * time.Second)//må endre tid, sikkert sende meldinger mye oftere
		case stateMsg := <- newState:
			ch.TransmittStateCh <- stateMsg
		
		}	
	}	
}