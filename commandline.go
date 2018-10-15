package main

func commandlinecmd(opt options) error {
	return coptyfiletree(opt.CmdLine.Source, opt.CmdLine.Target)
}
