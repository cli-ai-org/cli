# Homebrew Formula for cli-ai
class CliAi < Formula
  desc "Discover and explore CLI tools on your system - optimized for AI agents"
  homepage "https://github.com/cli-ai-org/cli-ai"
  url "https://github.com/cli-ai-org/cli-ai/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"
  version "0.1.0"

  depends_on "go" => :build

  def install
    # Build from source
    cd "cli" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"cli-ai", "."
    end

    # Install documentation
    doc.install "cli/README.md"
    doc.install "cli/COMMANDS.md"
    doc.install "cli/USAGE_EXAMPLES.md"
    doc.install "cli/PACKAGE_DETECTION.md"
    doc.install "cli/docs/AI_AGENT_USAGE.md" if File.exist?("cli/docs/AI_AGENT_USAGE.md")
  end

  test do
    # Test that the binary runs
    system "#{bin}/cli-ai", "help"

    # Test list command
    assert_match "Found", shell_output("#{bin}/cli-ai list 2>&1")
  end
end
