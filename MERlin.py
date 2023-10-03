import os
import socket

RED = "\033[91m" 
BLUE = "\033[94m"   
GREEN = "\033[92m"
DEFAULT = "\033[0m"

def createCauldronFolder():
    folder_name = "Cauldron"
    folder_path = os.path.join("C:\\", folder_name)
    
    try:
        if not os.path.exists(folder_path):
            os.mkdir(folder_path)
            print(f"Folder '{folder_name}' created at '{folder_path}'.")
        else:
            print(f"Folder '{folder_name}' already exists.")
    except Exception as e:
        print(f"Error creating folder: {e}")
    
    os.system("icacls C:\Cauldron /grant Everyone:(OI)(CI)F")
    os.system("net share Cauldron=C:\Cauldron")

def banner():
    ascii_art = """\033[34m
    ███╗░░░███╗███████╗██████╗░██╗░░░░░██╗███╗░░██╗
    ████╗░████║██╔════╝██╔══██╗██║░░░░░██║████╗░██║
    ██╔████╔██║█████╗░░██████╔╝██║░░░░░██║██╔██╗██║
    ██║╚██╔╝██║██╔══╝░░██╔══██╗██║░░░░░██║██║╚████║
    ██║░╚═╝░██║███████╗██║░░██║███████╗██║██║░╚███║
    ╚═╝░░░░░╚═╝╚══════╝╚═╝░░╚═╝╚══════╝╚═╝╚═╝░░╚══╝\033[0m
    """

    print(ascii_art)

def logsMenu(sysmon_logs, security_logs, system_logs, application_logs):
    while 1:
        os.system("cls")
        banner()
        print(f"          +-----------SpellBook-----------+")
        print(f"          |  [1] Sysmon      |   " + sysmon_logs + "    |")
        print(f"          |  [2] Security    |   " + security_logs + "    |")
        print(f"          |  [3] System      |   " + system_logs + "    |")
        print(f"          |  [4] Application |   " + application_logs + "    |")
        print(f"          |  [B] Back        |  Main Menu |")
        print(f"          +-------------------------------+")

        toggle = input("\n         Select Spell: ")

        if toggle == "1":
            if sysmon_logs == "\033[92mtrue \033[0m":
                sysmon_logs ="\033[91mfalse\033[0m"
            else:
                sysmon_logs = "\033[92mtrue \033[0m"
        elif toggle == "2":
            if security_logs == "\033[92mtrue \033[0m":
                security_logs ="\033[91mfalse\033[0m"
            else:
                security_logs = "\033[92mtrue \033[0m"
        elif toggle == "3":
            if system_logs == "\033[92mtrue \033[0m":
                system_logs ="\033[91mfalse\033[0m"
            else:
                system_logs = "\033[92mtrue \033[0m"
        elif toggle == "4":
            if application_logs == "\033[92mtrue \033[0m":
                application_logs ="\033[91mfalse\033[0m"
            else:
                application_logs = "\033[92mtrue \033[0m"
        elif toggle == "b" or toggle == "B":
            return sysmon_logs, security_logs, system_logs, application_logs
        else:
            print("bruh")

def generate_clear_logs_script():
    script_content = """
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
"""

    target_directory = "C:\\Cauldron"
    os.makedirs(target_directory, exist_ok=True)

    with open(os.path.join(target_directory, "clearLogs.ps1"), "w") as file:
        file.write(script_content)

