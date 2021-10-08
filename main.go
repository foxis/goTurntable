package main

//go:generate tinygo build -target=d1mini -o table.bin
//go:generate esptool -p COM10 write_flash 0x0 table.bin

import (
	"machine"
	"time"
)

const (
	SPEED_STEP = time.Microsecond * 900
	MAX_SPEED  = 8
)

var (
	speed int = 0
	steps     = [][4]bool{
		{true, false, false, false},
		{true, true, false, false},
		{false, true, false, false},
		{false, true, true, false},
		{false, false, true, false},
		{false, false, true, true},
		{false, false, false, true},
		{true, false, false, true},
	}
)

type Stepper struct {
	p1, p2, p3, p4 machine.Pin
	step           int
}

func (s *Stepper) Configure() {
	s.p1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s.p2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s.p3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s.p4.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (s *Stepper) Step(d int) {
	s.step += d
	if s.step > 0 {
		s.step = s.step % len(steps)
	} else if s.step < 0 {
		s.step = len(steps) + s.step
	}

	s.p1.Set(steps[s.step][0])
	s.p2.Set(steps[s.step][1])
	s.p3.Set(steps[s.step][2])
	s.p4.Set(steps[s.step][3])
}

func (s *Stepper) Off() {
	s.p1.Low()
	s.p2.Low()
	s.p3.Low()
	s.p4.Low()
}

type Button struct {
	pin  machine.Pin
	prev bool
	f    func(bool)
}

func (b *Button) Configure() {
	b.pin.Configure(machine.PinConfig{Mode: machine.PinInput})
}

func (b *Button) Update() {
	val := b.pin.Get()
	if val != b.prev {
		b.f(val)
	}
	b.prev = val
}

func main() {
	leftBtn := Button{
		machine.D1,
		false,
		func(v bool) {
			if v {
				speed += 1
			}
		},
	}
	rightBtn := Button{
		machine.D2,
		false,
		func(v bool) {
			if v {
				speed -= 1
			}
		},
	}

	leftBtn.Configure()
	rightBtn.Configure()
	motor := Stepper{machine.D5, machine.D6, machine.D7, machine.D8, 0}
	motor.Configure()

	i := 0
	for {
		leftBtn.Update()
		rightBtn.Update()
		if speed > MAX_SPEED {
			speed = MAX_SPEED
		}
		if speed < -MAX_SPEED {
			speed = -MAX_SPEED
		}

		d := 0
		sp := speed
		if sp < 0 {
			sp = -sp
		}

		if i > MAX_SPEED-sp {
			if speed > 0 {
				d = 1
			} else {
				d = -1
			}
			i = 0
		}
		if sp == 0 {
			d = 0
			motor.Off()
		} else {
			motor.Step(d)
		}

		i++
		time.Sleep(SPEED_STEP)
	}
}
