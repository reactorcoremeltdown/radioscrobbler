class Radioscrobbler < Formula
  desc "A source-agnostic song scrobbler"
  homepage "https://github.com/reactorcoremeltdown/radioscrobbler"
  license "GPLv3"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://tiredsysadmin.cc/media/files/bin/radioscrobbler/radioscrobbler-Darwin-arm64"

      def install
        bin.install "radioscrobbler-Darwin-arm64" => "radioscrobbler"
      end
    end

    if Hardware::CPU.intel?
      url "https://tiredsysadmin.cc/media/files/bin/radioscrobbler/radioscrobbler-Darwin-x86_64"

      def install
        bin.install "radioscrobbler-Darwin-x86_64" => "radioscrobbler"
      end
    end
  end

end
