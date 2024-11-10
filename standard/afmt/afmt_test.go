/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package afmt_test

import (
	"bytes"
	"fmt"
	"github.com/auvitly/go-tools/afmt"
	"io/ioutil"
	"os"
	"testing"
)

// afmtFunc is used to identify which public function of the afmt package or
// ConfigState a test applies to.
type afmtFunc int

const (
	fCSFdump afmtFunc = iota
	fCSFprint
	fCSFprintf
	fCSFprintln
	fCSPrint
	fCSPrintln
	fCSSdump
	fCSSprint
	fCSSprintf
	fCSSprintln
	fCSErrorf
	fCSNewFormatter
	fErrorf
	fFprint
	fFprintln
	fPrint
	fPrintln
	fSdump
	fSprint
	fSprintf
	fSprintln
)

// Map of afmtFunc values to names for pretty printing.
var afmtFuncStrings = map[afmtFunc]string{
	fCSFdump:        "ConfigState.Fdump",
	fCSFprint:       "ConfigState.Fprint",
	fCSFprintf:      "ConfigState.Fprintf",
	fCSFprintln:     "ConfigState.Fprintln",
	fCSSdump:        "ConfigState.Sdump",
	fCSPrint:        "ConfigState.Print",
	fCSPrintln:      "ConfigState.Println",
	fCSSprint:       "ConfigState.Sprint",
	fCSSprintf:      "ConfigState.Sprintf",
	fCSSprintln:     "ConfigState.Sprintln",
	fCSErrorf:       "ConfigState.Errorf",
	fCSNewFormatter: "ConfigState.NewFormatter",
	fErrorf:         "afmt.Errorf",
	fFprint:         "afmt.Fprint",
	fFprintln:       "afmt.Fprintln",
	fPrint:          "afmt.Print",
	fPrintln:        "afmt.Println",
	fSdump:          "afmt.Sdump",
	fSprint:         "afmt.Sprint",
	fSprintf:        "afmt.Sprintf",
	fSprintln:       "afmt.Sprintln",
}

func (f afmtFunc) String() string {
	if s, ok := afmtFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown afmtFunc (%d)", int(f))
}

// afmtTest is used to describe a test to be performed against the public
// functions of the afmt package or ConfigState.
type afmtTest struct {
	cs     *afmt.ConfigState
	f      afmtFunc
	format string
	in     interface{}
	want   string
}

// afmtTests houses the tests to be performed against the public functions of
// the afmt package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var afmtTests []afmtTest

// redirStdout is a helper function to return the standard output from f as a
// byte slice.
func redirStdout(f func()) ([]byte, error) {
	tempFile, err := ioutil.TempFile("", "ss-test")
	if err != nil {
		return nil, err
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName) // Ignore error

	origStdout := os.Stdout
	os.Stdout = tempFile
	f()
	os.Stdout = origStdout
	tempFile.Close()

	return ioutil.ReadFile(fileName)
}

