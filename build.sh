#! /bin/bash
set -e
mkdir -p output
cd app && fyne package -os darwin -icon ../static/LeetCode.png && mv app.app ../output/leetcode.app

