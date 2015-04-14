/*
Copyright 2015 Doubledutch

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mux

import (
	"os"
	"testing"
)

func TestStringReceiver(t *testing.T) {
	strCh := make(chan string, 1)
	expected := "hello world"

	strR := NewStringReceiver(strCh)
	defer strR.Close()

	enc := NewEncoder()
	enc.Encode(expected)

	if err := strR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	actual := <-strCh
	if actual != expected {
		t.Fatalf("actual '%s' != expected '%s'", actual, expected)
	}

}

func TestSignalReceiver(t *testing.T) {
	sigCh := make(chan os.Signal, 1)

	sigR := NewSignalReceiver(sigCh)
	defer sigR.Close()

	enc := NewEncoder()
	enc.Encode(os.Kill)

	if err := sigR.Receive(enc.Bytes()); err != nil {
		t.Fatal(err)
	}

	if <-sigCh != os.Kill {
		t.Fatal("expected os.Kill")
	}
}
