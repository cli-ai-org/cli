# Quick Publish to Homebrew - 5 Minutes

Fast track to get cli on Homebrew.

## Prerequisites

- GitHub account
- `gh` CLI installed (`brew install gh`)

## Steps

### 1. Create GitHub Repository (1 min)

```bash
cd /Users/op/code/cli

# Login to GitHub
gh auth login

# Create repo
gh repo create cli --public --description "Discover CLI tools on your system - optimized for AI agents"

# Push code
cd cli
git init
git add .
git commit -m "Initial release of cli"
git branch -M main
git remote add origin https://github.com/$(gh api user --jq .login)/cli.git
git push -u origin main
```

### 2. Create Release (1 min)

```bash
# Tag the release
git tag -a v0.1.0 -m "v0.1.0 - Initial release"
git push origin v0.1.0

# Create GitHub release
gh release create v0.1.0 \
  --title "v0.1.0 - Initial Release" \
  --notes "Initial release of cli

Features:
- Discover CLI tools across your PATH
- List packages (npm, pip, brew, cargo, gem)
- Export JSON catalog for AI agents
- Link tools to their source packages"
```

### 3. Calculate SHA256 (30 sec)

```bash
# Get your GitHub username
USERNAME=$(gh api user --jq .login)

# Download and calculate SHA
curl -L "https://github.com/$USERNAME/cli/archive/refs/tags/v0.1.0.tar.gz" -o /tmp/cli.tar.gz
shasum -a 256 /tmp/cli.tar.gz

# Copy the hash (first part before the filename)
```

### 4. Create Homebrew Tap (2 min)

```bash
# Create tap repository
gh repo create homebrew-cli --public --description "Homebrew tap for cli"

# Clone it
cd ~
git clone "https://github.com/$USERNAME/homebrew-cli.git"
cd homebrew-cli

# Copy formula template
cp /Users/op/code/cli/cli/Formula/cli.rb ./cli.rb

# Edit the formula (replace USERNAME and SHA256)
cat > cli.rb << 'EOF'
class Cliil < Formula
  desc "Discover and explore CLI tools on your system - optimized for AI agents"
  homepage "https://github.com/REPLACE_USERNAME/cli"
  url "https://github.com/REPLACE_USERNAME/cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_SHA256"
  license "MIT"
  version "0.1.0"

  depends_on "go" => :build

  def install
    cd "cli" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"cli", "."
    end
  end

  test do
    system "#{bin}/cli", "help"
  end
end
EOF

# Now edit with real values
echo "Edit cli.rb now:"
echo "1. Replace REPLACE_USERNAME with: $USERNAME"
echo "2. Replace REPLACE_WITH_SHA256 with the hash from step 3"
echo ""
read -p "Press enter when done editing..."

# Commit and push
git add cli.rb
git commit -m "Add cli formula"
git push origin main
```

### 5. Test Installation (1 min)

```bash
# Tap your repository
brew tap $USERNAME/cli

# Install cli
brew install cli

# Test it
cli help
cli list | head -20
```

## Done! 🎉

Users can now install with:

```bash
brew tap YOUR_USERNAME/cli
brew install cli
```

Or in one command:
```bash
brew install YOUR_USERNAME/cli/cli
```

## Share With Users

Add this to your README.md:

```markdown
## Installation

### Homebrew (macOS/Linux)

\`\`\`bash
brew install YOUR_USERNAME/cli/cli
\`\`\`

### From Source

\`\`\`bash
git clone https://github.com/YOUR_USERNAME/cli.git
cd cli/cli
go build -o cli .
\`\`\`
```

## Updating for New Releases

When you release v0.2.0:

```bash
# 1. Tag and release
cd /Users/op/code/cli/cli
git tag v0.2.0
git push origin v0.2.0
gh release create v0.2.0 --generate-notes

# 2. Calculate new SHA
USERNAME=$(gh api user --jq .login)
curl -L "https://github.com/$USERNAME/cli/archive/refs/tags/v0.2.0.tar.gz" | shasum -a 256

# 3. Update formula
cd ~/homebrew-cli
# Edit cli.rb:
# - Change version to "0.2.0"
# - Change url to v0.2.0
# - Update sha256
git add cli.rb
git commit -m "Update cli to 0.2.0"
git push

# 4. Users update with:
# brew update
# brew upgrade cli
```

## Troubleshooting

### Can't find gh command
```bash
brew install gh
gh auth login
```

### SHA256 mismatch
```bash
# Clear cache and try again
brew cleanup
brew untap YOUR_USERNAME/cli
brew tap YOUR_USERNAME/cli
brew install cli
```

### Formula not found
```bash
# Make sure repo is named "homebrew-cli"
gh repo view YOUR_USERNAME/homebrew-cli

# Update tap
brew update
brew tap YOUR_USERNAME/cli
```
