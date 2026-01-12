# Publishing cli to Homebrew

Complete guide to making cli installable via `brew install cli`.

## Two Approaches

### Option 1: Custom Tap (Recommended)
Users install with: `brew install cli-org/cli/cli`

**Pros:**
- You control the release cycle
- Easier to set up and maintain
- No external review process
- Can publish immediately

**Cons:**
- Users need to specify the tap name
- Less discoverable than official Homebrew

### Option 2: Official Homebrew Core
Users install with: `brew install cli`

**Pros:**
- Most discoverable (just `brew install cli`)
- Appears in `brew search`
- More "official" presence

**Cons:**
- Requires review by Homebrew maintainers
- Must meet strict guidelines
- Takes longer to get approved
- More maintenance requirements

---

## Option 1: Custom Tap Setup (Step-by-Step)

### Prerequisites

1. **GitHub Account** - You need a GitHub account
2. **GitHub Repository** - Create `https://github.com/cli-org/cli` with your code
3. **GitHub Token** - Create a personal access token for releases

### Step 1: Prepare Your Repository

```bash
cd /Users/op/code/cli

# Initialize git if not already done
git init
git add .
git commit -m "Initial commit"

# Create GitHub repo and push
# Go to github.com and create a new repo called "cli"
git remote add origin https://github.com/YOUR_USERNAME/cli.git
git branch -M main
git push -u origin main
```

### Step 2: Create a GitHub Release

#### Option A: Using GitHub CLI (gh)

```bash
# Install gh if needed
brew install gh

# Login
gh auth login

# Create a release
cd /Users/op/code/cli/cli
go build -o cli .

# Create tarball
tar -czf cli-0.1.0.tar.gz .

# Create GitHub release
gh release create v0.1.0 cli-0.1.0.tar.gz \
  --title "v0.1.0" \
  --notes "Initial release of cli - CLI tool discovery for AI agents"
```

#### Option B: Using GitHub Website

1. Go to `https://github.com/YOUR_USERNAME/cli/releases/new`
2. Tag: `v0.1.0`
3. Title: `v0.1.0`
4. Description: "Initial release"
5. Upload the tarball
6. Click "Publish release"

### Step 3: Calculate SHA256

```bash
# Download the release tarball
curl -L https://github.com/YOUR_USERNAME/cli/archive/refs/tags/v0.1.0.tar.gz -o cli-0.1.0.tar.gz

# Calculate SHA256
shasum -a 256 cli-0.1.0.tar.gz
# Copy this hash!
```

### Step 4: Create Homebrew Tap Repository

```bash
# Create a new GitHub repo called "homebrew-cli"
# (Homebrew taps must be named "homebrew-*")

mkdir ~/homebrew-cli
cd ~/homebrew-cli

# Copy the formula
cp /Users/op/code/cli/cli/Formula/cli.rb ./cli.rb

# Edit cli.rb - update the SHA256 and username
nano cli.rb
# Change:
# - url to your GitHub username
# - sha256 to the calculated hash
# - version if needed

# Example cli.rb content:
```

```ruby
class Cliil < Formula
  desc "Discover and explore CLI tools on your system - optimized for AI agents"
  homepage "https://github.com/YOUR_USERNAME/cli"
  url "https://github.com/YOUR_USERNAME/cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "abc123..." # <- Your calculated SHA256
  license "MIT"
  version "0.1.0"

  depends_on "go" => :build

  def install
    cd "cli" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"cli", "."
    end

    doc.install "cli/README.md"
    doc.install "cli/COMMANDS.md"
  end

  test do
    system "#{bin}/cli", "help"
    assert_match "Found", shell_output("#{bin}/cli list 2>&1")
  end
end
```

```bash
# Push to GitHub
git init
git add cli.rb
git commit -m "Add cli formula"
git remote add origin https://github.com/YOUR_USERNAME/homebrew-cli.git
git branch -M main
git push -u origin main
```

### Step 5: Test Installation

```bash
# Tap your repository
brew tap YOUR_USERNAME/cli

# Install cli
brew install cli

# Test it
cli help
cli list
```

### Step 6: Users Install With

```bash
# Add your tap
brew tap YOUR_USERNAME/cli

# Install
brew install cli
```

Or in one command:
```bash
brew install YOUR_USERNAME/cli/cli
```

---

## Option 1B: Automated with GoReleaser

GoReleaser automates building for multiple platforms and publishing.

