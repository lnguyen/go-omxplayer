package omxplayer

import (
	"io/ioutil"
	"os"
	"os/exec"

	dbus "github.com/hoffoo/go.dbus"
)

var sdbus *dbus.Object

type OmxPlayer struct {
	sdbus   *dbus.Object
	command *exec.Cmd
}

func New() OmxPlayer {
	return OmxPlayer{}
}

func (o *OmxPlayer) PlayFile(filename string) error {
	o.command = exec.Command("omxplayer", "-o", "hdmi", filename)
	err := o.command.Start()
	if err != nil {
		return err
	}
	env, err := ioutil.ReadFile("/tmp/omxplayerdbus")
	if err != nil {
		return err
	}
	os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(env))
	return nil
}

func (o *OmxPlayer) PlayPause() error {
	err := o.Method("org.mpris.MediaPlayer2.Player.PlayPause")
	return err
}

func (o *OmxPlayer) Method(method string) error {
	o.sdbus = connDbus()

	c := o.sdbus.Call(method, 0)
	return c.Err
}

func connDbus() *dbus.Object {
	conn, err := dbus.SessionBus()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2")
}
