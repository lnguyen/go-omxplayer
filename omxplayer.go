package omxplayer

import dbus "github.com/hoffoo/go.dbus"

var sdbus *dbus.Object

func connDbus() *dbus.Object {

	conn, err := dbus.SessionBus()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2")
}

func Status() *dbus.Variant {

	pstatus, err := sdbus.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")

	if err != nil {
		panic(err) // most likely spotify not running
	}

	return &pstatus
}

func OmxPlayerMethod(method string) error {
	if sdbus == nil {
		sdbus = connDbus()
	}

	c := sdbus.Call(method, 0)
	return c.Err
}
