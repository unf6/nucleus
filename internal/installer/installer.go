package installer

import (
	"fmt"
	"os"
	"os/exec"
)
var Dependencies = []string{
	"bluez-utils",
	"brightnessctl",
	"ddcutil",
	"fastfetch",
	"firefox",
	"grim",
	"hyprland",
	"hyprlock",
	"hyprpaper",
	"hyprpicker",
	"hyprsunset",
	"jq",
	"kitty",
	"kvantum",
	"kvantum-qt5",
	"matugen-bin",
	"nautilus",
	"nerd-fonts",
	"networkmanager",
	"papirus-icon-theme-git",
	"playerctl",
	"qt5-graphicaleffects",
	"qt5-wayland",
	"qt5ct",
	"qt6-5compat",
	"qt6-wayland",
	"qt6ct",
	"quickshell",
	"slurp",
	"starship",
	"ttf-fira-code",
	"ttf-fira-sans",
	"ttf-firacode-nerd",
	"ttf-font-awesome",
	"ttf-jetbrains-mono",
	"ttf-material-symbols-variable-git",
	"wf-recorder",
	"wireplumber",
	"xdg-desktop-portal-hyprland",
	"zenity",
}

func ensureYayInstalled() error {
	if _, err := exec.LookPath("yay"); err == nil {
		fmt.Println("✓ yay is already installed")
		return nil
	}

	fmt.Println("yay not found, installing yay...")

	// Install prerequisites
	preReq := exec.Command("sudo", "pacman", "-S", "--needed", "--noconfirm", "base-devel", "git")
	preReq.Stdout = os.Stdout
	preReq.Stderr = os.Stderr
	if err := preReq.Run(); err != nil {
		return fmt.Errorf("failed to install prerequisites: %w", err)
	}

	// Clone yay-bin
	tmpDir := "/tmp/yay-bin"
	_ = os.RemoveAll(tmpDir)

	clone := exec.Command("git", "clone", "https://aur.archlinux.org/yay.git", tmpDir)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	if err := clone.Run(); err != nil {
		return fmt.Errorf("failed to clone yay: %w", err)
	}

	// Build & install
	build := exec.Command("bash", "-c", fmt.Sprintf("cd %s && makepkg -si --noconfirm", tmpDir))
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("failed to build yay: %w", err)
	}

	fmt.Println("✓ yay installed successfully")
	return nil
}


func InstallDependencies() error {

    if err := ensureYayInstalled(); err != nil {
		return err
	}
	
	args := append([]string{"yay", "-S", "--needed", "--noconfirm"}, Dependencies...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
