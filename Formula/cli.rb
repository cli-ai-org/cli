# Homebrew Formula for cli
class Cli < Formula
  desc "Discover and explore CLI tools on your system - optimized for AI agents"
  homepage "https://github.com/cli-ai-org/cli"
  url "https://github.com/cli-ai-org/cli/archive/refs/tags/v0.1.0-beta.tar.gz"
  sha256 "a43f5c687c8e55e675db74f397f45e8bf48c5af811f48d16c5f35faf37848035"
  license "MIT"
  version "0.1.0-beta"

  depends_on "go" => :build

  def install
    # Build from source
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"cli", "."

    # Install documentation
    doc.install "README.md"
    doc.install "COMMANDS.md"
    doc.install "USAGE_EXAMPLES.md"
    doc.install "PACKAGE_DETECTION.md"
    doc.install "docs/AI_AGENT_USAGE.md" if File.exist?("docs/AI_AGENT_USAGE.md")
  end

  test do
    # Test that the binary runs
    system "#{bin}/cli", "help"

    # Test list command
    assert_match "Found", shell_output("#{bin}/cli list 2>&1")
  end
end
