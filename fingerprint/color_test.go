// Copyright 2014 Volker Dobler.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fingerprint

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strings"
	"testing"
)

var testfiles = []string{
	"boat.jpg", "clock.jpg", "lena.jpg", "lenamir.jpg", "lenablu.jpg",
	"baboon.jpg", "pepper.jpg", "logo.png",
}

var stringTests = []struct {
	ch   ColorHist
	want string
}{
	{
		ColorHist{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"Z00000000000000000000000",
	},
	{
		ColorHist{0, 5, 6, 11, 12, 18, 19, 25, 26, 32, 33, 0, 255, 248,
			247, 239, 238, 230, 229, 221, 220, 0, 102, 103},
		"022335577990ZXXVVTTRR0pq",
	},
	{
		ColorHist{0, 11, 22, 33, 44, 56, 67, 78, 89, 100, 111, 122,
			134, 145, 146, 167, 178, 189, 200, 211, 223, 234, 245, 255},
		"0369behjmpsux++EHJMPRUXZ",
	},
	{
		ColorHist{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44,
			134, 145, 146, 167, 178, 189, 200, 211, 223, 234, 245, 255},
		"012356789-abx++EHJMPRUXZ",
	},
}

func TestString(t *testing.T) {
	for i, tc := range stringTests {
		s := tc.ch.String()
		if s != tc.want {
			t.Errorf("%d: got %s, want %s", i, s, tc.want)
		}
		re, err := ColorHistFromString(s)
		if err != nil {
			t.Errorf("%d: unexpected error %q on %s", i, err, s)
		}
		r := re.String()
		if r != s {
			t.Errorf("%d: got %s, want %s", i, r, s)
		}
	}
}

func readImage(fn string) image.Image {
	file, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func TestBinColor(t *testing.T) {
	var c color.RGBA
	for _, file := range testfiles {
		img := readImage("testdata/" + file)
		bounds := img.Bounds()
		bined := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				bin := colorBin(img.At(x, y), true)
				if bin < 0 {
					c = color.RGBA{0, 0, 0, 0}
				} else {
					mb := macbeth[bin]
					c = color.RGBA{
						uint8(mb[0]),
						uint8(mb[1]),
						uint8(mb[2]),
						0xff}
				}
				bined.SetRGBA(x, y, c)
			}
		}

		out, err := os.Create("testdata/" + file + ".bined.jpg")
		if err != nil {
			t.Fatal(err.Error())
		}
		err = jpeg.Encode(out, bined, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		out.Close()
	}
}

func TestUniformColorHist(t *testing.T) {
	red := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			red.SetRGBA(x, y, color.RGBA{0xff, 0, 0, 0xff})
		}
	}
	green := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			green.SetRGBA(x, y, color.RGBA{0, 0xff, 0, 0xff})
		}
	}
	blue := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			if y < 50 {
				blue.SetRGBA(x, y, color.RGBA{0, 0, 0xff, 0xff})
			} else {
				blue.SetRGBA(x, y, color.RGBA{0, 0, 0xff, 0x10})
			}
		}
	}
	rh := NewColorHist(red)
	gh := NewColorHist(green)
	bh := NewColorHist(blue)

	if rh.String() != "00000000000000Z000000000" ||
		gh.String() != "0000000000000Z0000000000" ||
		bh.String() != "000000000000Z00000000000" {
		t.Fatalf("Got rh=%s gh=%s bh=%s", rh.String(), gh.String(), bh.String())
	}

	rg := rh.l1Norm(gh)
	// Two bins out of 24 differ completely
	if rg < 2.0/24-1e-6 || rg > 2.0/24+1e-6 {
		t.Errorf("Got %.6f, want 2/24=0.0833", rg)
	}
}

func TestL1Norm(t *testing.T) {
	ref, err := ColorHistFromString("0000000000000ZZ000000000")
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	for i, tc := range []struct {
		in   string
		want float64
	}{
		{"0000000000000ZZ000000000", 0.0 / 24},
		{"00000000000000Z000000000", 1.0 / 24},
		{"Z00000000000000000000000", 2.0 / 24},
		{"Z0000000000000000000000Z", 4.0 / 24},
		{"v000000000000Z000000000v", 2.0 / 24},
	} {

		ch, err := ColorHistFromString(tc.in)
		if err != nil {
			t.Fatalf("%d %q: Unexpected error %s", i, tc.in, err)
			continue
		}

		diff := ref.l1Norm(ch)
		ffid := ch.l1Norm(ref)
		if diff != ffid {
			t.Errorf("%d %q: l1 norm not symetric %f != %f",
				i, tc.in, diff, ffid)
		}
		if math.Abs(diff-tc.want) > 5e-4 {
			t.Errorf("%d %q: Got %.6f, want %.6f",
				i, tc.in, diff, tc.want)
		}
	}
}

