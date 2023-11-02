package main

import (
    "fmt"
    "os"
    "os/exec"
)

const (
    RED     = "\033[91m"
    BLUE    = "\033[94m"
    GREEN   = "\033[92m"
    DEFAULT = "\033[0m"
)

func createCauldronFolder() {
    folderName := "Cauldron"
    folderPath := "C:\\" + folderName

    if _, err := os.Stat(folderPath); os.IsNotExist(err) {
        os.Mkdir(folderPath, 0755)
        fmt.Printf("Folder '%s' created at '%s'.\n", folderName, folderPath)
    } else {
        fmt.Printf("Folder '%s' already exists.\n", folderName)
    }

    cmd1 := exec.Command("icacls", "C:\\Cauldron", "/grant", "Everyone:(OI)(CI)F")
    cmd1.Run()
    cmd2 := exec.Command("net", "share", "Cauldron=C:\\Cauldron")
    cmd2.Run()
}

func banner() {
    asciiArt := `
    ███╗░░░███╗███████╗██████╗░██╗░░░░░██╗███╗░░██╗
    ████╗░████║██╔════╝██╔══██╗██║░░░░░██║████╗░██║
    ██╔████╔██║█████╗░░██████╔╝██║░░░░░██║██╔██╗██║
    ██║╚██╔╝██║██╔══╝░░██╔══██╗██║░░░░░██║██║╚████║
    ██║░╚═╝░██║███████╗██║░░██║███████╗██║██║░╚███║
    ╚═╝░░░░░╚═╝╚══════╝╚═╝░░╚═╝╚══════╝╚═╝╚═╝░░╚══╝
    `
    fmt.Print(asciiArt)
}

func logsMenu(sysmonLogs, securityLogs, systemLogs, applicationLogs string) (string, string, string, string) {
    for {
        clearScreen()
        banner()
        fmt.Printf("          +-----------SpellBook-----------+\n")
        fmt.Printf("          |  [1] Sysmon      |   %s    |\n", sysmonLogs)
        fmt.Printf("          |  [2] Security    |   %s    |\n", securityLogs)
        fmt.Printf("          |  [3] System      |   %s    |\n", systemLogs)
        fmt.Printf("          |  [4] Application |   %s    |\n", applicationLogs)
        fmt.Printf("          |  [B] Back        |  Main Menu |\n")
        fmt.Printf("          +-------------------------------+\n")

        toggle := getInput("Select Spell: ")

        switch toggle {
        case "1":
            sysmonLogs = toggleBoolean(sysmonLogs)
        case "2":
            securityLogs = toggleBoolean(securityLogs)
        case "3":
            systemLogs = toggleBoolean(systemLogs)
        case "4":
            applicationLogs = toggleBoolean(applicationLogs)
        case "b", "B":
            return sysmonLogs, securityLogs, systemLogs, applicationLogs
        default:
            fmt.Println("bruh")
        }
    }
}

func toggleBoolean(currentValue string) string {
    if currentValue == "\033[92mtrue \033[0m" {
        return "\033[91mfalse\033[0m"
    }
    return "\033[92mtrue \033[0m"
}

func generateClearLogsScript() {
    scriptContent := `
# PowerShell script to clear logs
# Clear Security, System, Application, and Sysmon logs

# Clear Security Log
Clear-EventLog -LogName Security

# Clear System Log
Clear-EventLog -LogName System

# Clear Application Log
Clear-EventLog -LogName Application

# Clear Sysmon Log (if installed)
$sysmonLog = Get-WinEvent -LogName "Microsoft-Windows-Sysmon/Operational" -MaxEvents 1 -ErrorAction SilentlyContinue
if ($sysmonLog) {
    Clear-EventLog -LogName "Microsoft-Windows-Sysmon/Operational"
}

# Display a message
Write-Host "Logs cleared successfully!"
`

    targetDirectory := "C:\\Cauldron"
    os.MkdirAll(targetDirectory, os.ModePerm)

    scriptPath := targetDirectory + "\\clearLogs.ps1"
    scriptFile, err := os.Create(scriptPath)
    if err != nil {
        fmt.Printf("Error creating 'clearLogs.ps1': %s\n", err)
        return
    }
    defer scriptFile.Close()

    _, err = scriptFile.WriteString(scriptContent)
    if err != nil {
        fmt.Printf("Error writing to 'clearLogs.ps1': %s\n", err)
        return
    }

    fmt.Printf("Created 'clearLogs.ps1' in the Cauldron directory.\n")
}

