# Hosting Local Share on Homebrew

This guide explains how to host the Local Share application on Homebrew as a single binary.

## Steps to Publish to Homebrew

### 1. Create a GitHub Repository

If you haven't already, push your Local Share code to GitHub:

```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/local-share.git
git push -u origin main
```

### 2. Tag a Release

Create a tagged release on GitHub:

```bash
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

Then go to GitHub and create a new release based on this tag. Upload the compiled binary or create a source tarball.

### 3. Update the Formula File

Update the `local-share.rb` formula with:

- Your actual GitHub username
- The correct SHA256 checksum of your release tarball
- The correct license (if not MIT)

To get the SHA256 checksum:

```bash
curl -L https://github.com/yourusername/local-share/archive/v1.0.0.tar.gz | shasum -a 256
```

### 4. Test Your Formula Locally

```bash
brew install --build-from-source ./local-share.rb
```

### 5. Submit to Homebrew

You have two options:

#### Option A: Submit to Homebrew Core (Official Repository)

For widely-used applications that meet Homebrew's requirements:

1. Fork the Homebrew Core repository
2. Add your formula to `Formula/l/local-share.rb`
3. Submit a pull request

#### Option B: Create Your Own Tap (Recommended for Personal Projects)

```bash
# Create a GitHub repository for your tap
brew tap-new yourusername/tap
cp local-share.rb $(brew --repo)/Library/Taps/yourusername/homebrew-tap/Formula/
cd $(brew --repo)/Library/Taps/yourusername/homebrew-tap
git add Formula/local-share.rb
git commit -m "Add local-share formula"
git push
```

Users can then install your app with:

```bash
brew tap yourusername/tap
brew install local-share
```

## Usage After Installation

Once installed via Homebrew, users can use the application with:

```bash
# Start the server
local-share receiver

# Send text
local-share send text <server-ip> "Your message"

# Send file
local-share send file <server-ip> /path/to/file

# Get help
local-share help
```

## Updating Your Formula

When you release a new version:

1. Create a new tagged release on GitHub
2. Update the formula with the new version and SHA256
3. Push the changes to your tap repository

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Creating a Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Homebrew Documentation](https://docs.brew.sh/) 