func initSpewTests() {
	// Config states with various settings.
	scsDefault := afmt.NewDefaultConfig()
	scsNoMethods := &afmt.ConfigState{Indent: " ", DisableMethods: true}
	scsNoPmethods := &afmt.ConfigState{Indent: " ", DisablePointerMethods: true}
	scsMaxDepth := &afmt.ConfigState{Indent: " ", MaxDepth: 1}
	scsContinue := &afmt.ConfigState{Indent: " ", ContinueOnMethod: true}
	scsNoPtrAddr := &afmt.ConfigState{DisablePointerAddresses: true}
	scsNoCap := &afmt.ConfigState{DisableCapacities: true}

	// Variables for tests on types which implement Stringer interface with and
	// without a pointer receiver.
	ts := stringer("test")
	tps := pstringer("test")

	type ptrTester struct {
		s *struct{}
	}
	tptr := &ptrTester{s: &struct{}{}}

	// depthTester is used to test max depth handling for structs, array, slices
	// and maps.
	type depthTester struct {
		ic    indirCir1
		arr   [1]string
		slice []string
		m     map[string]int
	}
	dt := depthTester{indirCir1{nil}, [1]string{"arr"}, []string{"slice"},
		map[string]int{"one": 1}}

	// Variable for tests on types which implement error interface.
	te := customError(10)

	afmtTests = []afmtTest{
		{scsDefault, fCSFdump, "", int8(127), "(int8) 127\n"},
		{scsDefault, fCSFprint, "", int16(32767), "32767"},
		{scsDefault, fCSFprintf, "%v", int32(2147483647), "2147483647"},
		{scsDefault, fCSFprintln, "", int(2147483647), "2147483647\n"},
		{scsDefault, fCSPrint, "", int64(9223372036854775807), "9223372036854775807"},
		{scsDefault, fCSPrintln, "", uint8(255), "255\n"},
		{scsDefault, fCSSdump, "", uint8(64), "(uint8) 64\n"},
		{scsDefault, fCSSprint, "", complex(1, 2), "(1+2i)"},
		{scsDefault, fCSSprintf, "%v", complex(float32(3), 4), "(3+4i)"},
		{scsDefault, fCSSprintln, "", complex(float64(5), 6), "(5+6i)\n"},
		{scsDefault, fCSErrorf, "%#v", uint16(65535), "(uint16)65535"},
		{scsDefault, fCSNewFormatter, "%v", uint32(4294967295), "4294967295"},
		{scsDefault, fErrorf, "%v", uint64(18446744073709551615), "18446744073709551615"},
		{scsDefault, fFprint, "", float32(3.14), "3.14"},
		{scsDefault, fFprintln, "", float64(6.28), "6.28\n"},
		{scsDefault, fPrint, "", true, "true"},
		{scsDefault, fPrintln, "", false, "false\n"},
		{scsDefault, fSdump, "", complex(-10, -20), "(complex128) (-10-20i)\n"},
		{scsDefault, fSprint, "", complex(-1, -2), "(-1-2i)"},
		{scsDefault, fSprintf, "%v", complex(float32(-3), -4), "(-3-4i)"},
		{scsDefault, fSprintln, "", complex(float64(-5), -6), "(-5-6i)\n"},
		{scsNoMethods, fCSFprint, "", ts, "test"},
		{scsNoMethods, fCSFprint, "", &ts, "<*>test"},
		{scsNoMethods, fCSFprint, "", tps, "test"},
		{scsNoMethods, fCSFprint, "", &tps, "<*>test"},
		{scsNoPmethods, fCSFprint, "", ts, "stringer test"},
		{scsNoPmethods, fCSFprint, "", &ts, "<*>stringer test"},
		{scsNoPmethods, fCSFprint, "", tps, "test"},
		{scsNoPmethods, fCSFprint, "", &tps, "<*>stringer test"},
		{scsMaxDepth, fCSFprint, "", dt, "{{<max>} [<max>] [<max>] map[<max>]}"},
		{scsMaxDepth, fCSFdump, "", dt, "(afmt_test.depthTester) {\n" +
			" ic: (afmt_test.indirCir1) {\n  <max depth reached>\n },\n" +
			" arr: ([1]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" slice: ([]string) (len=1 cap=1) {\n  <max depth reached>\n },\n" +
			" m: (map[string]int) (len=1) {\n  <max depth reached>\n }\n}\n"},
		{scsContinue, fCSFprint, "", ts, "(stringer test) test"},
		{scsContinue, fCSFdump, "", ts, "(afmt_test.stringer) " +
			"(len=4) (stringer test) \"test\"\n"},
		{scsContinue, fCSFprint, "", te, "(error: 10) 10"},
		{scsContinue, fCSFdump, "", te, "(afmt_test.customError) " +
			"(error: 10) 10\n"},
		{scsNoPtrAddr, fCSFprint, "", tptr, "<*>{<*>{}}"},
		{scsNoPtrAddr, fCSSdump, "", tptr, "(*afmt_test.ptrTester)({\ns: (*struct {})({\n})\n})\n"},
		{scsNoCap, fCSSdump, "", make([]string, 0, 10), "([]string) {\n}\n"},
		{scsNoCap, fCSSdump, "", make([]string, 1, 10), "([]string) (len=1) {\n(string) \"\"\n}\n"},
	}
}

// TestSpew executes all of the tests described by afmtTests.
func TestSpew(t *testing.T) {
	initSpewTests()

	t.Logf("Running %d tests", len(afmtTests))
	for i, test := range afmtTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fCSFdump:
			test.cs.Fdump(buf, test.in)

		case fCSFprint:
			test.cs.Fprint(buf, test.in)

		case fCSFprintf:
			test.cs.Fprintf(buf, test.format, test.in)

		case fCSFprintln:
			test.cs.Fprintln(buf, test.in)

		case fCSPrint:
			b, err := redirStdout(func() { test.cs.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSPrintln:
			b, err := redirStdout(func() { test.cs.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSSdump:
			str := test.cs.Sdump(test.in)
			buf.WriteString(str)

		case fCSSprint:
			str := test.cs.Sprint(test.in)
			buf.WriteString(str)

		case fCSSprintf:
			str := test.cs.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fCSSprintln:
			str := test.cs.Sprintln(test.in)
			buf.WriteString(str)

		case fCSErrorf:
			err := test.cs.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fCSNewFormatter:
			fmt.Fprintf(buf, test.format, test.cs.NewFormatter(test.in))

		case fErrorf:
			err := afmt.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fFprint:
			afmt.Fprint(buf, test.in)

		case fFprintln:
			afmt.Fprintln(buf, test.in)

		case fPrint:
			b, err := redirStdout(func() { afmt.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fPrintln:
			b, err := redirStdout(func() { afmt.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fSdump:
			str := afmt.Sdump(test.in)
			buf.WriteString(str)

		case fSprint:
			str := afmt.Sprint(test.in)
			buf.WriteString(str)

		case fSprintf:
			str := afmt.Sprintf(test.format, test.in)
			buf.WriteString(str)

		case fSprintln:
			str := afmt.Sprintln(test.in)
			buf.WriteString(str)

		default:
			t.Errorf("%v #%d unrecognized function", test.f, i)
			continue
		}
		s := buf.String()
		if test.want != s {
			t.Errorf("ConfigState #%d\n got: %s want: %s", i, s, test.want)
			continue
		}
	}
}
