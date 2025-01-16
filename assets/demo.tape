# `stylish` demo
Output demo.gif
Output demo.mp4

# Settings

Set Width 700 
Set Height 1200
Set TypingSpeed 100ms

# Require `stylish` and its dependency
Require dircolors
Require stylish

# Hidden initialization
Hide

# Ensure program is up to date
Type "go install .."
Enter
Sleep 100ms

# Remove demo theme if it already exists
# Type "rm ~/.config/stylish/DemoTheme -rf"
# Enter

# Setup demo dir that has example files
Type "mkdir DemoDir"
Enter
Type "cd DemoDir"
Enter
Type "touch Nickleback.mp3 MyCoolRingtone.wav SlamPoetryDemo.ogg"
Enter
Type "touch my_first_wedding.mp4 my_seventh_wedding.avi"
Enter
Type "touch config.conf poketerm.ini jason.json mark_hamill.yaml tommy_lee.toml"
Enter
Type "cd .."
Enter

# Clear any existing styling
Type 'export LS_COLORS=""'
Enter
Type "clear"
Enter

# Get back to the action
Show

# Show boring files
Type "ls DemoDir --color=auto"
Enter
Sleep 2000ms

# Start program
Type "clear"
Enter
Type "stylish"
Enter
Sleep 2000ms

# Create new theme called "DemoTheme"
Type "n"
Sleep 250ms
Type "DemoTheme"
# Sleep 1000ms
Enter
Sleep 600ms

# Create new style called "Audio Files"
Type "n"
Sleep 400ms
Type "Audio Files"
# Sleep 1000ms
Ctrl+S
Sleep 600ms

# Make "Audio Files" bold
Type "1"
Sleep 1000ms

# Set foreground to #E69F66, leave background default
Type "f"
Sleep 250ms
Type "E69F66"
Sleep 800ms
Ctrl+S
Sleep 800ms

# Set filetypes
Type "t.mp3"
Enter
Type ".ogg"
Enter
Type ".wav"
Sleep 600ms
Ctrl+S
Sleep 600ms

# Create new style called "Video Files"
Type "n"
Sleep 500ms
Type "Video Files"
Ctrl+S
# Type "j"
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
Ctrl+S
# Type "j"
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
Sleep 750ms
Ctrl+S
Sleep 800ms

Sleep 1000ms
Ctrl+C
Sleep 400ms
Ctrl+C
Sleep 200ms

# Sleep 200ms
# Type "stylish apply DemoTheme"
# Enter

# Sleep 2000ms

#Type "eval $(stylish apply DemoTheme)"
#Enter
#Type "ls DemoDir --color=auto"
Type "cd DemoDir" Enter
Type "stylish preview DemoTheme" Enter
Sleep 4000ms

# Execution cleanup
Hide
Type "rm -rf ~/.config/stylish/DemoTheme"
Enter
Type "cd .."
Enter
Type "rm -rf DemoDir"
Enter
