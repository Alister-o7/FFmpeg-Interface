package main

import (
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var width float32 = 640
var height float32 = 480
var appWindow fyne.Window
var ui *fyne.Container
var content *fyne.Container
var progessInfinite *widget.ProgressBarInfinite
var button_1 *widget.Button
var button_2 *widget.Button
var button_3 *widget.Button
var button_4 *widget.Button
var button_5 *widget.Button
var inFile string = "EMPTY"
var vCodec string = "EMPTY"
var aCodec string = "EMPTY"
var outFile string = "EMPTY"

const inPrefix string = "Input - "
const vCodecPrefix string = "VCodec - "
const aCodecPrefix string = "ACodec - "
const outPrefix string = "Output - "

func screenInit() {
	progessInfinite = widget.NewProgressBarInfinite()
	progessInfinite.Hide()
	ui = container.New(layout.NewHBoxLayout(), button_1, button_2, button_3, button_4, button_5)
	content = container.New(layout.NewCenterLayout())
	content.Add(ui)
	content.Add(progessInfinite)
	appWindow.SetContent(content)
}

func ffmpegConvert(i string, o string, vc string, ac string) {
	if i == "EMPTY" || o == "EMPTY" {
		dialog.ShowInformation("Status", "Input or Output is EMPTY", appWindow)
		return

	} else if vc != "EMPTY" && ac != "EMPTY" {
		cmd := exec.Command("ffmpeg", "-y", "-i", i, "-vcodec", vc, "-acodec", ac, o)
		if err := cmd.Run(); err != nil && err != os.ErrProcessDone {
			dialog.ShowInformation("Status", "Incorrect Input or Output Filename", appWindow)
			return
		}

	} else if vc != "EMPTY" && ac == "EMPTY" {
		cmd := exec.Command("ffmpeg", "-y", "-i", i, "-vcodec", vc, o)
		if err := cmd.Run(); err != nil && err != os.ErrProcessDone {
			dialog.ShowInformation("Status", "Incorrect Input or Output Filename", appWindow)
			return
		}

	} else if vc == "EMPTY" && ac != "EMPTY" {
		cmd := exec.Command("ffmpeg", "-y", "-i", i, "-acodec", ac, o)
		if err := cmd.Run(); err != nil && err != os.ErrProcessDone {
			dialog.ShowInformation("Status", "Incorrect Input or Output Filename", appWindow)
			return
		}

	} else {
		cmd := exec.Command("ffmpeg", "-y", "-i", i, o)
		if err := cmd.Run(); err != nil && err != os.ErrProcessDone {
			dialog.ShowInformation("Status", "Incorrect Input or Output Filename", appWindow)
			return
		}
	}
	dialog.ShowInformation("Status", "Done", appWindow)
}

func main() {
	os.Setenv("FYNE_THEME", "dark")
	icon, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		log.Fatal(err)
	}

	fyneApp := app.NewWithID("ffmpeginterface")
	appWindow = fyneApp.NewWindow("FFmpeg Interface")
	appWindow.Resize(fyne.NewSize(width, height))
	appWindow.SetIcon(icon)
	appWindow.SetMaster()
	appWindow.CenterOnScreen()
	appWindow.SetPadded(false)

	button_1 = widget.NewButton(inPrefix+inFile, func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				log.Fatal(err)
			} else if file == nil {
				return
			}
			filePath := file.URI().Path()

			inFile = filePath
			button_1.SetText(inPrefix + inFile)
		}, appWindow)
	})

	button_2 = widget.NewButton(vCodecPrefix+vCodec, func() {
		vCodecList := widget.NewRadioGroup([]string{"av1", "h264", "hevc", "vp9", "mpeg2video"}, func(s string) {
			vCodec = s
		})
		dialog.ShowCustomConfirm("Video Codecs", "Confirm", "Cancel", vCodecList, func(b bool) {
			if !b {
				vCodec = "EMPTY"
			}
			button_2.SetText(vCodecPrefix + vCodec)
		}, appWindow)
	})

	button_3 = widget.NewButton(aCodecPrefix+aCodec, func() {
		aCodecList := widget.NewRadioGroup([]string{"flac", "aac", "mp3", "alac", "libvorbis"}, func(s string) {
			aCodec = s
		})
		dialog.ShowCustomConfirm("Audio Codecs", "Confirm", "Cancel", aCodecList, func(b bool) {
			if !b {
				aCodec = "EMPTY"
			}
			button_3.SetText(aCodecPrefix + aCodec)
		}, appWindow)
	})

	button_4 = widget.NewButton(outPrefix+outFile, func() {
		dialog.ShowEntryDialog("Output Filename", "", func(s string) {
			outFile = s
			button_4.SetText(outPrefix + outFile)
		}, appWindow)
	})

	button_5 = widget.NewButton("RUN", func() {
		progessInfinite.Show()
		progessInfinite.Start()
		ffmpegConvert(inFile, outFile, vCodec, aCodec)
		progessInfinite.Stop()
		progessInfinite.Hide()
	})

	screenInit()

	appWindow.ShowAndRun()
}
