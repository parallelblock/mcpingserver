#!/usr/bin/python3

import os
import subprocess

BUILD_COMBINATIONS = [
    ("darwin", "amd64"), ("darwin", "386"),
    ("linux", "amd64"), ("linux", "386"), ("linux", "arm"), ("linux", "arm64"),
    ("windows", "amd64"), ("windows", "386")]

BUILD_DIR = "bin/"
if not os.path.exists(BUILD_DIR):
    os.makedirs(BUILD_DIR)

for combo in BUILD_COMBINATIONS:
    build_location = BUILD_DIR + "mcpingserver-{}-{}".format(*combo)
    if combo[0] is "windows":
        build_location += ".exe" # Gotta love windows!
    print("Building for {} {} - {}".format(*combo, build_location))
    build_cmd = "env GOOS={} GOARCH={} go build -o {} .".format(*combo, build_location)
    subprocess.check_call(build_cmd, shell=True)
