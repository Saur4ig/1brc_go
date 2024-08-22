# The One Billion Row Challenge
My attempt to solve the **One Billion Row Challenge**.

## Overview

### v1
- **Time**: Total: 388.06s
- **Logic**: The most basic "read a file" line by line

### v1.1
- **Time**: Total: 370.11s
- **Logic**: added buffer size to scanner

### v2
- **Time**: Total: 329.07s
- **Logic**: added decode logic using bytes.IndexByte, looping - using buffer

### v2.1
- **Time**: Total: 308.56s
- **Logic**: pointer instead of value for result map

### v2.2
- **Time**: Total: 244.13s
- **Logic**: replaced float parsing with custom function

### v3
- **Time**: Total: 58.30s
- **Logic**: parallel processing of files, but with scanner and maps merging

### v3.1
- **Time**: Total: 32.77s
- **Logic**: replaced bufio scanner with faster implementation from v2

### v3.2
- **Time**: Total: 23.92s
- **Logic**: added capacity to the line slice, avoiding slice grow

### v3.3
- **Time**: Total: 20.32s
- **Logic**: replaced convert to string with fast string

### v3.4
- **Time**: Total: 18.57s
- **Logic**: instead of re-writing max or min each time, now rewrite only on change

### v3.5
- **Time**: Total: 15.45s
- **Logic**: reduced amount of interactions with the map

### v3.6
- **Time**: Total: 15.00s
- **Logic**: added map capacity, to reduce the alloc time

### v3.7
- **Time**: Total: 14.00s
- **Logic**: refactored map insertion logic, to reduce amount of allocations

### v3.8
- **Time**: Total: 13.12s
- **Logic**: custom semicolon search function, searching it from the right side

### v3.9
- **Time**: Total: 12.47
- **Logic**: reduced amount of iterations by 3 for semicolon search

### v3.10
- **Time**: Total: 11.08
- **Logic**: instead of converting []bytes to string all the time(quite expensive), moved to uint64 as a key in map