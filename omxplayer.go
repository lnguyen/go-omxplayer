package omxplayer

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	dbus "github.com/hoffoo/go.dbus"
)

var sdbus *dbus.Object

type OmxPlayer struct {
	sdbus    *dbus.Object
	command  *exec.Cmd
	filename string
}

func New() OmxPlayer {
	return OmxPlayer{}
}

func (o *OmxPlayer) IsPlaying() bool {
	if o.command != nil {
		return true
	}
	return false
}

func (o *OmxPlayer) FilePlaying() string {
	if o.filename != "" {
		return o.filename
	}
	return ""
}

func (o *OmxPlayer) PlayFile(filename string) error {
	if o.IsPlaying() {
		return errors.New("Error file is playing, please stop and try again")
	}
	o.command = exec.Command("omxplayer", "-b", "-o", "local", "--loop", filename)
	o.filename = filename
	err := o.command.Start()
	if err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond)
	env, err := ioutil.ReadFile("/tmp/omxplayerdbus")
	if err != nil {
		return err
	}
	os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(env))
	return nil
}

func (o *OmxPlayer) StopFile() error {
	o.command.Process.Kill()
	o.filename = ""
	err := exec.Command("killall", "omxplayer.bin").Run()
	return err
}

//@TODO DOESNT WORK ATM :(
func (o *OmxPlayer) PlayPause() error {
	err := o.Method("org.mpris.MediaPlayer2.omxplayer.Player.PlayPause")
	return err
}

func (o *OmxPlayer) Method(method string) error {
	o.sdbus = connDbus()

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
