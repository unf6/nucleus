package installer

import (
	"fmt"
	"os"
	"os/exec"
	"time" 
	"path/filepath" 
	"github.com/unf6/nucleus/internal/config"
	spinner "gabe565.com/spinners"
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
	"wl-color-picker",
	"imagemagick",
	"qt6-svg",
	"hyprsunset",
	"jq",
	"kitty",
	"kvantum",
	"kvantum-qt5",
	"matugen-git",
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

	preReq := exec.Command("sudo", "pacman", "-S", "--needed", "--noconfirm", "base-devel", "git")
	preReq.Stdout = os.Stdout
	preReq.Stderr = os.Stderr
	if err := preReq.Run(); err != nil {
		return fmt.Errorf("failed to install prerequisites: %w", err)
	}

	tmpDir := "/tmp/yay-bin"
	_ = os.RemoveAll(tmpDir)

	clone := exec.Command("git", "clone", "https://aur.archlinux.org/yay.git", tmpDir)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	if err := clone.Run(); err != nil {
		return fmt.Errorf("failed to clone yay: %w", err)
	}

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
	
	args := append([]string{"-S", "--needed", "--noconfirm"}, Dependencies...)
	cmd := exec.Command("yay", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func CopyToQuickShellConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	repoRoot, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	src := filepath.Join(repoRoot, "quickshell", "nucleus-shell")
	dst := filepath.Join(home, ".config", "quickshell", "nucleus-shell")

	// Sanity check
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("source quickshell config not found: %s", src)
	}

	fmt.Printf("\nCopying %s → %s\n", src, dst)

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	
	_ = os.RemoveAll(dst)

	cmd := exec.Command("cp", "-r", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}



func RunWithSpinner(label string, fn func() error) error {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	sp := spinner.Dots
	frame := 0

	// Hide cursor
	fmt.Print("\x1B[?25l")
	defer fmt.Print("\x1B[?25h")

	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	for {
		select {
		case err := <-done:
			fmt.Printf("\r\x1B[K✓ %s\n", label)
			return err

		default:
			fmt.Printf(
				"\r\x1B[K%s %s",
				sp.Frames[frame],
				label,
			)
			frame = (frame + 1) % len(sp.Frames)
			time.Sleep(sp.Interval)
		}
	}
}
