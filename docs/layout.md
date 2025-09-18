# UI Layout Specification

This document describes the visual layout and component placement for the Moodle Prototype Manager desktop application.

---

## 1. Main Window

### 1.1 Header
- **Logo**
  Centered at the top of the window.  with size 128x128 pixel
- **Version Label**
  Directly beneath the logo, displays the app name and version (e.g. “Moodle Prototype v0.0.1”).

### 1.2 Primary Action Area
- **Run/Stop Button**
  - When no container is running: a large green “Run Moodle” button.
  - Once the container is up: the button changes to a red “Stop Moodle” label.
  - This button is disabled until health checks (Docker + Internet) pass.

### 1.3 Credentials Display
- Appears after the container is up and running.
- Shows three lines of information in monospace or highlighted text:
  - `username: admin`
  - `password: <generated-password>`
  - `url: http://localhost:8080`
- Beneath the credentials, a styled “Open Browser” button invites the user to launch their default browser at the Moodle URL.

---

## 2. Footer Status Bar

- **Left Side**
  Displays a status message:
  - While checking: “Checking…”
  - After checks complete: either “Ready” or an error message.
- **Right Side**
  Two status icons, each with a colored circle (●):
  - Docker status: green if the Docker daemon is available, red otherwise.
  - Internet status: green if a quick connectivity check succeeds, red otherwise.

---

## 3. Progress Modal (Image Download)

- **Title Text**
  “Downloading Image…”
- **Progress Bar**
  A horizontal bar filling from 0% to 100%.
- **Percentage Labels**
  - Left: “0%”
  - Right: “100%”

This modal overlays the main window during the `docker pull` operation to show live progress feedback.

---

## 4. Startup-Wait Modal

- **Spinner or Loading Indicator**
  Communicates that the app is waiting for Moodle to finish initializing.
- No user interaction—automatically closes when credentials are detected in the container logs.

---

## 5. Visual Style Notes

- Use a clean, minimal palette to keep the focus on status and action items.
- Buttons use high-contrast borders and subtle fill to indicate states:
  - Enabled (solid color), Disabled (greyed out).
- Status circles use flat red/green colors for immediate recognition.
- All text should be legible at standard desktop resolutions; use a monospaced font for credentials.

---

_End of layout specification._
