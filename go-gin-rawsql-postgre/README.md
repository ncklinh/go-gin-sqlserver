# 🎬 Film Rental API

A backend service for managing film rentals, built with Go (Gin), PostgreSQL, and Docker, with hot-reloading support via Air. Optimized for development using WSL2 on Ubuntu.

-----
- To test concurrency of MQTT subscriber, run /test/mqtt_concurrency.go
-----

## 🚀 Prerequisites

- Windows 10/11 with WSL2
- Ubuntu (via WSL)
- Docker Desktop with WSL2 integration
- VS Code with **Remote - WSL** extension
- Go (version 1.18 or later)
- Air (for hot-reloading)

---

## ✅ Setup Instructions

### 1. Set Up WSL2 with Ubuntu

- Open PowerShell (Admin) and install WSL2:
  ```powershell
  wsl --install
  ```
- Reboot if prompted.
- Verify WSL2 is the default version:
  ```powershell
  wsl --status
  ```
  Ensure "Default Version: 2" is displayed.
- If Ubuntu is not installed, run:
  ```powershell
  wsl --install -d Ubuntu
  ```
- On first Ubuntu launch, set up your username and password.

### 2. Install and Configure Docker Desktop

- Download and install from: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
- Enable WSL2 integration in Docker Desktop settings for Ubuntu.
- Verify Docker installation in Ubuntu terminal:
  ```bash
  docker --version
  ```

### 3. Install Go and Air in WSL2

In Ubuntu terminal:

```bash
sudo apt update
sudo apt install golang-go
go install github.com/air-verse/air@latest
```

- Add Air to PATH:
  ```bash
  echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
  source ~/.bashrc
  ```
- Verify installations:
  ```bash
  go version
  air -v
  ```

### 4. Move Project to WSL

To avoid performance issues, store your project in WSL:

- Create a projects folder:
  ```bash
  mkdir -p ~/projects
  ```
- Copy your project from Windows (e.g., `D:\path\to\project`) to `\\wsl$\Ubuntu\home\<yourusername>\projects` using Windows Explorer.
- Replace `<yourusername>` with your WSL username.

### 5. Configure Air for Hot-Reloading

In your project folder, create or update `.air.toml`:

```toml
[build]
cmd = "go build -o tmp/main -mod=mod -buildvcs=false ."
bin = "tmp/main"
full_bin = "tmp/main"
```

Leave other `.air.toml` settings unchanged. This ensures compatibility with WSL/Linux.

---

## 🏃 Running the Project

1. **Open Project in VS Code**

   - In Ubuntu terminal:
     ```bash
     cd ~/projects/your-project-folder
     code .
     ```
   - Confirm "WSL: Ubuntu" appears in the bottom-left corner of VS Code.

2. **Start PostgreSQL**

   - Run Docker Compose:
     ```bash
     bash ./scripts/start-dev.sh
     ```
   - Verify container is running:
     ```bash
     docker ps
     ```

3. **Start the API with Air**
   - Run:
     ```bash
     air
     ```
   - Save any `.go` file to trigger automatic rebuild and restart.

---

## 🛠 Troubleshooting

- **Air not found**:

  ```bash
  go install github.com/air-verse/air@latest
  ```

  Recheck PATH with `echo $PATH`.

- **Docker containers not starting**:
  Verify with:

  ```bash
  docker ps
  ```

  Ensure containers are running and healthy.

- **Performance issues**:
  Always edit code in the WSL folder (`~/projects`), not directly on the Windows `C:\` drive, to avoid permission and performance problems.