def launchSpell(sysmon, security, system, application):
    hostname = socket.gethostname()
    if not hostname.startswith("\\\\"):
        hostname = f"\\\\{hostname}"
    if security == "\033[92mtrue \033[0m":
        securityGather = '''
$securityLogPath = Join-Path -Path $desktop_folder -ChildPath "Security.evtx"
Write-Host "Copying Security log to $securityLogPath"
Copy-Item -Path "C:\Windows\System32\winevt\Logs\Security.evtx" -Destination $securityLogPath -Force
'''
    else:
        securityGather = " "
    
    if sysmon == "\033[92mtrue \033[0m":
        sysmonGather = '''
# Copy the Sysmon log to the destination folder
$sysmonLogPath = Join-Path -Path $filePath -ChildPath "Sysmon.evtx"
$sysmonLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\\Winevt\\Logs\\Microsoft-Windows-Sysmon%4Operational.evtx"

Write-Host "Copying Sysmon log to $sysmonLogPath"
Copy-Item -Path $sysmonLogSourcePath -Destination $sysmonLogPath -Force
        '''
    else:
        sysmonGather = " "

    if system == "\033[92mtrue \033[0m":
        systemGather = '''
$systemLogPath = Join-Path -Path $desktop_folder -ChildPath "System.evtx"
$systemLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\Winevt\Logs\System.evtx"
Write-Host "Copying System log to $systemLogPath"
Copy-Item -Path $systemLogSourcePath -Destination $systemLogPath -Force
        '''
    else:
        systemGather = ""

    if application == "\033[92mtrue \033[0m":
        applicationGather = '''
$applicationLogPath = Join-Path -Path $desktop_folder -ChildPath "Application.evtx"
$applicationLogSourcePath = Join-Path -Path $env:SystemRoot -ChildPath "System32\Winevt\Logs\Application.evtx"
Write-Host "Copying Application log to $applicationLogPath"
Copy-Item -Path $applicationLogSourcePath -Destination $applicationLogPath -Force
        '''
    else:
        applicationGather = ""
#================================================================================
    script_content = '''
    
# Get the hostname of the machine
$hostname = $env:COMPUTERNAME

# Define the destination folder path on the desktop
$desktop_folder = Join-Path -Path $env:USERPROFILE -ChildPath "Desktop\$hostname"

# Create the folder on the desktop if it doesn't exist
if (-not (Test-Path -Path $desktop_folder -PathType Container)) {
    New-Item -Path $desktop_folder -ItemType Directory
}

'''+ securityGather +'''
'''+ sysmonGather +'''
'''+ systemGather +'''
'''+ applicationGather +'''

$cauldron_folder = "'''+ hostname +'''\Cauldron"

# Copy the entire folder to the Cauldron directory
Write-Host "Copying logs folder to $cauldron_folder"
Copy-Item -Path $desktop_folder -Destination $cauldron_folder -Recurse -Force

Write-Host "Casted logs to Cauldron completed successfully!"

Write-Host "Deleting Desktop folder: $desktop_folder"
Remove-Item -Path $desktop_folder -Recurse -Force'''

#=================================================================================

    script_path = os.path.join("C:\\Cauldron", "castingSpell.ps1")

    try:
        with open(script_path, "w") as script_file:
            script_file.write(script_content)
        print(f"Created 'castingSpell.ps1' in the Cauldron directory.")
    except Exception as e:
        print(f"Error creating 'castingSpell.ps1': {e}")

def statusMenu():
    sysmon = "\033[91mfalse\033[0m"
    security = "\033[91mfalse\033[0m"
    system = "\033[91mfalse\033[0m"
    application = "\033[91mfalse\033[0m"
    status = "\033[91mIdle\033[0m"
    ballStatus = "\033[91mIdle\033[0m"

    while 1:
        os.system("cls")
        banner()
        print(f"          +-----------Main Menu-----------+")
        print(f"          |   MERlin V.1   |   her3tic    |")
        print(f"          |   Cauldron     |   {GREEN}Sharing{DEFAULT}    |")
        print("          |   Spells       |   "+ status +"       |")
        print(f"          +-------------------------------+\n")
        print("          [1] Set Spell Rules")
        print("          [2] Generate Cast Spells Script")
        print("          [3] Generate Delete Logs Script")
        print("          [4] Exit\n")


        menuOption = input("         Enter Selection: ")

        if menuOption == '1':
            sysmon, security, system, application = logsMenu(sysmon, security, system, application)
        elif menuOption == '2':
            launchSpell(sysmon, security, system, application)
            status = "\033[92mLive\033[0m"
        elif menuOption == '3':
            generate_clear_logs_script()
        elif menuOption == '4':
            directory_to_delete = "C:\\Cauldron"
            os.system(f'rd /s /q "{directory_to_delete}"')
            exit(0)
        else:
            print("bruh")
            os.system("cls")


if __name__ == "__main__":
    os.system("cls")
    createCauldronFolder()
    os.system("cls")
    statusMenu()
