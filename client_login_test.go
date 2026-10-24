package routeros

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"testing"
	"time"
)

func loginTestClient(t *testing.T) (*Client, net.Conn) {
	t.Helper()
	cfg := NewConfig().
		WithCredentials("admin", "secret").
		WithTimeout(time.Second).
		WithReadTimeout(time.Second).
		WithWriteTimeout(time.Second).
		WithRetry(0, 0)
	clientConn, serverConn := net.Pipe()
	client := NewClient(cfg)
	client.conn = clientConn
	client.protocol = NewProtocol(clientConn, clientConn)
	t.Cleanup(func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
	})
	return client, serverConn
}

func checkLoginSentence(sentence *Sentence, withPassword bool) error {
	if sentence == nil || len(sentence.Words) != 3 {
		return fmt.Errorf("login sentence should have 3 words")
	}
	if sentence.Words[0].Value != "/login" {
		return fmt.Errorf("first word = %q, want /login", sentence.Words[0].Value)
	}
	if sentence.Words[1].Key != "name" || sentence.Words[1].Value != "admin" {
		return fmt.Errorf("username attribute = %#v, want admin", sentence.Words[1])
	}
	wantKey := "response"
	if withPassword {
		wantKey = "password"
	}
	if sentence.Words[2].Key != wantKey {
		return fmt.Errorf("third attribute key = %q, want %q", sentence.Words[2].Key, wantKey)
	}
	if withPassword && sentence.Words[2].Value != "secret" {
		return fmt.Errorf("password attribute was not sent in the first login sentence")
	}
	return nil
}

func TestLoginSendsCredentialsBeforeReadingResponse(t *testing.T) {
	client, serverConn := loginTestClient(t)
	server := NewProtocol(serverConn, serverConn)
	serverDone := make(chan error, 1)

	go func() {
		sentence, err := server.ReadSentence()
		if err != nil {
			serverDone <- err
			return
		}
		if err := checkLoginSentence(sentence, true); err != nil {
			serverDone <- err
			return
		}
		serverDone <- server.WriteSentence("!done")
	}()

	if err := client.login(); err != nil {
		t.Fatalf("login() error = %v", err)
	}
	if err := <-serverDone; err != nil {
		t.Fatalf("test server: %v", err)
	}
}

func TestLoginFallsBackToLegacyChallengeResponse(t *testing.T) {
	client, serverConn := loginTestClient(t)
	server := NewProtocol(serverConn, serverConn)
	challengeHex := "00112233445566778899aabbccddeeff"
	challenge, _ := hex.DecodeString(challengeHex)
	serverDone := make(chan error, 1)

	go func() {
		first, err := server.ReadSentence()
		if err != nil {
			serverDone <- err
			return
		}
		if err := checkLoginSentence(first, true); err != nil {
			serverDone <- err
			return
		}
		if err := server.WriteSentence("!done", "=ret="+challengeHex); err != nil {
			serverDone <- err
			return
		}

		second, err := server.ReadSentence()
		if err != nil {
			serverDone <- err
			return
		}
		if second == nil || len(second.Words) != 3 ||
			second.Words[0].Value != "/login" ||
			second.Words[1].Key != "name" || second.Words[1].Value != "admin" ||
			second.Words[2].Key != "response" {
			serverDone <- fmt.Errorf("unexpected legacy login response: %#v", second)
			return
		}

		hasher := md5.New()
		_, _ = hasher.Write([]byte{0})
		_, _ = hasher.Write([]byte("secret"))
		_, _ = hasher.Write(challenge)
		wantResponse := "00" + hex.EncodeToString(hasher.Sum(nil))
		if second.Words[2].Value != wantResponse {
			serverDone <- fmt.Errorf("legacy response = %q, want %q", second.Words[2].Value, wantResponse)
			return
		}
		serverDone <- server.WriteSentence("!done")
	}()

	if err := client.login(); err != nil {
		t.Fatalf("legacy login() error = %v", err)
	}
	if err := <-serverDone; err != nil {
		t.Fatalf("test server: %v", err)
	}
}
