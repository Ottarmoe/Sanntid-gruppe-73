package main 

type hallOrderState int

const (
	HallO hallOrderState = iota
	HallNO
	HallOPR
)

var hallOrders [4][2]hallOrderState

type cabOrderState int

const (
	CabO cabOrderState = iota
	CabNO
	CabUO //unknown order
)

var cabOrders [4]cabOrderState

type State struct {
    HallOrders    	[4][2]hallOrderState 
    CabOrders     	[4]cabOrderState      
    CabFloor   	 	string    
    CabDir      
	CabBehavoiur
	CabMechError
	OnNetwork		bool
}