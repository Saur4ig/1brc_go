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