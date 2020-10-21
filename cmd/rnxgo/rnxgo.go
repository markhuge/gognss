// Command-line tool for handling RINEX files - TODO -
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/de-bkg/gognss/pkg/rinex"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Erwin Wiesensarter",
				Email: "Erwin.Wiesensarter@bkg.bund.de",
			},
		},
		Copyright: "(c) 2020 BKG Frankfurt",
		HelpName:  "rnxgo",
		Usage:     "one more RINEX toolkit",
		//UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "starttime, start",
				Usage: "consider epochs beginning at this starttime",
			},
			&cli.StringFlag{
				Name:  "endtime, end",
				Usage: "consider epochs up to this endtime",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "diff",
				//Category:    "motion",
				Usage:       "Compare two RINEX files",
				UsageText:   "diff - compare two RINEX files",
				Description: "no really, there is a lot of dooing to be done",
				ArgsUsage:   "[arrrgh]",
				// Flags: []cli.Flag{
				// 	&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				// },
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "doo!",
				BashComplete: func(c *cli.Context) {
					fmt.Fprintf(c.App.Writer, "--better\n")
				},
				Before: func(c *cli.Context) error {
					fmt.Fprintf(c.App.Writer, "brace for impact\n")
					return nil
				},
				After: func(c *cli.Context) error {
					fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
					return nil
				},
				Action: func(c *cli.Context) error {
					// c.Command.FullName()
					// c.Command.HasName("wop")
					// c.Command.Names()
					// c.Command.VisibleFlags()
					// fmt.Fprintf(c.App.Writer, "dodododododoodododddooooododododooo\n")
					// if c.Bool("forever") {
					//   c.Command.Run(c)
					// }

					if c.NArg() != 2 {
						fmt.Fprintf(c.App.Writer, "ERROR: diff needs two files to compare\n\n")
						cli.ShowCommandHelpAndExit(c, "diff", 1)
					}

					fil1 := c.Args().Get(0)
					fil2 := c.Args().Get(1)
					obs1, err := rinex.NewObsFile(fil1)
					if err != nil {
						log.Fatal(err)
					}

					obs2, err := rinex.NewObsFile(fil2)
					if err != nil {
						log.Fatal(err)
					}

					//obs1.Opts.SatSys = []rune("GR")
					err = obs1.Diff(obs2)
					return nil
				},
				OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
					fmt.Fprintf(c.App.Writer, "for shame\n")
					return err
				},
			},
			{
				Name: "compress",
				//Category:    "motion",
				Usage:     "compress RINEX files using Hatanaka and gzip",
				UsageText: "rnxgo compress [path]",
				Description: `Compress all RINEX files recursively starting by the given path. 
		Gzip format is used for compression. RINEX observation files are furthermore Hatanaka compressed.`,
				ArgsUsage: "[arrrgh]",
				// Flags: []cli.Flag{
				// 	&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				// },
				// Subcommands: []*cli.Command{
				// 	&cli.Command{
				// 		Name:   "wop",
				// 		Action: wopAction,
				// 	},
				// },
				SkipFlagParsing: false,
				HelpName:        "rnxgo compress",
				BashComplete: func(c *cli.Context) {
					fmt.Fprintf(c.App.Writer, "--better\n")
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						fmt.Fprintf(c.App.Writer, "ERROR: missing root path as argument\n\n")
						cli.ShowCommandHelpAndExit(c, "compress", 1)
					}
					rootPath := c.Args().First()
					err := compressRINEXFiles(rootPath)
					return err
				},
				OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
					fmt.Fprintf(c.App.Writer, "for shame\n")
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func compressRINEXFiles(rootdir string) error {
	err := os.Chdir(rootdir)
	if err != nil {
		return err
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			log.Printf("visited dir: %q", path)
			return nil
		}
		fmt.Printf("visited file: %q\n", path)
		rnxFil, err := rinex.NewFile(path)
		if err != nil {
			log.Printf("ERROR: %s: %v", path, err)
		}
		log.Printf("compress file %q", path)
		err = rnxFil.Compress()
		if err != nil {
			log.Printf("ERROR: %s: %v", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path %q: %v", rootdir, err)
	}

	return nil
}
