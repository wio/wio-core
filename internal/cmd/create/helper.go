package create

import (
	"os"
	"path/filepath"
	"strings"
	"wio/internal/constants"
	"wio/pkg/log"
	"wio/pkg/util/sys"
	"wio/pkg/util/template"
)

// when all tag is specified, platform will become "" so it can be omitted
func getPlatform(platformProvided string) string {
	if platformProvided == "all" {
		return ""
	} else {
		return platformProvided
	}
}

// when all tag is specified, framework will become "" so it can be omitted
func getFramework(frameworkProvided string) string {
	if frameworkProvided == "all" {
		return ""
	} else {
		return frameworkProvided
	}
}

// when all tag is specified, board will become "" so it can be omitted
func getBoard(boardProvided string) string {
	if boardProvided == "all" {
		return ""
	} else {
		return boardProvided
	}
}

func (info createInfo) fillReadMe(queue *log.Queue, readmeFile string) error {
	log.Verb(queue, "filling README file ... ")

	platform := info.platform
	if info.projectType == constants.App {
		if platform == "all" {
			platform = getPlatform(platform)
		} else {
			platform += " "
		}
	}

	if err := template.IOReplace(readmeFile, map[string]string{
		"PLATFORM":     platform,
		"FRAMEWORK":    info.framework,
		"PROJECT_NAME": info.name,
	}); err != nil {
		log.WriteFailure(queue, log.VERB)
		return err
	}
	log.WriteSuccess(queue, log.VERB)
	return nil
}

func (info createInfo) toLowerCase() {
	info.projectType = strings.ToLower(info.projectType)
	info.platform = strings.ToLower(info.platform)
	info.framework = strings.ToLower(info.framework)
	info.board = strings.ToLower(info.board)
	info.ide = strings.ToLower(info.ide)
}

func (info *createInfo) generateConstraints() (map[string]bool, map[string]bool) {
	context := info.context
	dirConstraints := map[string]bool{
		"tests":        false,
		"!header-only": !context.Bool("header-only"),
		"header-only":  context.Bool("header-only"),
	}
	fileConstraints := map[string]bool{
		"extra":        !context.Bool("no-extras"),
		"example":      context.Bool("create-example"),
		"!header-only": !context.Bool("header-only"),
		"header-only":  context.Bool("header-only"),
	}

	if info.ide != "none" {
		fileConstraints["ide="+info.ide] = true
		dirConstraints["ide="+info.ide] = true
	}

	return dirConstraints, fileConstraints
}

// This uses a structure.json file and creates a project structure based on that. It takes in consideration
// all the constrains and copies files. This should be used for creating project for any type of app/pkg
func copyProjectAssets(queue *log.Queue, info *createInfo, data *StructureTypeData, updateOverride bool) error {
	dirConstraints, fileConstraints := info.generateConstraints()
	for _, path := range data.Paths {
		directoryPath := sys.Path(info.directory, path.Entry)
		skipDir := false
		log.Verbln(queue, "copying assets to directory: %s", directoryPath)
		// handle directory constraints
		for _, constraint := range path.Constraints {
			_, exists := dirConstraints[constraint]
			if exists && !dirConstraints[constraint] {
				log.Verbln(queue, "constraint not satisfied and hence skipping this directory")
				skipDir = true
				break
			} else if !exists {
				log.Verbln(queue, "constraint not specified and hence skipping this directory")
				skipDir = true
				break
			}
		}
		if skipDir {
			continue
		}

		if !sys.Exists(directoryPath) {
			if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
				return err
			}
			log.Verbln(queue, "created directory: %s", directoryPath)
		}

		log.Verbln(queue, "copying asset files for directory: %s", directoryPath)
		for _, file := range path.Files {
			toPath := sys.Path(directoryPath, file.To)
			skipFile := false
			// handle file constraints
			for _, constraint := range file.Constraints {
				_, exists := fileConstraints[constraint]
				if exists && !fileConstraints[constraint] {
					log.Verbln(queue, "constraint not satisfied and hence skipping this file")
					skipFile = true
					break
				} else if !exists {
					log.Verbln(queue, "constraint not specified and hence skipping this file")
					skipFile = true
					break
				}
			}
			if skipFile {
				continue
			}

			if info.updateOnly && !file.Update {
				log.Verbln(queue, "file doesn't support updates, hence skipping update for path: %s", toPath)
				continue
			}

			if info.updateOnly && file.Update && file.AllowFull && updateOverride {
				file.Override = true
			}

			// copy assets
			if err := sys.AssetIO.CopyFile(file.From, toPath, file.Override); err != nil {
				return err
			} else {
				log.Verbln(queue, `copied asset file "%s" TO: %s: `, filepath.Base(file.From), toPath)
			}
		}
	}
	return nil
}
