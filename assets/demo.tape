# `stylish` demo
Output demo.gif

# Settings

Set Width 700 
Set Height 1050
Set TypingSpeed 100ms

# Require `stylish` and its dependency
Require dircolors
Require stylish

Hide
Type "go install .."
Enter
Sleep 100ms
Type "rm ~/.config/stylish/DemoTheme -rf"
Enter
Type "clear"
Enter
Show

# Start program
Type "stylish"
Enter
Sleep 600ms

# Create new theme called "DemoTheme"
Type "n"
Sleep 250ms
Type "DemoTheme"
# Sleep 1000ms
Enter
Sleep 600ms

# Create new style called "Audio Files"
Type "n"
Sleep 800ms
Type "Audio Files"
# Sleep 1000ms
Ctrl+S
Sleep 600ms

# Make "Audio Files" bold
Type "1"
Sleep 1000ms

# Set foreground to #E69F66, leave background default
Type "fE69F66"
Sleep 1000ms
Ctrl+S
Sleep 800ms

# Set filetypes
Type "t.mp3"
Enter
Type ".ogg"
Enter
Type ".wav"
Sleep 800ms
Ctrl+S
Sleep 800ms

# Create new style called "Video Files"
Type "n"
Sleep 500ms
Type "Video Files"
Sleep 800ms
Ctrl+S
Type "j"
Sleep 800ms

# Make "Video Files" underlined
Type "2"
Sleep 800ms

# Set background to #AE81FF, leave foreground default
Type "bAE81FF"
Sleep 800ms
Ctrl+S
Sleep 800ms

# Set filetypes
Type "t.avi"
Enter
Type ".mp4"
Sleep 800ms
Ctrl+S
Sleep 800ms

# Create new style called "Config Files"
Type "n"
Sleep 500ms
Type "Config Files"
Sleep 800ms
Ctrl+S
Type "j"
Sleep 800ms

# Make "Config Files" bold and underlined (vhs doesn't render blinking)
Type "12"
Sleep 800ms

# Set background to #75715E, set foreground to #F92672
Type "b75715E"
Ctrl+S
Sleep 800ms
Type "fF92672"
Sleep 800ms
Ctrl+S
Sleep 800ms

# Set filetypes
Type "t.conf"
Enter
Type ".yaml"
Enter
Type ".json"
Enter
Type ".toml"
Enter
Type ".ini"
Sleep 800ms
Ctrl+S
Sleep 800ms

Sleep 2000ms
Ctrl+C
Sleep 800ms
Ctrl+C

Sleep 200ms
Type "stylish apply DemoTheme"
Enter

Sleep 3000ms