func TestMaxDiffColorHist(t *testing.T) {
	first := image.NewRGBA(image.Rect(0, 0, 100, 12))
	for x := 0; x < 100; x++ {
		for y := 0; y < 12; y++ {
			c := color.RGBA{uint8(macbeth[y][0]), uint8(macbeth[y][1]),
				uint8(macbeth[y][2]), 0xff}
			first.SetRGBA(x, y, c)
		}
	}
	second := image.NewRGBA(image.Rect(0, 0, 100, 12))
	for x := 0; x < 100; x++ {
		for y := 0; y < 12; y++ {
			c := color.RGBA{uint8(macbeth[y+12][0]),
				uint8(macbeth[y+12][1]), uint8(macbeth[y+12][2]), 0xff}
			second.SetRGBA(x, y, c)
		}
	}
	h := NewColorHist(first)
	g := NewColorHist(second)
	if h.l1Norm(g) < 0.99999999 {
		t.Errorf("g=%s  h=%s  delte=%f\n", h.String(), g.String(), h.l1Norm(g))
	}
}

func basename(s string) string {
	return s[:strings.LastIndex(s, ".")]
}

func TestNewColorHist(t *testing.T) {
	hists := make(map[string]ColorHist)
	bmvs := make(map[string]BMVHash)

	for _, file := range testfiles {
		img := readImage("testdata/" + file)
		ch := NewColorHist(img)
		chstr := ch.String()
		ch2, err := ColorHistFromString(chstr)
		if err != nil {
			t.Errorf("%s: %s", file, err.Error())
		}
		ch2str := ch2.String()
		if chstr != ch2str {
			t.Errorf("%s %s != %s", file, chstr, ch2str)
		}
		fmt.Printf("%8s: %s\n", basename(file), chstr)
		hists[file] = ch
		bmvs[file] = NewBMVHash(img)
	}
	fmt.Println()

	fmt.Printf(" Δ ColorHist | ")
	for _, a := range testfiles {
		fmt.Printf("%-8s", basename(a))
	}
	fmt.Println()
	fmt.Println(" Δ BMV       |  ")
	fmt.Println("-------------+--------------------------------------------------------------")
	for _, a := range testfiles {
		h := hists[a]
		fmt.Printf(" %-11s | ", basename(a))
		for _, b := range testfiles {
			g := hists[b]
			fmt.Printf("%3.0f     ", 1000*ColorHistDelta(h, g))
		}
		fmt.Println()
		fmt.Printf("             | ")
		for _, b := range testfiles {
			fmt.Printf("%3.0f     ", 1000*BMVDelta(bmvs[a], bmvs[b]))

		}
		fmt.Println()
		fmt.Println("-------------+--------------------------------------------------------------")
	}

}

func savePngImage(t *testing.T, img *image.RGBA, name string) {
	out, err := os.Create(name)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer out.Close()
	err = png.Encode(out, img)
	if err != nil {
		t.Fatal(err.Error())
	}
	out.Close()
}

func TestColorImage(t *testing.T) {
	for _, file := range testfiles {
		img := readImage("testdata/" + file)
		ch := NewColorHist(img)
		reconstructed := ch.Image(64, 64)
		savePngImage(t, reconstructed, "testdata/"+file+".colrec.png")
		ch2, err := ColorHistFromString(ch.String())
		if err != nil {
			t.Errorf("%s", err)
		}
		reconstructed2 := ch2.Image(64, 64)
		savePngImage(t, reconstructed2, "testdata/"+file+".colrec2.png")
	}
}

func TestColorImageSpecial(t *testing.T) {
	ch, err := ColorHistFromString("3102000000f002e000021006")
	if err != nil {
		t.Fatalf("Ooops: %v", err)
	}

	ch.Image(64, 64)
}

func TestColorHistTransparent(t *testing.T) {
	transparentGreen := color.RGBA{0, 0xff, 0, 0}
	transparent := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			transparent.SetRGBA(x, y, transparentGreen)
		}
	}
	ch := NewColorHist(transparent)
	if got := ch.String(); got != "0000000000000Z0000000000" {
		t.Errorf("not green: %s", got)
	}
}
