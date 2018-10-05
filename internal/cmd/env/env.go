package env

import (
    "os/user"
    "regexp"
    "strings"
    "wio/internal/constants"
    "wio/pkg/log"
    "wio/pkg/util"
    "wio/pkg/util/sys"

    "github.com/joho/godotenv"
    "github.com/urfave/cli"
)

var constantEnv = map[string]bool{
    "OS":      true,
    "WIOROOT": true,
}

const (
    RESET = 0
    UNSET = 1
    VIEW  = 2
)

type Env struct {
    Context *cli.Context
    Command byte
}

// get context for the command
func (env Env) GetContext() *cli.Context {
    return env.Context
}

// Runs the build command when cli build option is provided
func (env Env) Execute() error {
    currUser, err := user.Current()
    if err != nil {
        log.WriteFailure()
        return err
    }
    envFilePath := sys.Path(currUser.HomeDir, constants.WioUsr, constants.EnvFileName)

    switch env.Command {
    case RESET:
        log.Write(log.Cyan, "resetting wio environment... ")

        if err := CreateEnv(envFilePath); err != nil {
            log.WriteFailure()
            return err
        }
        log.WriteSuccess()
        break
    case UNSET:
        if env.Context.NArg() < 1 {
            return util.Error("Need minimum one variable to unset")
        }
        if err := unsetEnvironment(envFilePath, env.Context.NArg(), env.Context.Args()); err != nil {
            return err
        }
        break
    case VIEW:
        if err := showAndUpdateEnvironment(envFilePath, env.Context.NArg(), env.Context.Args()); err != nil {
            return err
        }
        break
    }

    return nil
}

func unsetEnvironment(envFilePath string, numKeys int, keys cli.Args) error {
    envData, err := godotenv.Read(envFilePath)
    if err != nil {
        return err
    }

    for i := 0; i < numKeys; i++ {
        if _, ok := envData[keys.Get(i)]; ok {
            if _, ok := constantEnv[keys.Get(i)]; ok {
                log.Errln("%s => env cannot be edited and is read only", keys.Get(i))
                continue
            }

            delete(envData, keys.Get(i))
            log.Write(log.Green, "%s", keys.Get(i))
            log.Writeln(" variable removed")

            // update wio.env file
            if err := godotenv.Write(envData, envFilePath); err != nil {
                return err
            }
        } else {
            log.Errln("%s => no such environment variable found", keys.Get(i))
        }
    }

    return nil
}

func showAndUpdateEnvironment(envFilePath string, numKeys int, keys cli.Args) error {
    envData, err := godotenv.Read(envFilePath)
    if err != nil {
        return err
    }

    if numKeys == 0 {
        for key, val := range envData {
            log.Write(log.Cyan, "%s=", key)
            log.Writeln(log.Green, "%s", val)
        }
    }

    checkIfSettingReg := regexp.MustCompile(`^[^=]+=[^=]+$`)
    keyChanged := false

    for i := 0; i < numKeys; i++ {
        givenToken := keys.Get(i)
        if checkIfSettingReg.MatchString(givenToken) {
            newKey := strings.Split(givenToken, "=")[0]
            newValue := strings.Split(givenToken, "=")[1]

            if _, ok := constantEnv[newKey]; ok {
                log.Errln("%s => env cannot be edited and is read only", newKey)
                continue
            }

            log.Write(log.Cyan, "%s=", newKey)
            log.Write(log.Green, "%s", newValue)
            log.Writeln(log.Cyan, " environment variable added/updated")

            envData[newKey] = newValue
            keyChanged = true
        } else {
            if val, ok := envData[keys.Get(i)]; ok {
                log.Write(log.Cyan, "%s=", keys.Get(i))
                log.Writeln(log.Green, "%s", val)
            } else {
                log.Errln("%s => no such environment key found", keys.Get(i))
            }
        }
    }

    if keyChanged {
        // update wio.env file
        if err := godotenv.Write(envData, envFilePath); err != nil {
            return err
        }
    }

    return nil
}

// Creates environment and overrides if there is an old environment
func CreateEnv(envFilePath string) error {
    wioRoot, err := sys.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    // create wio.env file if it does not exist
    if err := godotenv.Write(map[string]string{
        "WIOROOT": wioRoot, "OS": sys.GetOS()}, envFilePath); err != nil {
        return err
    }

    return nil
}

func LoadEnv(envFilePath string) error {
    return godotenv.Load(envFilePath)
}
