package omxplayer

import (
	"os/exec"

	dbus "github.com/hoffoo/go.dbus"
)

var sdbus *dbus.Object

type OmxPlayer struct {
	sdbus   *dbus.Object
	command exec.Cmd
}

func (o *OmxPlayer) PlayFile(filename string) error {
	o.command = exec.Command("omxplayer", "-o", "hdmi", filename)
	err := o.command.Run()
	if err != nil {
		return err
	}
	return nil
}

func (o *OmxPlayer) PlayPause() error {
	err := o.Method("PlayPause")
	return err
}

func (o *OmxPlayer) Method(method string) error {
	if o.sdbus == nil {
		o.sdbus = connDbus()
	}

	c := o.sdbus.Call(method, 0)
	return c.Err
}

func connDbus() *dbus.Object {
	conn, err := dbus.SessionBusPrivate()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2")
}
