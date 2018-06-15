// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
package main

import (
    "github.com/urfave/cli"
    "os"
    "path/filepath"
    "time"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/config"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/log"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/errors"
)

func main() {
    log.Init()

    // read help templates
    appHelpText, err := io.AssetIO.ReadFile("cli-helper/app-help.txt")
    log.WriteErrorlnExit(err)

    commandHelpText, err := io.AssetIO.ReadFile("cli-helper/command-help.txt")
    log.WriteErrorlnExit(err)

    subCommandHelpText, err := io.AssetIO.ReadFile("cli-helper/subcommand-help.txt")
    log.WriteErrorlnExit(err)

    // override help templates
    cli.AppHelpTemplate = string(appHelpText)
    cli.CommandHelpTemplate = string(commandHelpText)
    cli.SubcommandHelpTemplate = string(subCommandHelpText)

    // command that will be executed
    var command commands.Command

    app := cli.NewApp()
    app.Name = config.ProjectMeta.Name
    app.Version = config.ProjectMeta.Version
    app.EnableBashCompletion = config.ProjectMeta.EnableBashCompletion
    app.Compiled = time.Now()
    app.Copyright = config.ProjectMeta.Copyright
    app.Usage = config.ProjectMeta.UsageText

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create pkg [command options]",
                    Subcommands: cli.Commands{
                        cli.Command{
                            Name: "avr",
                            Usage: "Creates an AVR package.",
                            UsageText: "wio create pkg avr <DIRECTORY> [BOARD] [command options]",
                            Flags: []cli.Flag{
                                cli.BoolFlag{Name: "header-only",
                                    Usage: "This flag can be used to specify that the package is header only"},
                                cli.StringFlag{Name: "framework",
                                    Usage: "Framework being used for this project. Framework is Cosa/Arduino SDK",
                                    Value: config.ProjectDefaults.Framework},
                                cli.BoolFlag{Name: "create-example",
                                    Usage: "This will create an example project that user can build and upload"},
                                cli.BoolFlag{Name: "no-extras",
                                    Usage: "This will restrict wio from creating .gitignore, README.md, etc files"},
                                cli.BoolFlag{Name: "verbose",
                                    Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                                cli.BoolFlag{Name: "disable-warnings",
                                    Usage: "Disables all the warning shown by wio"},
                            },
                            Action: func(c *cli.Context) {
                                command = create.Create{Context: c, Type: create.PKG, Platform: create.AVR, Update: false}
                            },
                        },

                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app [command options]",
                    Subcommands: cli.Commands{
                        cli.Command{
                            Name:      "avr",
                            Usage:     "Creates an AVR application.",
                            UsageText: "wio create app avr <DIRECTORY> [BOARD] [command options]",
                            Flags: []cli.Flag{
                                cli.StringFlag{Name: "framework",
                                    Usage: "Framework being used for this project. Framework contains the core libraries",
                                    Value: config.ProjectDefaults.Framework},
                                cli.BoolFlag{Name: "create-example",
                                    Usage: "This will create an example project that user can build and upload"},
                                cli.BoolFlag{Name: "no-extras",
                                    Usage: "This will restrict wio from creating .gitignore, README.md, etc files"},
                                cli.BoolFlag{Name: "verbose",
                                    Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                                cli.BoolFlag{Name: "disable-warnings",
                                    Usage: "Disables all the warning shown by wio"},
                            },
                            Action: func(c *cli.Context) {
                                command = create.Create{Context: c, Type: create.APP, Platform: create.AVR, Update: false}
                            },
                        },
                    },
                },
            },
        },
        {
            Name:  "update",
            Usage: "Updates the current project and fixes any issues.",
            Flags: []cli.Flag {
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed"},
                cli.BoolFlag{Name: "disable-warnings",
                    Usage: "Disables all the warning shown by wio"},
            },
            Action: func(c *cli.Context) {
                command = create.Create{Context: c, Update: true}
            },
        },
        /*
        {
            Name:      "build",
            Usage:     "Builds the wio project.",
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Build a specified target instead of building the default",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = build.Build{Context: c}
            },
        },
        {
            Name:      "clean",
            Usage:     "Cleans all the build files for the project.",
            UsageText: "wio clean",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Cleans build files for a specified target instead of cleaning all the targets",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = clean.Clean{Context: c}
            },
        },
        {
            Name:      "run",
            Usage:     "Builds and Uploads the project to a device (provide port flag to trigger upload)",
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: config.ProjectDefaults.DefaultTarget,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: config.ProjectDefaults.Port,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                validateWioProject(c.String("dir"))
                command = run.Run{Context: c}
            },
        },
        {
            Name:      "monitor",
            Usage:     "Runs the serial monitor.",
            UsageText: "wio monitor [command options]",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "open",
                    Usage:     "Opens a Serial monitor.",
                    UsageText: "wio monitor open [command options]",
                    Flags: []cli.Flag{
                        cli.IntFlag{Name: "baud",
                            Usage: "Baud rate for the Serial port.",
                            Value: config.ProjectDefaults.Baud},
                        cli.StringFlag{Name: "port",
                            Usage: "Serial Port to open.",
                            Value: config.ProjectDefaults.Board},
                        cli.BoolFlag{Name: "gui",
                            Usage: "Runs the GUI version of the serial monitor tool",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = monitor.Monitor{Context: c, Type: monitor.OPEN}
                    },
                },
                cli.Command{
                    Name:      "ports",
                    Usage:     "Lists all the open ports and provides information about them.",
                    UsageText: "wio monitor ports [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "basic",
                            Usage: "Shows only the name of the ports."},
                        cli.BoolFlag{Name: "show-all",
                            Usage: "Shows all the ports, closed or open."},
                    },
                    Action: func(c *cli.Context) {
                        command = monitor.Monitor{Context: c, Type: monitor.PORTS}
                    },
                },
            },
        },
        */
        /*
           {
               Name:      "test",
               Usage:     "Runs unit tests available in the project.",
               UsageText: "wio test",
               Flags: []cli.Flag{
                   cli.BoolFlag{Name: "clean",
                       Usage: "Clean the project before building it",
                   },
                   cli.StringFlag{Name: "port",
                       Usage: "Port to upload the project to, (default: automatically select)",
                       Value: defaults.Port,
                   },
                   cli.StringFlag{Name: "target",
                       Usage: "Builds, and uploads a specified target instead of the main/default target",
                       Value: defaults.Utarget,
                   },
                   cli.BoolFlag{Name: "verbose",
                       Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                   },
               },
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "doctor",
               Usage:     "Guide development tooling and system configurations.",
               UsageText: "wio doctor",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "analyze",
               Usage:     "Analyzes C/C++ code statically.",
               UsageText: "wio analyze",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
           {
               Name:      "doxygen",
               Usage:     "Runs doxygen tool to create documentation for the code.",
               UsageText: "wio doxygen",
               Action: func(c *cli.Context) error {
                   return nil
               },
           },
        */
        /*
        {
            Name:  "pac",
            Usage: "Package manager for Wio projects.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "add",
                    Usage:     "Add/Update dependencies.",
                    UsageText: "wio pac add [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "vendor",
                            Usage: "Adds the dependency as vendor",
                        },
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.ADD}
                    },
                },
                cli.Command{
                    Name:      "rm",
                    Usage:     "Remove dependencies.",
                    UsageText: "wio pac rm [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "A",
                            Usage: "Delete all the dependencies",
                        },
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.RM}
                    },
                },
                cli.Command{
                    Name:      "list",
                    Usage:     "List all the dependencies",
                    UsageText: "wio pac list [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.LIST}
                    },
                },
                cli.Command{
                    Name:      "info",
                    Usage:     "Get information about a dependency being used",
                    UsageText: "wio pac info [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.INFO}
                    },
                },
                cli.Command{
                    Name:      "publish",
                    Usage:     "Publish the wio package to the package manager site (npm site)",
                    UsageText: "wio pac publish [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.PUBLISH}
                    },
                },
                cli.Command{
                    Name:      "get",
                    Usage:     "Gets all the packages mentioned in wio.yml file and vendor folder.",
                    UsageText: "wio pac get [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "clean",
                            Usage: "Cleans all the current packages and re get all of them.",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        validateWioProject(c.String("dir"))
                        command = pac.Pac{Context: c, Type: pac.GET}
                    },
                },
                /*
                   cli.Command{
                       Name:      "update",
                       Usage:     "Updates all the packages mentioned in wio.yml file and vendor folder.",
                       UsageText: "wio pac update [command options]",
                       Flags: []cli.Flag{
                           cli.StringFlag{Name: "dir",
                               Usage: "Directory for the project (default: current working directory)",
                               Value: getCurrDir(),
                           },
                           cli.BoolFlag{Name: "verbose",
                               Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                           },
                       },
                       Action: func(c *cli.Context) {
                           command = pac.Pac{Context: c, Type: pac.UPDATE}
                       },
                   },
                   cli.Command{
                       Name:      "collect",
                       Usage:     "Creates vendor folder and puts all the packages in that folder.",
                       UsageText: "wio pac collect [command options]",
                       Flags: []cli.Flag{
                           cli.StringFlag{Name: "dir",
                               Usage: "Directory for the project (default: current working directory)",
                               Value: getCurrDir(),
                           },
                           cli.StringFlag{Name: "pkg",
                               Usage: "Packages to collect instead of collecting all of the packages.",
                               Value: "none",
                           },
                           cli.BoolFlag{Name: "verbose",
                               Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                           },
                       },
                       Action: func(c *cli.Context) {
                           command = pac.Pac{Context: c, Type: pac.COLLECT}
                       },
                   },*/
            /*},
        },
        */
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    if err = app.Run(os.Args); err != nil {
        log.WriteErrorlnExit(err)
    }

    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }

        if command.GetContext().Bool("disable-warnings") {
            log.DisableWarnings()
        }

        command.Execute()
    }
}

func validateWioProject(directory string) {
    directory, err := filepath.Abs(directory)
    log.WriteErrorlnExit(err)

    if !utils.PathExists(directory) {
        err := errors.PathDoesNotExist{
            Path: directory,
        }

        log.WriteErrorlnExit(err)
    }

    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        err := errors.ConfigMissing{}

        log.WriteErrorlnExit(err)
    }
}

// returns the current directory from where wio is being called
func getCurrDir() string {
    directory, err := os.Getwd()
    log.WriteErrorlnExit(err)
    return directory
}
