package utils

import (
	"fmt"
	"getpack/config"
	"getpack/sources/fabric"
	"getpack/sources/forge"
	"getpack/util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func Instructions(instructions []string) error {
	for _, instruction := range instructions {
		parts := strings.Split(instruction, ":")
		start := time.Now()
		switch parts[0] {
		case "delete":
			target := parts[1]
			if strings.HasSuffix(target, "/") {
				err := os.RemoveAll(target)
				if err != nil {
					return fmt.Errorf("failed to execute remove folder instruction %s: %s", instruction, err)
				}
			} else {
				err := os.Remove(target)
				if err != nil {
					return fmt.Errorf("failed to execute remove instruction %s: %s", instruction, err)
				}
			}
		case "deletereg":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					err = os.Remove(file.Name())
					if err != nil {
						return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
					}
				}
			}
		case "rename":
			err := os.Rename(parts[1], parts[2])
			if err != nil {
				return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
			}
		case "renamereg":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					err = os.Rename(file.Name(), parts[2])
					if err != nil {
						return fmt.Errorf("failed to execute rename instruction %s: %s", instruction, err)
					}
					continue
				}
			}
		case "cdreg":
			folders, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, folder := range folders {
				if folder.IsDir() && r.MatchString(folder.Name()) {
					err = os.Chdir(folder.Name())
					if err != nil {
						return fmt.Errorf("failed to execute cd instruction %s: %s", instruction, err)
					}
					continue
				}
			}
		case "javarun":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				return fmt.Errorf("cant read dir for regex match: %s", err)
			}
			r, err := regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("cant compile regex: %s", err)
			}
			for _, file := range files {
				if r.MatchString(file.Name()) {
					output, err := exec.Command("java", "-jar", file.Name(), parts[2]).Output()
					if err != nil {
						log.Println(string(output))
						return fmt.Errorf("failed to javarun for %s: %s", instruction, err)
					}
					continue
				}
			}
		case "forgegrep":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			forgeversion := ""
			mcversion := ""
			url := ""
			mirror := ""

			for _, line := range strings.Split(string(file), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "FORGE_URL=") || strings.HasPrefix(line, "FORGE_INSTALLER_URL=") {
					url = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "MIRROR=") {
					mirror = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "FORGE_VERSION=") {
					forgeversion = strings.Split(line, "=")[1]
					continue
				}
				if strings.HasPrefix(line, "MODLOADER_VERSION=") || strings.HasPrefix(line, "FORGE=") {
					forgeversion = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "MINECRAFT_VERSION=") || strings.HasPrefix(line, "MINECRAFT=") {
					mcversion = strings.Split(line, "\"")[1]
				}
			}

			if mcversion != "" {
				version := fmt.Sprintf("%s-%s", mcversion, forgeversion)
				url = fmt.Sprintf("%s%s/forge-%s-installer.jar", forge.Base, version, version)
			} else {
				url = strings.ReplaceAll(url, "$FORGE_VERSION", forgeversion)
			}

			if mirror != "" {
				url = strings.ReplaceAll(url, "${MIRROR}", mirror)
			}

			err = util.Download(url, "forge-installer.jar")
			if err != nil {
				return fmt.Errorf("failed to download forge installer: %s", err)
			}

			output, err := exec.Command("java", "-jar", "forge-installer.jar", "--installServer").Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to run installer: %s", err)
			}

			err = os.Remove("forge-installer.jar")
			if err != nil {
				return fmt.Errorf("failed to remove installer: %s", err)
			}

			_, err = os.Stat("installer.log")
			if err == nil {
				err = os.Remove("installer.log")
				if err != nil {
					return fmt.Errorf("failed to remove installer log: %s", err)
				}
			}

			_, err = os.Stat("forge-installer.jar.log")
			if err == nil {
				err = os.Remove("forge-installer.jar.log")
				if err != nil {
					return fmt.Errorf("failed to remove installer log: %s", err)
				}
			}
		case "fabricgrep":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			fabricinstallerversion := ""
			fabricurl := ""
			mcversion := ""

			for _, line := range strings.Split(string(file), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "FABRIC_INSTALLER_VERSION") {
					fabricinstallerversion = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "FABRIC_INSTALLER_URL") {
					fabricurl = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "MINECRAFT_VERSION") {
					mcversion = strings.Split(line, "\"")[1]
				}
			}

			url := strings.ReplaceAll(fabricurl, "${FABRIC_INSTALLER_VERSION}", fabricinstallerversion)

			err = util.Download(url, "installer.jar")
			if err != nil {
				return fmt.Errorf("failed to download forge installer: %s", err)
			}

			output, err := exec.Command("java", "-jar", "installer.jar", "server", "-mcversion", mcversion, "-downloadMinecraft").Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to run installer for %s: %s", mcversion, err)
			}

			err = os.Rename("server.jar", "vanilla.jar")
			if err != nil {
				return fmt.Errorf("failed to rename server jar: %s", err)
			}

			err = os.Rename("fabric-server-launch.jar", "server.jar")
			if err != nil {
				return fmt.Errorf("failed to rename server jar: %s", err)
			}

			fabricfile, err := os.Create("fabric-server-launcher.properties")
			if err != nil {
				return fmt.Errorf("failed to create fabricfile: %s", err)
			}

			_, err = fmt.Fprintln(fabricfile, "serverJar=vanilla.jar")
			if err != nil {
				return fmt.Errorf("failed to write text to fabricfile: %s", err)
			}

			err = fabricfile.Close()
			if err != nil {
				return fmt.Errorf("failed to close fabricfile: %s", err)
			}

			err = os.Remove("installer.jar")
			if err != nil {
				return fmt.Errorf("failed to remove installer: %s", err)
			}
		case "fabricrun":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			fabricversion := ""
			installerversion := ""
			mcversion := ""

			for _, line := range strings.Split(string(file), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "FABRIC_LOADER_VERSION") {
					fabricversion = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "FABRIC_VERSION") {
					installerversion = strings.Split(line, "\"")[1]
					continue
				}
				if strings.HasPrefix(line, "MC_VERSION") {
					mcversion = strings.Split(line, "\"")[1]
				}
			}

			installer := installerversion + "-installer.jar"

			output, err := exec.Command("java", "-jar", installer, "server", "-mcversion", mcversion, "-loader", fabricversion, "-downloadMinecraft").Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to run installer for %s: %s", mcversion, err)
			}

			err = fabric.PostInstall()
			if err != nil {
				return fmt.Errorf("failed to run fabric postinstall: %s", err)
			}

			err = os.Remove(installer)
			if err != nil {
				return fmt.Errorf("failed to remove installer: %s", err)
			}
		case "variablesgrep":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			modloader := ""
			mcversion := ""
			modloaderversion := ""

			for _, line := range strings.Split(string(file), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "MODLOADER=") {
					modloader = strings.Split(line, "=")[1]
					continue
				}
				if strings.HasPrefix(line, "MODLOADER_VERSION=") {
					modloaderversion = strings.Split(line, "=")[1]
					continue
				}
				if strings.HasPrefix(line, "MINECRAFT_VERSION=") {
					mcversion = strings.Split(line, "=")[1]
				}
			}

			if modloader == "Fabric" {
				err := fabricInstall(mcversion, modloaderversion)
				if err != nil {
					return fmt.Errorf("failed to install fabric: %s", err)
				}
			}

			if modloader == "Forge" {
				err := forge.Install(mcversion, modloaderversion)
				if err != nil {
					return fmt.Errorf("failed to install forge: %s", err)
				}
			}
		case "bashrun":
			output, err := exec.Command("bash", "-c", parts[1]).Output()
			if err != nil {
				log.Println(string(output))
				return fmt.Errorf("failed to bashrun for %s: %s", instruction, err)
			}
		case "shforgegrep":
			file, err := ioutil.ReadFile(parts[1])
			if err != nil {
				return fmt.Errorf("cant read file: %s", err)
			}

			if parts[1] == "run.sh" {
				err = os.Remove("run.sh")
				if err != nil {
					return fmt.Errorf("failed to remove duplicate run.sh: %s", err)
				}
			}

			lines := strings.Split(string(file), "\n")
			for _, line := range lines {
				if strings.Contains(line, "@libraries") {
					parts := strings.Split(line, " ")
					for _, part := range parts {
						if strings.HasPrefix(part, "@libraries/net") {
							startfile, err := os.Create("run.sh")
							if err != nil {
								return fmt.Errorf("cant create file: %s", err)
							}

							_, err = startfile.WriteString(fmt.Sprintf("#!/usr/bin/env sh\njava @user_jvm_args.txt %s \"$@\"", part))
							if err != nil {
								return fmt.Errorf("cant write to file: %s", err)
							}

							err = startfile.Close()
							if err != nil {
								return fmt.Errorf("cant close file: %s", err)
							}
							continue
						}
					}
				}
			}
		case "forge":
			mcversion := parts[1]
			err := forge.Install(mcversion, "")
			if err != nil {
				return fmt.Errorf("failed to install forge: %s", err)
			}
		case "fabric":
			mcversion := parts[1]
			if len(parts) >= 3 {
				err := fabricInstall(mcversion, parts[2])
				if err != nil {
					return fmt.Errorf("failed to install fabric: %s", err)
				}
			} else {
				err := fabricInstall(mcversion, "")
				if err != nil {
					return fmt.Errorf("failed to install fabric: %s", err)
				}
			}
		}
		if config.Global.Debug {
			log.Printf("Instruction '%s' took %.2fs", instruction, time.Now().Sub(start).Seconds())
		}
	}
	return nil
}

func fabricInstall(mcversion string, fabricversion string) error {
	url, err := fabric.GetInstaller()
	if err != nil {
		return fmt.Errorf("failed to get fabric installer url: %s", err)
	}

	err = util.Download(url, "installer.jar")
	if err != nil {
		return fmt.Errorf("failed to download installer: %s", err)
	}

	args := []string{"-jar", "installer.jar", "server", "-mcversion", mcversion, "-downloadMinecraft"}
	if fabricversion != "" {
		args = append(args, "-loader", fabricversion)
	}

	output, err := exec.Command("java", args...).Output()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("failed to run installer: %s", err)
	}

	err = os.Remove("installer.jar")
	if err != nil {
		return fmt.Errorf("failed to delete installer: %s", err)
	}

	err = fabric.PostInstall()
	if err != nil {
		return fmt.Errorf("failed to run fabric postinstall: %s", err)
	}

	return nil
}
