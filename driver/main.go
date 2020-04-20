package main

import (
  "./elevio"
  "./fsm"
  cf "./config"
  oh "./orderhandler"
  "./networkmod"
  "./networkmod/network/peers"
  "./networkmod/network/bcast"
  "flag"
  "strconv"
)

func main(){
  
  var (
    localHostId        string
    id                 int
    ElevatorList       [cf.NumElevators]cf.ElevatorState
    ActiveElevatorList [cf.NumElevators]bool
  )

    flag.StringVar(&localHostId, "hostId", "", "hostId of this peer")
    flag.IntVar(&id, "id", 1234, "id of this peer")
    flag.Parse()
    elevio.Init("localhost:"+localHostId, cf.NumFloors)

    idAsString := strconv.Itoa(id)


    driverChannels := cf.DriverChannels{
      DrvButtons: make(chan elevio.ButtonEvent),
      DrvFloors: make(chan int),
      DrvStop: make(chan bool),
    }

    networkChannels := cf.NetworkChannels{
      PeerTxEnable: make(chan bool),
      PeerUpdateCh: make(chan peers.PeerUpdate),
      TransmittOrderCh: make(chan cf.ElevatorOrder),
      TransmittStateCh: make(chan map[string][cf.NumElevators]cf.ElevatorState),
      RecieveOrderCh: make(chan cf.ElevatorOrder),
      RecieveStateCh: make(chan map[string][cf.NumElevators]cf.ElevatorState),
    }

    timerChannels := cf.TimerChannels{
      Open_door: make(chan bool),
    }


    orderChannels := cf.OrderChannels{
      LostConnection: make(chan cf.ElevatorState),
      SendState: make(chan map[string][cf.NumElevators]cf.ElevatorState),
      SendOrder: make(chan cf.ElevatorOrder),
    }

    go peers.Transmitter(22349, idAsString, networkChannels.PeerTxEnable)
    go peers.Receiver(22349, networkChannels.PeerUpdateCh)
    go bcast.Transmitter(22367, networkChannels.TransmittOrderCh)
    go bcast.Receiver(22367, networkChannels.RecieveOrderCh)
    go bcast.Transmitter(22378, networkChannels.TransmittStateCh)
    go bcast.Receiver(22378, networkChannels.RecieveStateCh)

    go elevio.PollButtons(driverChannels.DrvButtons)
    go elevio.PollFloorSensor(driverChannels.DrvFloors)
    go elevio.PollStopButton(driverChannels.DrvStop)

    go networkmod.RecieveData(id, networkChannels, orderChannels.LostConnection, &ElevatorList, &ActiveElevatorList)
    go networkmod.SendData(networkChannels, orderChannels) 

    go oh.OrderHandler(driverChannels.DrvButtons, orderChannels.SendOrder, orderChannels.SendState, networkChannels.RecieveStateCh, 
      networkChannels.RecieveOrderCh, orderChannels.LostConnection, id, &ElevatorList, &ActiveElevatorList)
    
    fsm.ElevStateMachine(driverChannels, id, orderChannels.SendOrder, orderChannels.SendState, &ElevatorList,
         timerChannels, orderChannels.LostConnection)

}