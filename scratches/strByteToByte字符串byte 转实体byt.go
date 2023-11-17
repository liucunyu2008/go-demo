package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	b:="safdasfdafasfd"
	j,err:=json.Marshal(b)
	if err!=nil{
		fmt.Printf("Error marshalling:%#v\n",err)
		return
	}
	fmt.Printf("===================j=============:%v\n",j)
	a:="[34 115 97 102 100 97 115 102 100 97 102 97 115 102 100 34]"
	fmt.Printf("===retstr==:%#v\n",string(strByteToByte(a)))

}


func strByteToByte(sb string)[]byte  {
	var bb []byte
	ps:=strings.Split(strings.Trim(sb, "[]"), " ")
	for _, v:=range ps  {
		pi,_:=strconv.Atoi(v)
		fmt.Printf("==========================pi======:%#v\n",pi)
		bb=append(bb,byte(pi))
	}
	return bb
}
