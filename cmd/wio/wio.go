// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
// Author: Deep Dhillon
// Last Modified: 03/18/2018

/*
    Wio is a tool to make development of embedded system applications easier and simpler.
    It allows for building, testing, and uploading AVR applications for Commandline.
 */


package main

import (
    "time"
    "os"
    "log"
    "fmt"
    "../util"
    "github.com/urfave/cli"
)

func main()  {
    // override help template
    cli.AppHelpTemplate =
`Wio a simplified development process for embedded applications.
Create, Build, Test, and Upload AVR projects from Commandline.

Common Commands:
    
    wio create <app type> [options] <output directory>
        Create a new Wio project in the specified directory.
    
    wio build [options]
        Build the Wio project based on all the configurations defined
    
    wio upload [options]
        Upload the Wio project to an attached embedded device
    
    wio run [options]
        Builds, Tests, and Uploads the Wio projects

Usage: {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
Global options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
Available commands:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
Copyright:
   {{.Copyright}}
   {{end}}{{if .Version}}
Vesrion:
   {{.Version}}
   {{end}}
Run "wio command <help>" for more information about a command.
`

cli.CommandHelpTemplate =
`{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

Category:
   {{.Category}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .VisibleFlags}}

Available commands:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`

cli.SubcommandHelpTemplate =
`{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

Available commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`
    defaults := Config{}
    data, _ := util.FileToString("cmd/wio/wio-config.yml")
    util.ToYmlStruct(data, &defaults)
    
    app := cli.NewApp()
    app.Name = "wio"
    app.Version = defaults.Version
    app.EnableBashCompletion = true
    app.Compiled = time.Now()
    app.Authors = []cli.Author{
        {
            Name:  "Deep Dhillon",
            Email: "deep@deepdhillon.ca",
        },
        {
            Name:  "Jeff Niu",
            Email: "jeffniu22@gmail.com",
        },
    }
    app.Copyright = "Copyright (c) 2018 Waterloop"
    app.Usage = "Create, Build and Upload AVR projects"

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.\n\nIt also works as an updater when called on already created projects.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "package",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create package <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},

                    },
                    Action: func(c *cli.Context) error {
                        if len(c.Args()) == 0 {
                            fmt.Println("A Directory path/name is needed to create this wio application!")
                            fmt.Println("Execute `wio create app -h` for more details and help")
                            os.Exit(1)
                        }
                        fmt.Println("Hello package")
                        fmt.Println(c.String("ide"))
                        return nil
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app <BOARD> <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                    },
                    Action: func(c *cli.Context) error {
                        // check if user defined a board
                        if len(c.Args()) == 0 {
                            fmt.Println("A Board is needed to create this wio application!")
                            fmt.Println("Execute `wio create app -h` for more details and help")
                            os.Exit(1)
                        } else if len(c.Args()) == 1 {
                            fmt.Println("A Directory path/name is needed to create this wio application!")
                            fmt.Println("Execute `wio create app -h` for more details and help")
                            os.Exit(1)
                        }
                        fmt.Println("Hello app")
                        fmt.Println(c.String("platform"))
                        return nil
                    },
                },
            },
        },
        {
            Name:      "build",
            Usage:     "Builds the project",
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",

                },
            },
            Action: func(c *cli.Context) error {
                /*
                build based on the type of project (from config file)
                 */
                return nil
            },
        },
        {
            Name:      "clean",
            Usage:     "Cleans all the build files for the project",
            UsageText: "wio clean",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "upload",
            Usage:     "Uploads the project to a device",
            UsageText: "wio upload [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to",
                    Value: defaults.Port,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "run",
            Usage:     "Builds, Tests, and Uploads the project to a device",
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "test",
            Usage:     "Runs unit tests available in the project",
            UsageText: "wio test",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "monitor",
            Usage:     "Runs the serial monitor",
            UsageText: "wio monitor [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "gui",
                    Usage: "Runs the GUI version of the serial monitor tool",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doctor",
            Usage:     "Show information about the installed tooling",
            UsageText: "wio doctor",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "configure",
            Usage:     "Configures paths for the tools used for development",
            UsageText: "wio configure [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "arduino-sdk-dir",
                    Usage: "path to Arduino SDK",
                },
                cli.StringFlag{Name: "make-path",
                    Usage: "path to `make` tool",
                },
                cli.StringFlag{Name: "cmake-path",
                    Usage: "Path to `cmake` tool",
                },
                cli.StringFlag{Name: "avr-path",
                    Usage: "Path to AVR libraries",
                },
                cli.StringFlag{Name: "arm-path",
                    Usage: "Path to ARM libraries",
                },
            },
            Action: func(c *cli.Context) error {
                /*
                If no flag provided, show current settings
                 */
                return nil
            },
        },
        {
            Name:      "analyze",
            Usage:     "Analyzes C/C++ code statically",
            UsageText: "wio analyze",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doxygen",
            Usage:     "Runs doxygen tool to create documentation for the code",
            UsageText: "wio doxygen",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "upgrade",
            Usage:     "Upgrades the current version of the program",
            UsageText: "wio upgrade [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "version",
                    Usage: "Specify the exact version to upgrade/downgrade wio to",
                    Value: defaults.Version,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    err := app.Run(os.Args)

    if err != nil {
        log.Fatal(err)
    }
}