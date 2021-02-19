package cli

//func TestCLI_RunShell(t *testing.T) {
//	c := New()
//	_ = c.SetCommands(
//		lime.Command{
//			Keyword: "test",
//			Func: func(_ []string, _ io.Writer) error {
//				return nil
//			},
//		},
//		lime.Command{
//			Keyword: "error",
//			Func: func(_ []string, _ io.Writer) error {
//				return errors.New("failed successfully")
//			},
//		},
//	)
//	cc := c.(*cli)
//
//	outBuffer := &bytes.Buffer{}
//	in, input := os.Pipe()
//
//	c.SetOutput(outBuffer)
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	go func() {
//		c.Run()
//		wg.Done()
//	}()
//
//	if outBuffer.String() != "entering shell mode\n> " {
//		t.Error("the shell mode initialization out was unexpected:")
//	}
//
//	//test empty input
//	fmt.Fprintln(inputWriter, "")
//	if outBuffer.String() != "> " {
//		t.Error("the shell mode new line was unexpected")
//	}
//
//	fmt.Fprintln(inputWriter, "test")
//	if outBuffer.String() != "> " {
//		t.Error("the shell non-error, empty out new line was unexpected")
//	}
//
//	fmt.Fprintln(inputWriter, "error")
//	if outBuffer.String() != "failed successfully\n> " {
//		t.Error("an error from a `lime.Func` was not out in the shell")
//	}
//
//	fmt.Fprintln(inputWriter, "invalid")
//	if outBuffer.String() != fmt.Sprintf("%s\n> ", errNoMatch.Error()) {
//		t.Error("a `limeErrors.ErrNoMatch` was not out in the shell")
//	}
//
//	fmt.Fprintln(inputWriter, cc.exitWord)
//
//	wg.Wait()
//}
//
//func TestCLI_RunNamedShell(t *testing.T) {
//	c := New()
//	c.SetName("myCli")
//	cc := c.(*cli)
//
//	// Capture the input and out
//	inputReader, inputWriter, _ := os.Pipe()
//	os.Stdin = inputReader
//
//	buffer := &bytes.Buffer{}
//	c.SetOutput(buffer)
//
//	//Ensure the CLI enters and exits shell mode with no args
//	os.Args = []string{"myCli"}
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	go func() {
//		_ = c.Run()
//		wg.Done()
//	}()
//
//	if s := buffer.String(); s != "entering shell mode for myCli\n> " {
//		t.Errorf("the shell mode initialization out was unexpected: %s", s)
//	}
//
//	fmt.Fprintln(inputWriter, cc.exitWord)
//
//	wg.Wait()
//}
