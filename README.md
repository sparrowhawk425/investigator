# International Investigators

An investigation game where the player attempts to track and catch criminals in different locations around the world. The player can choose locations to travel to and then pursue a number of members of a criminal Syndicate through a country or region. Using clues left behind at crime scenes, they try to identify Syndicate members and arrest them before they leave the area.

## Motivation
This project uses Go to build a CLI game inspired by the classic Carmen Sandiego game. It uses Go HTTP calls to generate names and places from an API tool, https://randomuser.me/api.

## Quick Start
Use the Go toolchain to install the application.
```bash
# Install International Investigators
go install github.com/sparrowhawk425/investigator@latest

# Run the game
investigators
```

## Playing the game
The game generates fake people and place names from different countries around the world. After entering your name, you can select a country/region to investigate and it will generate people and places to investigate. The game uses a CLI to select actions. You can use `help` to see the available actions.

The goal is to find and arrest all of the Syndicate members before they achieve their goal and escape. To aid in their capture, new crimes are reported each day and you can investigate locations in hopes criminals left clues behind. Once you find something worth recording, you can create dossiers on the criminals you are hunting, helping you identify which character you are looking for.

Currently, the game ends when all the Syndicate members in the area are gone, either arrested or escaped. Each criminal has a goal and if they achieve it, they will leave the area.

## Usage
The game is based on a CLI design, offering the player different actions for aiding in their hunt for criminals.

```
visit: Travel to one of the locations to look for clues
crimes: View a list of reported crimes
dossiers: View and create dossiers on the Syndicate members
arrest: Try to arrest a Syndicate member
help: Display the list of available commands
```

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/sparrowhawk425/investigator@latest
cd zipzod
```

### Build the compiled binary

```bash
go build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## Future Plans

The eventual goal is to allow the player to pursue Syndicate members up the chain toward the Big Boss (Carmen Sandiego-style). When you arrest a member, it would give a clue for a new country/region to investigate and lead you to new Syndicate members, until you find the head honcho.

I want to add more variety in terms of criminals to chase and personality traits to affect their behavior.

In terms of QoL, I want to build robust filtering options since the lists of locations and people will quickly get unwieldy to look at, but I'm still trying to figure out how to build the generic filter menus.