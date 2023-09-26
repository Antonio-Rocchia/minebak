# Minebak

A tool to backup your Minecraft world from your MultiCraft server directly from the command line.

Minebak tries to be [human centric and discoverable](https://clig.dev/) such that the program can be used by user with different level of command line experience while also be powerful enough to be used in scripts to automate the backups


## Table of content
1. [Why this tool](#why-this-tool)
2. [Features](#features)
3. [How to install](#install)
4. [Examples](#examples)
5. [Licence](#licence)

## Why this tool
Me and my friends used to play a lot of minecraft worlds and everytime we got bored we always forgot to backup the world to play in the future.

That's because using Multicraft's suggestend method, FileZilla, is really tedious in my opinion. This tool helps us to enjoy our minecraft world for as long as we want.

A big win for me was seeing my friends who have never used a terminal use this tool almost weekly to get a copy of our world. I know that you can achieve similar result using wget or even curl but minebak can guide you interactively, for example it could be expanded to help you find your world even if you don't directly know its name.

## Features

1. Interactive input
2. Scriptable with flags, interactive by default
3. Automatic renaming of the backup folder
4. Timestamps can be added to the backup folder name (using flags) 

[Check the Examples](#Examples)

## Install
On the right of this page click on "Releases" and download the appropriate version for your operating system.

To use minebak, place the executatable file in a folder and run it. [Check the examples](#examples)

## Examples
The program was tested on Linux Ubuntu 22.04, WSL 2 Ubuntu 22.04, Windows 11
### Windows (powershell)
```shell
# To use minebak use it from the terminal
# WorldName is a required argument, it refers to your minecraft's world name on the server
> .\minebak WorldName

# Interactive by default
# minebak will ask you the MultiCraft ip, port, username and password
> .\minebak WorldName
# Prompt for missing information
> .\minebak WorldName --addr 127.0.0.1 # missing port, username and password

# But interactive input can be disabled
> .\minebak WorldName --no-input

# Helpful output for the user can be disabled (ex. progress bars)
> .\minebak WorldName --quiet

# You can save the backup anywhere
# By default a folder named like your world is created in the folder you call minebak from
# In this example the backup is saved into the "MyBackup" Folder
> .\minebak WorldName --output D:\anton\Download\MyBackup 
> .\minebak WorldName --output ./MyBackup # relative links are supported

# You can add timestamps to your backup folder name to differentiate between multiple backups
# The result is a folder named ./MyBackup20230921
# The number is a date with the format YYYY-MM-DD
> .\minebak WorldName --output ./MyBackup --with-timestamp
> .\minebak WorldName --with-timestamp # Also valid without the --output flag

# Easy to use in scripts
# You can save your password in a file. The files needs to contain only one line: the password  
> .\minebak WorldName `
      --addr 127.0.0.1 `
      --port 21 `
      --user youremail@example.com `
      --password-file .\pass.txt `
      --with-timestamp
```
### Linux
```shell
# To use minebak use it from the terminal
# WorldName is a required argument, it refers to your minecraft's world name on the server
$ ./minebak WorldName

# Interactive by default
# minebak will ask you the MultiCraft ip, port, username and password
$ ./minebak WorldName
# Prompt for missing information
$ ./minebak WorldName --addr 127.0.0.1 # missing port, username and password

# But interactive input can be disabled
$ ./minebak WorldName --no-input

# Helpful output for the user can be disabled (ex. progress bars)
$ ./minebak WorldName --quiet

# You can save the backup anywhere
# By default a folder named like your world is created in the folder you call minebak from
# In this example the backup is saved into the "MyBackup" Folder
$ ./minebak WorldName --output /home/antonio/Downloads/MyBackup
$ ./minebak WorldName --output ./MyBackup # relative links are supported

# You can add timestamps to your backup folder name to differentiate between multiple backups
# The result is a folder named ./MyBackup20230921
# The number is a date with the format YYYY-MM-DD
$ ./minebak WorldName --output ./MyBackup --with-timestamp
$ ./minebak WorldName --with-timestamp # Also valid without the --output flag

# Easy to use in scripts
# You can save your password in a file. The files needs to contain only one line: the password  
$ ./minebak WorldName \
      --addr 127.0.0.1 \
      --port 21 \
      --user youremail@example.com \
      --password-file ./pass.txt \
      --with-timestamp
```
## Licence
This project follows the MIT licence, see LICENSE
