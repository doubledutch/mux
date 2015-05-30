package mux

import (
	"os"
	"syscall"
)

// SignalReceiver receives signals
type SignalReceiver struct {
	dec *Decoder
	ch  chan os.Signal
}

// NewSignalReceiver creates a new signal receiver
func NewSignalReceiver(ch chan os.Signal) SignalReceiver {
	return SignalReceiver{
		dec: NewDecoder(),
		ch:  ch,
	}
}

// Receive decodes bytes into signal and puts it on ch
func (r SignalReceiver) Receive(b []byte) error {
	var sig syscall.Signal

	r.dec.Write(b)
	err := r.dec.Decode(&sig)
	r.dec.Reset()
	if err != nil {
		return err
	}
	r.ch <- sig
	return nil
}

// Close and cleans up SignalReceiver
func (r SignalReceiver) Close() error {
	close(r.ch)
	return nil
}

// StringReceiver receives strings
type StringReceiver struct {
	dec *Decoder
	ch  chan string
}

// NewStringReceiver returns a StringReceiver
func NewStringReceiver(ch chan string) StringReceiver {
	return StringReceiver{
		dec: NewDecoder(),
		ch:  ch,
	}
}

// Receive decodes bytes into string and puts it on ch
func (r StringReceiver) Receive(b []byte) error {
	var str string

	r.dec.Write(b)
	err := r.dec.Decode(&str)
	r.dec.Reset()
	if err != nil {
		return err
	}
	r.ch <- str
	return nil
}

// Close and cleans up StringReceiver
func (r StringReceiver) Close() error {
	close(r.ch)
	return nil
}
