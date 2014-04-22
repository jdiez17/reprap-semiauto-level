package main

import (
	"strconv"
)

type GCodeCommand interface {
	ToString() string
}

type BasicGCodeCommand struct {
	command string
}

type MoveToOriginCommand struct {
	BasicGCodeCommand
	Axes []string
}

type MoveToCommand struct {
	BasicGCodeCommand
	Positions map[string]float64
}

func (c BasicGCodeCommand) ToString() string {
	return c.command
}

func (c MoveToOriginCommand) ToString() string {
	c.command = "G28"

	for _, axis := range c.Axes {
		c.command += " " + axis + "0"
	}

	return c.command
}

func (c MoveToCommand) ToString() string {
	c.command = "G1"

	for axis, pos := range c.Positions {
		c.command += " " + axis + strconv.FormatFloat(pos, 'f', 10, 64)
	}

	return c.command
}

// Convenience
var SetCurrentLineNumber = BasicGCodeCommand{command: "M110"}
var GetTemperature = BasicGCodeCommand{command: "M105"}
var SetAbsolutePositioning = BasicGCodeCommand{command: "G90"}
var SetRelativePositioning = BasicGCodeCommand{command: "G91"}
var HomeAllAxes = MoveToOriginCommand{Axes: []string{"X", "Y", "Z"}}

var MoveXY = func(X, Y float64) MoveToCommand {
	return MoveToCommand{
		Positions: map[string]float64{"X": X, "Y": Y},
	}
}
var MoveZ = func(Z float64) MoveToCommand {
	return MoveToCommand{
		Positions: map[string]float64{"Z": Z},
	}
}
