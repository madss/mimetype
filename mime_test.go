package mimetype

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabriel-vasile/mimetype/internal/matchers"
)

const testDataDir = "testdata"

var files = map[string]*node{
	// archives
	"a.pdf":  pdf,
	"a.zip":  zip,
	"a.tar":  tar,
	"a.xls":  xls,
	"a.xlsx": xlsx,
	"a.doc":  doc,
	"a.docx": docx,
	"a.ppt":  ppt,
	"a.pptx": pptx,
	"a.epub": epub,
	"a.7z":   sevenZ,
	"a.jar":  jar,
	"a.gz":   gzip,
	"a.fits": fits,
	"a.xar":  xar,
	"a.bz2":  bz2,

	// images
	"a.png":  png,
	"a.jpg":  jpg,
	"a.psd":  psd,
	"a.webp": webp,
	"a.tif":  tiff,
	"a.ico":  ico,
	"a.bmp":  bmp,

	// video
	"a.mp4":  mp4,
	"b.mp4":  mp4,
	"a.webm": webM,
	"a.3gp":  threeGP,
	"a.3g2":  threeG2,
	"a.flv":  flv,
	"a.avi":  avi,
	"a.mov":  quickTime,
	"a.mqv":  mqv,
	"a.mpeg": mpeg,
	"a.mkv":  mkv,
	"a.asf":  asf,

	// audio
	"a.mp3":  mp3,
	"a.wav":  wav,
	"a.flac": flac,
	"a.midi": midi,
	"a.ape":  ape,
	"a.aiff": aiff,
	"a.au":   au,
	"a.ogg":  ogg,
	"a.amr":  amr,
	"a.mpc":  musePack,
	"a.m4a":  m4a,
	"a.m4b":  aMp4,

	// source code
	"a.html":    html,
	"a.svg":     svg,
	"b.svg":     svg,
	"a.txt":     txt,
	"a.php":     php,
	"a.ps":      ps,
	"a.json":    json,
	"a.geojson": geoJson,
	"b.geojson": geoJson,
	"a.csv":     csv,
	"a.tsv":     tsv,
	"a.rtf":     rtf,
	"a.js":      js,
	"a.lua":     lua,
	"a.pl":      perl,
	"a.py":      python,
	"a.tcl":     tcl,

	// binary
	"a.class": class,
	"a.swf":   swf,
	"a.crx":   crx,
	"a.wasm":  wasm,
	"a.exe":   exe,

	// fonts
	"a.woff":  woff,
	"a.woff2": woff2,

	// XML and subtypes of XML
	"a.xml": xml,
	"a.kml": kml,
	"a.dae": collada,
	"a.gml": gml,
	"a.gpx": gpx,
	"a.tcx": tcx,
	"a.x3d": x3d,

	"a.shp": shp,
	"a.shx": shx,
	"a.dbf": dbf,
}

func TestMatching(t *testing.T) {
	errStr := "File: %s; Mime: %s != DetectedMime: %s; err: %v"
	for fName, node := range files {
		fileName := filepath.Join(testDataDir, fName)
		f, err := os.Open(fileName)
		if err != nil {
			t.Fatal(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if dMime, _ := Detect(data); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, nil)
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			t.Errorf(errStr, fName, node.mime, root.mime, err)
		}

		if dMime, _, err := DetectReader(f); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, err)
		}
		f.Close()

		if dMime, _, err := DetectFile(fileName); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, err)
		}
	}
}

func TestFaultyInput(t *testing.T) {
	inexistent := "inexistent.file"
	if _, _, err := DetectFile(inexistent); err == nil {
		t.Errorf("%s should not match successfully", inexistent)
	}

	f, _ := os.Open(inexistent)
	if _, _, err := DetectReader(f); err == nil {
		t.Errorf("%s reader should not match successfully", inexistent)
	}
}

func TestEmptyInput(t *testing.T) {
	if m, _ := Detect([]byte{}); m != "inode/x-empty" {
		t.Errorf("failed to detect empty file")
	}
}

func TestGenerateSupportedMimesFile(t *testing.T) {
	f, err := os.OpenFile("supported_mimes.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	nodes := root.flatten()
	header := fmt.Sprintf(`## %d Supported MIME types
This file is automatically generated when running tests. Do not edit manually.

Extension | MIME type
--------- | --------
`, len(nodes))

	if _, err := f.WriteString(header); err != nil {
		t.Fatal(err)
	}
	for _, n := range nodes {
		ext := n.extension
		if ext == "" {
			ext = "n/a"
		}
		str := fmt.Sprintf("**%s** | %s\n", ext, n.mime)
		if _, err := f.WriteString(str); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkMatchDetect(b *testing.B) {
	files := []string{"a.png", "a.jpg", "a.pdf", "a.zip", "a.docx", "a.doc"}
	data, fLen := [][matchers.ReadLimit]byte{}, len(files)
	for _, f := range files {
		d := [matchers.ReadLimit]byte{}

		file, err := os.Open(filepath.Join(testDataDir, f))
		if err != nil {
			b.Fatal(err)
		}

		io.ReadFull(file, d[:])
		data = append(data, d)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Detect(data[n%fLen][:])
	}
}
