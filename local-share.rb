class LocalShare < Formula
  desc "CLI app for sharing text and files between computers in a LAN with encryption"
  homepage "https://github.com/yourusername/local-share"
  url "https://github.com/yourusername/local-share/archive/v1.0.0.tar.gz"
  sha256 "YOUR_TARBALL_SHA256_CHECKSUM"
  license "MIT"
  
  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"local-share", "./cmd/..."
  end

  test do
    assert_match "local-share", shell_output("#{bin}/local-share --help", 0)
  end
end 