# International Investigators

An investigation game where the player attempts to track and catch criminals in different locations around the world.

## Install
After downloading the source, use `go build` and `go install` to generate the executable.

## Playing the game
The game uses https://randomuser.me/api to generate fake people and place names from different countries around the world. After entering your name, you can select a country/region to investigate and it will generate people and places to investigate. The game uses a CLI to select actions. You can use `help` to see the available actions.

The goal is to find and arrest all of the Syndicate members before they achieve their goal and escape. To aid in their capture, new crimes are reported each day and you can investigate locations in hopes criminals left clues behind. Once you find something worth recording, you can create dossiers on the criminals you are hunting, helping you identify which character you are looking for.

Currently, the game ends when all the Syndicate members in the area are gone, either arrested or escaped. Each criminal has a goal and if they achieve it, they will leave the area.

## Future Plans

The eventual goal is to allow the player to pursue Syndicate members up the chain toward the Big Boss (Carmen Sandiego-style). When you arrest a member, it would give a clue for a new country/region to investigate and lead you to new Syndicate members, until you find the head honcho.

I want to add more variety in terms of criminals to chase and personality traits to affect their behavior.

In terms of QoL, I want to build robust filtering options since the lists of locations and people will quickly get unwieldy to look at, but I'm still trying to figure out how to build the generic filter menus.