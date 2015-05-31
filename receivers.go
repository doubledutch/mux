package mux

import (
	"os"
	"reflect"
	"syscall"
)

// ValueReceiver uses reflection to send and receive values
type ValueReceiver struct {
	dec BufferDecoder
	ch  reflect.Value
	t   reflect.Type
}

// NewReceiver creates a new Receiver.
func NewReceiver(ch interface{}, pool Pool) Receiver {
	v := reflect.TypeOf(ch)

	if v.Kind() != reflect.Chan {
		panic("Receiver requires a channel")
	}

	chV := reflect.ValueOf(ch)
	t := reflect.TypeOf(ch).Elem()

	return &ValueReceiver{
		dec: pool.NewBufferDecoder(),
		ch:  chV,
		t:   t,
	}
}

func (r *ValueReceiver) Receive(b []byte) error {
	x := reflect.New(r.t)

	r.dec.Write(b)
	err := r.dec.Decode(x.Interface())
	r.dec.Reset()
	if err != nil {
		return err
	}

	r.ch.Send(x.Elem())
	return nil
}

func (r *ValueReceiver) Close() error {
	r.ch.Close()
	return nil
}

// SignalReceiver receives signals
type SignalReceiver struct {
	dec BufferDecoder
	ch  chan os.Signal
}

// NewSignalReceiver creates a new signal receiver
func NewSignalReceiver(ch chan os.Signal, pool Pool) SignalReceiver {
	return SignalReceiver{
		dec: pool.NewBufferDecoder(),
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
	dec BufferDecoder
	ch  chan string
}

// NewStringReceiver returns a StringReceiver
func NewStringReceiver(ch chan string, pool Pool) StringReceiver {
	return StringReceiver{
		dec: pool.NewBufferDecoder(),
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
