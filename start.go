package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"unsafe"
)

const (
	HOST="0.0.0.0"
	PORT = "179"
	TYPE = "tcp"
)

func main() {
	// Listen for incoming connections
	l, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error opening Port",err.Error())
	}
	// Close the socket when the application closes.
	defer l.Close()
	fmt.Println("BGP Daemon Running" + HOST + ":" + PORT)
	for {
		connection, err :=l.Accept()
		if(err != nil) {
			fmt.Println("Error Accepting:", err.Error())
			os.Exit(1)
		}
		go handleRequest(connection)

	}
}
type BGPOPEN struct {
	bgp_type int
	bgp_version int
	bgpAS int

}
func handleRequest(conn net.Conn){
	buf := make([]byte,1024) // Make a buffer to hold incoming data
	// Read the incoming connection into buffer
	reqLen, err := conn.Read(buf)
	_ = reqLen
	fmt.Println(buf[:reqLen])
	//bgpLen := buf[17]
	bgpType := buf[18]
	bgpVersion := buf[19]
	bgpAS := ByteArrayToInt(buf[20:22])
	//bgpHT := buf[24]
	bgpNeighborID := strconv.Itoa(int(buf[24]))+"."+strconv.Itoa(int(buf[25]))+"."+strconv.Itoa(int(buf[26]))+"."+strconv.Itoa(int(buf[27]))
	if(bgpType==1) {
		fmt.Println("BGP OPEN Message received. NeighborID " + string(bgpNeighborID) + " Type: " + strconv.Itoa(int(bgpType)) + ", BGP Version: " + strconv.Itoa(int(bgpVersion)) + ",ASN: " + strconv.Itoa(int(bgpAS)))
		buf[18] = 4
		buf[24] = 10
		buf[25] = 0
		buf[26] = 1
		buf[27] = 100
		conn.Write(buf)
		fmt.Println(buf)
	}
	if err != nil {
		fmt.Println("Error Reading:",err.Error())
	}

	conn.Close()
}


func ByteArrayToInt(arr []byte) int64{
	val := int64(0)
	size := len(arr)
	for i := 0 ; i < size ; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}