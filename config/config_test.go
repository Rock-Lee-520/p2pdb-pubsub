package config_test

import (
	"testing"

	"github.com/Rock-liyi/p2pdb-log/config"
	debug "github.com/favframework/debug"
	_ "github.com/joho/godotenv/autoload"
)

func TestDatabase_Name(t *testing.T) {
	//require := require.New(t)
	path := config.GetDataPath()
	env := config.GetEnv()
	isdebu := config.IsDebug()
	debug.Dump(path)
	debug.Dump(env)
	debug.Dump(isdebu)
	//require.Emptyf(path)
	// val, ex := os.LookupEnv("DATAPATH")
	// debug.Dump(val)
	// debug.Dump(ex)
}
