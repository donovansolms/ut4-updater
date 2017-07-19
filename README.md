# Unreal Tournament 4 updater

An incremental updater for Unreal Tournament 4 on Linux

[![Build Status](https://travis-ci.org/donovansolms/ut4-updater.svg?branch=master)](https://travis-ci.org/donovansolms/ut4-updater)
[![Go Report Card](https://goreportcard.com/badge/github.com/donovansolms/ut4-updater)](https://goreportcard.com/report/github.com/donovansolms/ut4-updater)
[![Current Version](https://img.shields.io/badge/version-development-orange.svg)](https://img.shields.io/badge/version-development-orange.svg)

## About

Unreal Tournament 4 is available on Linux. Currently, no incremental updater exists which requires you to download the full game (around 10GB) for every patch release.

This updater aims to make life a bit simpler by only applying files that have been added, removed or modified since the previous version - allowing for much faster and smaller updates to the game.

## How it works

1. When the launcher (either [cli](https://github.com/donovansolms/ut4-launcher) or [GUI](https://github.com/donovansolms/ut4-launcher)) is opened the updater checks for a new release against [https://ut4.donovansolms.com](https://ut4.donovansolms.com) (the [ut4-update-server](https://github.com/donovansolms/ut4-update-server) is also open source)
2. If an update is available you can download and install, download and install in the background (while playing) or simply ignore
3. If you decide to install, the upgrader will create a clone of the current installation and apply the updates to the cloned version only.
4. The updater keeps track of installed versions. The option `version` allows you to specify the version to run, the default it to run the latest version available.

### Options

* `InstallPath` (required)

InstallPath is the base path for creating new installations. Must be specified. If the path doesn't exist, it will be created.

* `Versioning.Keep` (Defaults to 2 in the launcher)

Keep specifies the clones to keep. **Warning** if set to 0, the updates will be applied to your current version which could break the game and cause you to download the full game again.

* `Versioning.Run` (Defaults to latest in the launcher)

Run allows you to run any previously downloaded version. This is handy in case something is broken or you need to check performance between versions

* `SendStats` (Defaults to true in the launcher)

Basic information is collected to improve the updater and display stats about Unreal Tournament players using Linux

## GUI and CLI Launchers

* CLI Launcher: [ut4-launcher](https://github.com/donovansolms/ut4-launcher)

* GUI Launcher: [ut4-launcher-gui](https://github.com/donovansolms/ut4-launcher)

## Privacy

I respect your privacy. The updater collects the following information to improve the updater:

1. Your installed Unreal Tournament 4 versions
2. Your public IP is saved by the update server for country install stats
3. Kernel version and Distribution using `/etc/*-release` and `uname -r`. Only used for stats

You can disable Kernel version, distribution and installed Unreal Tournament version collection
by setting the `SendStats` option to `false`.

## Contact

You can get in contact on Twitter [@donovansolms](https://twitter.com/donovansolms) or by [creating an issue](https://github.com/donovansolms/ut4-updater/issues/new)