### Setup GoReleaser

```bash
# Install goreleaser
brew install goreleaser

# Test the configuration
cd /Users/op/code/cli/cli
goreleaser check

# Create a GitHub token
# Go to: https://github.com/settings/tokens/new
# Scopes needed: repo, write:packages
export GITHUB_TOKEN=your_token_here

# Dry run (test without publishing)
goreleaser release --snapshot --clean

# Real release
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

GoReleaser will:
- Build for macOS, Linux, Windows (amd64, arm64)
- Create GitHub release
- Upload binaries
- Update Homebrew tap automatically
- Generate checksums

---

## Option 2: Official Homebrew Core

For official Homebrew (advanced):

### Prerequisites

1. **Notable/stable project** - Homebrew prefers established projects
2. **Stable release** - At least one stable version
3. **CI/CD** - Automated tests are recommended
4. **License** - Must have an open source license

### Steps

1. **Meet Requirements**:
   - Have a stable v1.0.0+ release
   - Project on GitHub with good documentation
   - Passing tests
   - Active maintenance

2. **Fork homebrew-core**:
   ```bash
   # Fork https://github.com/Homebrew/homebrew-core
   git clone https://github.com/YOUR_USERNAME/homebrew-core.git
   cd homebrew-core
   ```

3. **Create Formula**:
   ```bash
   # Use brew create
   brew create --set-name cli https://github.com/YOUR_USERNAME/cli/archive/refs/tags/v1.0.0.tar.gz

   # Edit the generated formula
   brew edit cli
   ```

4. **Test Thoroughly**:
   ```bash
   brew install --build-from-source cli
   brew test cli
   brew audit --new-formula cli
   ```

5. **Submit PR**:
   ```bash
   git checkout -b cli
   git add Formula/cli.rb
   git commit -m "cli 1.0.0 (new formula)"
   git push origin cli
   # Create PR on GitHub
   ```

6. **Address Feedback**:
   - Homebrew maintainers will review
   - Address any issues they raise
   - May take days/weeks

---

## Updating Your Formula

When you release a new version:

### Custom Tap

```bash
cd ~/homebrew-cli

# Update cli.rb with new version and SHA256
nano cli.rb

git add cli.rb
git commit -m "Update cli to v0.2.0"
git push
```

Users update with:
```bash
brew update
brew upgrade cli
```

### With GoReleaser

```bash
cd /Users/op/code/cli/cli

# Create new tag
git tag v0.2.0
git push origin v0.2.0

# Release (automatically updates formula)
goreleaser release --clean
```

---

## Recommended Workflow

**For starting out:**

1. ✅ Create custom tap (`homebrew-cli`)
2. ✅ Set up GoReleaser for automation
3. ✅ Test with early users
4. ⏭️ Later: Submit to homebrew-core when stable

**Quick Start Commands:**

```bash
# 1. Create GitHub repo
gh repo create cli --public --source=. --remote=origin

# 2. Create release
git tag v0.1.0
git push origin v0.1.0

# 3. Calculate SHA256
curl -L https://github.com/YOUR_USERNAME/cli/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256

# 4. Create tap repo
gh repo create homebrew-cli --public

# 5. Update and push formula
# Edit Formula/cli.rb with correct URL and SHA256
cd Formula
git add cli.rb
git commit -m "Add cli formula"
git push

# 6. Test install
brew tap YOUR_USERNAME/cli
brew install cli
```

---

## Example: Real-World Tap

See how other projects do it:
- Vercel: `brew install vercel/tap/vercel-cli`
- Stripe: `brew install stripe/stripe-cli/stripe`
- HashiCorp: `brew install hashicorp/tap/terraform`

Your tap will be:
```bash
brew install cli-org/cli/cli
```

---

## Troubleshooting

### "Formula not found"
- Check repo name is `homebrew-cli`
- Check formula file is `cli.rb`
- Try `brew update`

### "Checksum mismatch"
- Recalculate SHA256
- Make sure URL is correct
- Clear cache: `brew cleanup`

### "Build failed"
- Test locally: `brew install --build-from-source --verbose cli`
- Check Go version requirements
- Verify source structure

### Test Formula Locally

```bash
# Install from local file
brew install --build-from-source ./Formula/cli.rb

# Or test directly
brew test cli
```

---

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GoReleaser Homebrew](https://goreleaser.com/customization/homebrew/)
- [Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [How to Create Homebrew Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