func launchSpell(sysmon, security, system, application string) {
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Printf("Error getting hostname: %s\n", err)
        return
    }
    if hostname[0] != '\\' {
        hostname = "\\" + hostname
    }

    securityGather := ""
    if security == "\033[92mtrue \033[0m" {
        securityGather = `
$securityLogPath = Join-Path -Path $desktop_folder -ChildPath "Security.evtx"
Write-Host "Copying Security log to $securityLogPath"
Copy-Item -Path "C:\Windows\System32\winevt\Logs\Security.evtx" -Destination $securityLogPath -Force
`
    }

    sysmonGather := ""
    if sysmon == "\033[92mtrue \033[0m" {
        sysmonGather = `
# Copy the Sysmon log to the destination folder
$sysmonLogPath = Join-Path -Path $filePath -ChildPath "Sysmon.evtx"
$sysmonLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\\Winevt\\Logs\\Microsoft-Windows-Sysmon%4Operational.evtx"

Write-Host "Copying Sysmon log to $sysmonLogPath"
Copy-Item -Path $sysmonLogSourcePath -Destination $sysmonLogPath -Force
`
    }

    systemGather := ""
    if system == "\033[92mtrue \033[0m" {
        systemGather = `
$systemLogPath = Join-Path -Path $desktop_folder -ChildPath "System.evtx"
$systemLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\Winevt\Logs\System.evtx"
Write-Host "Copying System log to $systemLogPath"
Copy-Item -Path $systemLogSourcePath -Destination $systemLogPath -Force
`
    }

    applicationGather := ""
    if application == "\033[92mtrue \033[0m" {
        applicationGather = `
$applicationLogPath = Join-Path -Path $desktop_folder -ChildPath "Application.evtx"
$applicationLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\Winevt\Logs\Application.evtx"
Write-Host "Copying Application log to $applicationLogPath"
Copy-Item -Path $applicationLogSourcePath -Destination $applicationLogPath -Force
`
    }

    scriptContent := fmt.Sprintf(`
# Get the hostname of the machine
$hostname = $env:COMPUTERNAME

# Define the destination folder path on the desktop
$desktop_folder = Join-Path -Path $env:USERPROFILE -ChildPath "Desktop\$hostname"

# Create the folder on the desktop if it doesn't exist
if (-not (Test-Path -Path $desktop_folder -PathType Container)) {
    New-Item -Path $desktop_folder -ItemType Directory
}

%s
%s
%s
%s

$cauldron_folder = "%s\\Cauldron"

# Copy the entire folder to the Cauldron directory
Write-Host "Copying logs folder to $cauldron_folder"
Copy-Item -Path $desktop_folder -Destination $cauldron_folder -Recurse -Force

Write-Host "Casted logs to Cauldron completed successfully!"

Write-Host "Deleting Desktop folder: $desktop_folder"
Remove-Item -Path $desktop_folder -Recurse -Force
`, securityGather, sysmonGather, systemGather, applicationGather, hostname)

    scriptPath := "C:\\Cauldron\\castingSpell.ps1"
    scriptFile, err := os.Create(scriptPath)
    if err != nil {
        fmt.Printf("Error creating 'castingSpell.ps1': %s\n", err)
        return
    }
    defer scriptFile.Close()

    _, err = scriptFile.WriteString(scriptContent)
    if err != nil {
        fmt.Printf("Error writing to 'castingSpell.ps1': %s\n", err)
        return
    }

    fmt.Printf("Created 'castingSpell.ps1' in the Cauldron directory.\n")
}

func statusMenu() {
    sysmon := "\033[91mfalse\033[0m"
    security := "\033[91mfalse\033[0m"
    system := "\033[91mfalse\033[0m"
    application := "\033[91mfalse\033[0m"
    status := "\033[91mIdle\033[0m"

    for {
        clearScreen()
        banner()
        fmt.Printf("          +-----------Main Menu-----------+\n")
        fmt.Printf("          |   MERlin V.1   |   her3tic    |\n")
        fmt.Printf("          |   Cauldron     |   %sSharing%s    |\n", GREEN, DEFAULT)
        fmt.Printf("          |   Spells       |   %s       |\n", status)
        fmt.Printf("          +-------------------------------+\n\n")
        fmt.Println("[1] Set Spell Rules")
        fmt.Println("[2] Generate Cast Spells Script")
        fmt.Println("[3] Generate Delete Logs Script")
        fmt.Println("[4] Exit\n")

        menuOption := getInput("Enter Selection: ")

        switch menuOption {
        case "1":
            sysmon, security, system, application = logsMenu(sysmon, security, system, application)
        case "2":
            launchSpell(sysmon, security, system, application)
            status = "\033[92mLive\033[0m"
        case "3":
            generateClearLogsScript()
        case "4":
            directoryToDelete := "C:\\Cauldron"
            cmd := exec.Command("cmd", "/C", "rd", "/s", "/q", directoryToDelete)
            cmd.Run()
            os.Exit(0)
        default:
            fmt.Println("bruh")
            clearScreen()
        }
    }
}

func clearScreen() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func getInput(prompt string) string {
    fmt.Print(prompt)
    var input string
    fmt.Scanln(&input)
    return input
}

func main() {
    clearScreen()
    createCauldronFolder()
    clearScreen()
    statusMenu()
}
