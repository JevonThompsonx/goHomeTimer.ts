# Go Home Timer
---

## Overview

Go Home Timer is a simple, interactive tool designed to help you determine your departure time based on your arrival and break schedule, then automatically add an "out for the day" event to your Google Calendar. This project aims to streamline the process of tracking your work hours and ensuring accurate calendar entries.

## Features

*   **Interactive Time Input:** Easily input your arrival time and break duration.
*   **Calculated Departure Time:** Automatically determine your recommended departure time based on a standard 8-hour workday.
*   **Google Calendar Integration:** Seamlessly adds an "Out for the day" event to your primary Google Calendar.

## Getting Started

### Prerequisites

*   Go (version 1.18 or later)
*   A Google account

### Installation and Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/JevonThompsonx/goHomeTimer.ts.git
    cd goHomeTimer.ts
    ```

2.  **Enable the Google Calendar API and get credentials:**

    *   Go to the [Google Cloud Console](https://console.cloud.google.com/).
    *   Create a new project.
    *   Enable the **Google Calendar API**.
    *   Create an **OAuth 2.0 Client ID** for a **Desktop app**.
    *   Download the JSON file and rename it to `credentials.json`.
    *   Place the `credentials.json` file in the root of the project directory.

3.  **Build the application:**

    ```bash
    go build -o goHomeTimer
    ```

## Usage

1.  **Run the application:**

    ```bash
    ./goHomeTimer
    ```

2.  **Enter your arrival time** when prompted (e.g., `9:00 AM`).

3.  **Enter your break duration** (e.g., `0.5` for 30 minutes, `1` for 1 hour).

4.  **Authorize the application:**

    *   The first time you run the application, you will be prompted to authorize access to your Google Calendar.
    *   Copy the URL from the terminal and paste it into your browser.
    *   Sign in to your Google account and grant the application permission.
    *   You will be redirected to a page that may show an error or a default web server page. This is expected.
    *   Copy the `code` from the URL in your browser's address bar (it will look something like `4/0AVMBs...`).
    *   Paste the code back into the terminal.

5.  **Event Creation:**

    *   The application will calculate your departure time and create an "Out for the day" event in your primary Google Calendar.
    *   A link to the event will be printed in the terminal.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---