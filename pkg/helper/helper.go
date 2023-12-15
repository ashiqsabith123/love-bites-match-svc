package helper

import "github.com/fatih/color"


var bold = color.Bold

var Green = color.New(color.FgGreen, bold).SprintFunc()
var Blue = color.New(color.FgBlue, bold).SprintFunc()
var Red = color.New(color.FgRed, bold).SprintFunc()
var Yellow = color.New(color.FgYellow, bold).SprintFunc()