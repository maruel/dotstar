// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import (
	"errors"
	"strings"
	"sync"
	"unicode/utf8"
)

// LocalBus is a Bus implementation that runs locally in the process.
// http://www.hivemq.com/blog/mqtt-essentials-part-8-retained-messages
type LocalBus struct {
	mu               sync.Mutex
	persistentTopics map[string][]byte
	subscribers      []*subscription
}

func (l *LocalBus) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range l.subscribers {
		l.subscribers[i].Close()
	}
	return nil
}

func (l *LocalBus) Publish(msg Message, qos QOS, retained bool) error {
	p := parseTopic(msg.Topic)
	if p == nil || p.isQuery() {
		return errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.persistentTopics == nil {
		l.persistentTopics = map[string][]byte{}
	}
	if len(msg.Payload) == 0 {
		// delete
		delete(l.persistentTopics, msg.Topic)
		return nil
	}
	// Save it first.
	if retained {
		l.persistentTopics[msg.Topic] = msg.Payload
	}
	for i := range l.subscribers {
		if l.subscribers[i].topic.match(msg.Topic) {
			l.subscribers[i].publish(msg)
		}
	}
	return nil
}

func (l *LocalBus) Subscribe(topic string, qos QOS) (<-chan Message, error) {
	p := parseTopic(topic)
	if p == nil {
		return nil, errors.New("invalid topic")
	}
	c := make(chan Message)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subscribers = append(l.subscribers, &subscription{topic: p, channel: c})
	return c, nil
}

func (l *LocalBus) Unsubscribe(topic string) error {
	p := parseTopic(topic)
	if p == nil {
		return errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range l.subscribers {
		if l.subscribers[i].topic.isEqual(p) {
			l.subscribers[i].Close()
			copy(l.subscribers[i:], l.subscribers[i+1:])
			l.subscribers = l.subscribers[:len(l.subscribers)-1]
			return nil
		}
	}
	return errors.New("subscription not found")
}

func (l *LocalBus) Get(topic string, qos QOS) ([]Message, error) {
	p := parseTopic(topic)
	if p == nil {
		return nil, errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Message
	for k, v := range l.persistentTopics {
		if p.match(k) {
			out = append(out, Message{k, v})
		}
	}
	return out, nil
}

/*
// Settle waits for all messages in transit to be fully published.
func (l *LocalBus) Settle() {
	l.mu.Lock()
	s := l.subscribers
	l.mu.Unlock()
	for i := range s {
		l.subscribers[i].settle()
	}
}
*/

//

type subscription struct {
	topic parsedTopic
	wg    sync.WaitGroup

	mu      sync.Mutex
	channel chan Message
}

func (s *subscription) publish(msg Message) {
	s.mu.Lock()
	s.mu.Unlock()
	if s.channel == nil {
		return
	}
	// s.wg.Add() and s.wg.Wait() can only be called with the lock held.
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.channel <- msg
	}()
}

/*
func (s *subscription) settle() {
	// TODO(maruel): Add(1) must be called in the same goroutine as Wait(). Dang.
	// Since publishing can happend in goroutine, this doesn't work, as it
	// doesn't make sense to grab the lock here.
	s.wg.Wait()
}
*/

func (s *subscription) Close() {
	s.mu.Lock()
	s.mu.Unlock()
	for ok := true; ok; {
		select {
		case _, ok = <-s.channel:
		default:
			ok = false
		}
	}
	// TODO(maruel): Some even may not be able to complete as the lock is held.
	// s.wg.Wait()
	// TODO(maruel): This may blow up the hung publishers.
	// close(s.channel)
	// s.channel = nil
}

// parsedTopic is either a query or a static topic.
type parsedTopic []string

func parseTopic(topic string) parsedTopic {
	if len(topic) == 0 || len(topic) > 65535 || strings.ContainsRune(topic, rune(0)) || !utf8.ValidString(topic) {
		return nil
	}
	p := parsedTopic(strings.Split(topic, "/"))
	if !p.isValid() {
		return nil
	}
	return p
}

func (p parsedTopic) isValid() bool {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180921
	// section 4.7.2 about '$' prefix and section 4.7.3
	if len(p[0]) != 0 && p[0][0] == '$' {
		return false
	}
	for i, e := range p {
		// As per the spec, empty sections are valid.
		if i != len(p)-1 && e == "#" {
			// # can only appear at the end.
			return false
		} else if e != "+" && e != "#" {
			if strings.HasSuffix(e, "#") || strings.HasSuffix(e, "+") {
				return false
			}
		}
	}
	return true
}

func (p parsedTopic) isQuery() bool {
	for _, e := range p {
		if e == "#" || e == "+" {
			return true
		}
	}
	return false
}

func (p parsedTopic) isEqual(other parsedTopic) bool {
	if len(other) != len(p) {
		return false
	}
	for i := range p {
		if p[i] != other[i] {
			return false
		}
	}
	return true
}

// match follows rules as defined at section 4.7:
// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180919
func (p parsedTopic) match(topic string) bool {
	t := strings.Split(topic, "/")
	if len(t) == 0 {
		return false
	}
	for i, e := range p {
		if e == "#" {
			return true
		}
		if e == "+" {
			if i == len(p)-1 && len(t) == len(p) || len(t) == len(p)-1 {
				return true
			}
			continue
		}
		if len(t) <= i || t[i] != e {
			return false
		}
	}
	return len(t) == len(p)
}

var _ Bus = &LocalBus{}
