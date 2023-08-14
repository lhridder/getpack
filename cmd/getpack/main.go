package main

import (
	"encoding/json"
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
	"getpack/util"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type packversion struct {
	Name        string
	Packname    string
	Version     string
	Java        int
	Description string
	Clientlink  string
}

type packlist struct {
	Packs []string        `json:"packs"`
	Data  map[string]data `json:"data"`
}

type data struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Versions    []string `json:"versions"`
	Clientlink  string   `json:"clientlink"`
	Java        int      `json:"java"`
}

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
		err = os.Remove(discord.Logfile)
		if err != nil {
			log.Printf("Failed to delete %s: %s", discord.Logfile, err)
		}

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
	var packs []packversion

	modpacksfolder := config.Global.Target + "modpacks/"
	_, err := os.Stat(modpacksfolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(modpacksfolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create modpacks folder: %s", err)
		}
	}

	if config.Global.Enabled.Curse {
		err := os.Mkdir("curseinstaller", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %s", err)
		}

		err = os.Chdir("curseinstaller")
		if err != nil {
			return fmt.Errorf("failed to go to folder: %s", err)
		}

		cursefolder := modpacksfolder + "curse/"
		_, err = os.Stat(cursefolder)
		if os.IsNotExist(err) {
			err = os.Mkdir(cursefolder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for packname, packid := range cfg.Curse.Modpacks {
			failed := false

			log.Println("")
			log.Printf("Starting install of curse modpack %s", packname)
			cursepack, err := curseforge.Get(packid)
			if err != nil {
				log.Printf("Failed to get %s: %s", packname, err)
				continue
			}
			log.Printf("Found pack %s with version %s", cursepack.Data.Name, cursepack.Version)

			file := fmt.Sprintf("%s/%s/%s.zip", cursefolder, packname, cursepack.Version)
			_, err = os.Stat(file)
			if err == nil {
				log.Printf("Version %s already installed, continueing...", cursepack.Version)
				continue
			}

			err = cursepack.Install(packname)
			if err != nil {
				log.Printf("Failed to install %s: %s", packname, err)
				failed = true
			}

			err = os.Chdir(dir + "/curseinstaller")
			if err != nil {
				return fmt.Errorf("failed to change directories back: %s", err)
			}

			err = os.RemoveAll(packname)
			if err != nil {
				return fmt.Errorf("failed to delete %s folder: %s", cursepack.Data.Name, err)
			}

			if cfg.Deploy && !failed {
				targetfolder := fmt.Sprintf("%s%s/", cursefolder, packname)
				_, err = os.Stat(targetfolder + "/cover.png")
				if os.IsNotExist(err) {
					err = os.Chdir(targetfolder)
					if err != nil {
						return fmt.Errorf("failed to go to folder: %s", err)
					}

					err = util.Download(cursepack.Data.Logo.URL, "cover.png")
					if err != nil {
						return fmt.Errorf("failed to download image for %s: %s", packname, err)
					}

					err = os.Chdir(dir + "/curseinstaller")
					if err != nil {
						return fmt.Errorf("failed to change directories back: %s", err)
					}
				}

				javaversion := 0
				latestfile := cursepack.Data.LatestFiles[len(cursepack.Data.LatestFiles)-1]
				for _, gameversion := range latestfile.GameVersions {
					if strings.HasPrefix(gameversion, "1.") {
						javaversion = util.JavaVersion(gameversion)
					}
				}

				pack := packversion{
					Name:        cursepack.Data.Name,
					Packname:    "curse/" + packname,
					Version:     cursepack.Version,
					Java:        javaversion,
					Description: cursepack.Data.Summary,
					Clientlink:  cursepack.Data.Links.WebsiteURL,
				}
				packs = append(packs, pack)
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

		technicfolder := modpacksfolder + "technic/"
		_, err = os.Stat(technicfolder)
		if os.IsNotExist(err) {
			err = os.Mkdir(technicfolder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for _, pack := range cfg.Technic.Modpacks {
			log.Println("")
			log.Printf("Starting install of technic modpack %s", pack)
			technicpack, err := technic.Get(pack)
			if err != nil {
				log.Printf("Failed to get %s: %s", pack, err)
				continue
			}

			log.Printf("Found pack %s with version %s", technicpack.DisplayName, technicpack.Version)

			file := fmt.Sprintf("%s/%s/%s.zip", technicfolder, pack, technicpack.Version)
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

			if cfg.Deploy {
				targetfolder := fmt.Sprintf("%s%s/", technicfolder, pack)
				_, err = os.Stat(targetfolder + "/cover.png")
				if os.IsNotExist(err) {
					err = os.Chdir(targetfolder)
					if err != nil {
						return fmt.Errorf("failed to go to folder: %s", err)
					}

					err = util.Download(technicpack.Logo.URL, "cover.png")
					if err != nil {
						return fmt.Errorf("failed to download image for %s: %s", pack, err)
					}

					err = os.Chdir(dir + "/technicinstaller")
					if err != nil {
						return fmt.Errorf("failed to change directories back: %s", err)
					}
				}

				pack := packversion{
					Name:        technicpack.DisplayName,
					Packname:    "technic/" + pack,
					Version:     technicpack.Version,
					Java:        0,
					Description: technicpack.Description,
					Clientlink:  technicpack.PlatformURL,
				}
				packs = append(packs, pack)
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

		ftbfolder := modpacksfolder + "ftb/"
		_, err = os.Stat(ftbfolder)
		if os.IsNotExist(err) {
			err = os.Mkdir(ftbfolder, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create target folder: %s", err)
			}
		}

		for packname, packid := range cfg.FTB.Modpacks {
			log.Println("")
			log.Printf("Starting install of ftb modpack %s", packname)
			ftbpack, err := ftb.Get(packid)
			if err != nil {
				log.Printf("Failed to get %s: %s", packname, err)
				continue
			}

			log.Printf("Found pack %s with version %s", ftbpack.Name, ftbpack.Version.Name)

			file := fmt.Sprintf("%s/%s/%s.zip", ftbfolder, packname, ftbpack.Version.Name)
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

			if cfg.Deploy {
				targetfolder := fmt.Sprintf("%s%s/", ftbfolder, packname)
				_, err = os.Stat(targetfolder + "/cover.png")
				if os.IsNotExist(err) {
					err = os.Chdir(targetfolder)
					if err != nil {
						return fmt.Errorf("failed to go to folder: %s", err)
					}

					err = util.Download(ftbpack.Art[0].URL, "cover.png")
					if err != nil {
						return fmt.Errorf("failed to download image for %s: %s", packname, err)
					}

					err = os.Chdir(dir + "/ftbinstaller")
					if err != nil {
						return fmt.Errorf("failed to change directories back: %s", err)
					}
				}

				mcversion := ""
				version := ftbpack.Versions[len(ftbpack.Versions)-1]
				for _, target := range version.Targets {
					if target.Name == "minecraft" {
						mcversion = target.Version
					}
				}
				javaversion := util.JavaVersion(mcversion)

				pack := packversion{
					Name:        ftbpack.Name,
					Packname:    "ftb/" + packname,
					Version:     version.Name,
					Java:        javaversion,
					Description: ftbpack.Description,
					Clientlink:  "https://www.feed-the-beast.com/modpacks/" + strconv.Itoa(packid),
				}
				packs = append(packs, pack)
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

	if cfg.Deploy {
		if len(packs) == 0 {
			return nil
		}

		err := os.Chdir(modpacksfolder)
		if err != nil {
			return fmt.Errorf("failed to go to folder: %s", err)
		}

		var list packlist
		_, err = os.Stat("list.json")
		if err != nil {
			log.Println("Can't find old list.json, creating new one")
			_, err := os.Create("list.json")
			if err != nil {
				return fmt.Errorf("failed to create list.json: %s", err)
			}
		} else {
			oldlist, err := ioutil.ReadFile("list.json")
			if err != nil {
				return fmt.Errorf("failed to read list.json: %s", err)
			}
			err = json.Unmarshal(oldlist, &list)
			if err != nil {
				return fmt.Errorf("failed to unmarshal list.json: %s", err)
			}
		}

		if list.Data == nil {
			list.Data = make(map[string]data)
		}

		for _, pack := range packs {
			packdata, ok := list.Data[pack.Packname]
			if ok {
				present := false
				for _, version := range packdata.Versions {
					if version == pack.Version {
						present = true
					}
				}
				if !present {
					packdata.Versions = append(packdata.Versions, pack.Version)
					list.Data[pack.Packname] = packdata
				}
			} else {
				list.Data[pack.Packname] = data{
					Name:        pack.Name,
					Description: pack.Description,
					Versions:    []string{pack.Version},
					Clientlink:  pack.Clientlink,
					Java:        pack.Java,
				}
			}

			present := false
			for _, modpack := range list.Packs {
				if modpack == pack.Packname {
					present = true
					break
				}
			}

			if !present {
				list.Packs = append(list.Packs, pack.Packname)
			}
		}

		newlist, err := json.Marshal(list)
		if err != nil {
			return fmt.Errorf("failed to marshal list.json: %s", err)
		}

		err = ioutil.WriteFile("list.json", newlist, 0644)
		if err != nil {
			return fmt.Errorf("failed to write list.json: %s", err)
		}

		log.Println("Deployed modpacks to list.json")
	}

	return nil
}
