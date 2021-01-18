package main

import "testing"

func TestParseCfg(t *testing.T) {
	parseCfg()
}

func TestSaveCfg(t *testing.T){
	saveCfg()
}

func TestGetUsrHomeDir(t *testing.T) {
 homeDir:=getUsrHomeDir()
 if	homeDir != `C:\Users\dimon`{
	 t.Error("wrong homedir")
 }
}
