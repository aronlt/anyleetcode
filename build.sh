#! /bin/bash
set -e
cd app && fyne package -os darwin -icon ../static/LeetCode.png && mv app.app ../output/leetcode.app && rm app

