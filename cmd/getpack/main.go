package main

import (
	"fmt"
	"getpack/config"
	"getpack/discord"
	"getpack/sources/bedrock"
	"getpack/sources/curseforge"
	"getpack/sources/fabric"
	"getpack/sources/forge"
	"getpack/sources/ftb"
	"getpack/sources/paper"
	"getpack/sources/purpur"
	"getpack/sources/spigot"
	"getpack/sources/technic"
	"log"
	"os"
	"strconv"
)

var (
	cfg config.Config
	dir string
)

func main() {
	err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %s", err)
		return
	}
	cfg = config.Global
	dir, _ = os.Getwd()

	if cfg.Discord.Enabled {
		file, err := os.OpenFile(discord.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log.txt: %s", err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	log.Printf("Starting getpack with debug: %s", strconv.FormatBool(cfg.Debug))

	err = getVersions()
	if err != nil {
		log.Println(err)
	}

	err = getPacks()
	if err != nil {
		log.Println(err)
	}

	if cfg.Discord.Enabled {
		err = discord.SendLog()
		if err != nil {
			log.Printf("Failed to log to discord: %s", err)
		}
	}
}

func getVersions() error {
	var err error

	if cfg.Enabled.Paper {
		err = paper.Get()
		if err != nil {
			log.Printf("failed to get paper: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("paperinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete paperinstaller folder: %s", err)
		}

	}

	if cfg.Enabled.Purpur {
		err = purpur.Get()
		if err != nil {
			log.Printf("failed to get purpur: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("purpurinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete purpurinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Spigot {
		err = spigot.Get()
		if err != nil {
			log.Printf("failed to get spigot: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("spigotinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete spigotinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Forge {
		err = forge.Get()
		if err != nil {
			log.Printf("failed to get forge: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("forgeinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete forgeinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Fabric {
		url, err := fabric.GetInstaller()
		if err == nil {
			err = fabric.Install(url)
			if err != nil {
				log.Printf("failed to install fabric: %s", err)
			}
		} else {
			log.Printf("failed to fabric installer: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("fabricinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete fabricinstaller folder: %s", err)
		}
	}

	if cfg.Enabled.Bedrock {
		err = bedrock.Get()
		if err != nil {
			return fmt.Errorf("failed to get bedrock: %s", err)
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("bedrockinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete bedrockinstaller folder: %s", err)
		}
	}

	return nil
}

func getPacks() error {
	if config.Global.Enabled.Curse {
		err := os.Mkdir("curseinstaller", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %s", err)
		}

		err = os.Chdir("curseinstaller")
		if err != nil {
			return fmt.Errorf("failed to go to folder: %s", err)
		}

		folder := config.Global.Target + "curse/"
		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for packname, packid := range cfg.Curse.Modpacks {
			log.Printf("Starting install of curse modpack %s", packname)
			cursepack, err := curseforge.Get(packid)
			if err != nil {
				log.Printf("Failed to get %s: %s", packname, err)
				continue
			}
			log.Printf("Found pack %s with version %s", cursepack.Data.Name, cursepack.Version)

			file := fmt.Sprintf("%scurse/%s/%s.zip", config.Global.Target, packname, cursepack.Version)
			_, err = os.Stat(file)
			if err == nil {
				log.Printf("Version %s already installed, continueing...", cursepack.Version)
				continue
			}

			err = cursepack.Install(packname)
			if err != nil {
				log.Printf("Failed to install %s: %s", packname, err)
			}

			err = os.Chdir(dir + "/curseinstaller")
			if err != nil {
				return fmt.Errorf("failed to change directories back: %s", err)
			}

			err = os.RemoveAll(packname)
			if err != nil {
				return fmt.Errorf("failed to delete %s folder: %s", cursepack.Data.Name, err)
			}
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("curseinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete curseinstaller folder: %s", err)
		}
	}

	if config.Global.Enabled.Technic {
		err := os.Mkdir("technicinstaller", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %s", err)
		}

		err = os.Chdir("technicinstaller")
		if err != nil {
			return fmt.Errorf("failed to go to folder: %s", err)
		}

		folder := config.Global.Target + "technic/"
		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for _, pack := range cfg.Technic.Modpacks {
			log.Printf("Starting install of technic modpack %s", pack)
			technicpack, err := technic.Get(pack)
			if err != nil {
				log.Printf("Failed to get %s: %s", pack, err)
				continue
			}

			log.Printf("Found pack %s with version %s", technicpack.DisplayName, technicpack.Version)

			file := fmt.Sprintf("%stechnic/%s/%s.zip", config.Global.Target, pack, technicpack.Version)
			_, err = os.Stat(file)
			if err == nil {
				log.Printf("Version %s already installed, continueing...", technicpack.Version)
				continue
			}

			err = technic.Install(technicpack)
			if err != nil {
				log.Printf("Failed to install %s: %s", pack, err)
			}

			err = os.Chdir("..")
			if err != nil {
				return fmt.Errorf("failed to change directories back: %s", err)
			}

			err = os.RemoveAll(technicpack.Name)
			if err != nil {
				return fmt.Errorf("failed to delete %s folder: %s", technicpack.Name, err)
			}
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("technicinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete technicinstaller folder: %s", err)
		}
	}

	if config.Global.Enabled.FTB {
		err := os.Mkdir("ftbinstaller", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %s", err)
		}

		err = os.Chdir("ftbinstaller")
		if err != nil {
			return fmt.Errorf("failed to go to folder: %s", err)
		}

		folder := config.Global.Target + "ftb/"
		_, err = os.Stat(folder)
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for packname, packid := range cfg.FTB.Modpacks {
			log.Printf("Starting install of ftb modpack %s", packname)
			ftbpack, err := ftb.Get(packid)
			if err != nil {
				log.Printf("Failed to get %s: %s", packname, err)
				continue
			}

			log.Printf("Found pack %s with version %s", ftbpack.Name, ftbpack.Version.Name)

			file := fmt.Sprintf("%sftb/%s/%s.zip", config.Global.Target, packname, ftbpack.Version.Name)
			_, err = os.Stat(file)
			if err == nil {
				log.Printf("Version %s already installed, continueing...", ftbpack.Version.Name)
				continue
			}

			err = ftbpack.Install(packname)
			if err != nil {
				log.Printf("Failed to install %s: %s", packname, err)
			}

			err = os.Chdir("..")
			if err != nil {
				return fmt.Errorf("failed to change directories back: %s", err)
			}

			err = os.RemoveAll(packname)
			if err != nil {
				return fmt.Errorf("failed to delete %s folder: %s", ftbpack.Name, err)
			}
		}

		err = os.Chdir(dir)
		if err != nil {
			return fmt.Errorf("failed to change directories to home: %s", err)
		}

		err = os.RemoveAll("ftbinstaller")
		if err != nil {
			return fmt.Errorf("failed to delete ftbinstaller folder: %s", err)
		}
	}

	return nil
}
