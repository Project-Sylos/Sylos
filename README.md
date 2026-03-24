# Sylos

This is the main repository in terms of user facing releases and overall links to everything. It's more of a hub repo, the main actual code of the software is mostly hosted on the other repositories though.

## Under Construction Note
This repo is mostly under construction, but if you're really interested in running Sylos, for now, see the official discord server at https://chat.sylos.io/ for further instructions on how to get started contributing to get us to early alpha testing, or see the below message about it. (But you should still join the server please!)

## Contributing
There's a lot of work that needs done and I would recommend taking a look at the pinned messages in discord in the #bounties channel for more info on what needs the most work done. You can also check out the [Active ongoing issues and roadmap](https://codeberg.org/Sylos/Discussions-and-Issues/issues/5) page for more information as to what really needs doing. This includes *both* Dev related work and Creative related work (like technical writing, music composition, and artistic assets needed).

## Getting Started (For contributors and extremely early testers). 
Head on over to the [Developer Utilities](https://codeberg.org/Sylos/Sylos-Dev-Utils) page to see our setup scripts. I have a git script in their called 'clone_repos.sh' that will clone all of the repos for you quickly and in the right places for the go replace functions to work for things like Sylos-API and the Migration-Engine repos for example. Once you do that, cd over to your Sylos-UI folder and type ./build.sh (if on linux) or cd to 'frontend' and type 'npm run dev'. You'll also need to run the API, and for that, open a new terminal or terminal tab and cd over to the Sylos-API folder you cloned and type 'go run main.go' once you've installed the dependencies.

In theory this software should also work on Windows but I've stopped developing for Windows temporarily so I am not up to date on testing it lately. 