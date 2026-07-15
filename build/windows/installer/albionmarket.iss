; Albion Market Data Client - Inno Setup installer.
; Bundles WinPcap 4.1.3 (redistributable) and the Microsoft-signed WebView2
; bootstrapper, so nothing is downloaded/executed from the internet at install
; time (avoids the PowerShell-download dropper pattern that trips AV heuristics).
;
; Build:  ISCC /DMyAppVersion=0.9.5 albionmarket.iss
; Paths are relative to this .iss file (build/windows/installer/).

#define MyAppName "Albion Market Data Client"
#define MyAppPublisher "Albion Market"
#define MyAppURL "https://albionmarket.gg"
#define MyAppExeName "AlbionMarketDataClient.exe"
#ifndef MyAppVersion
  #define MyAppVersion "0.0.0"
#endif

[Setup]
AppId={{A6F3B2E1-5C4D-4A9E-9B7F-albionmarketgg}}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
UninstallDisplayIcon={app}\{#MyAppExeName}
OutputDir=..\..\bin
OutputBaseFilename=AlbionMarketDataClient-Setup
SetupIconFile=..\icon.ico
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

[Files]
Source: "..\..\bin\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\thirdparty\WinPcap_4_1_3.exe"; DestDir: "{tmp}"; Flags: deleteafterinstall
Source: "..\thirdparty\MicrosoftEdgeWebview2Setup.exe"; DestDir: "{tmp}"; Flags: deleteafterinstall

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Run]
; WebView2 runtime (Microsoft-signed) - only if the runtime is missing.
Filename: "{tmp}\MicrosoftEdgeWebview2Setup.exe"; Parameters: "/silent /install"; StatusMsg: "Installing WebView2 runtime..."; Flags: waituntilterminated; Check: not IsWebView2Installed
; WinPcap (packet-capture driver) - only if no compatible capture driver exists.
Filename: "{tmp}\WinPcap_4_1_3.exe"; Parameters: "/S"; StatusMsg: "Installing WinPcap (packet capture)..."; Flags: waituntilterminated; Check: not IsPcapInstalled
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#MyAppName}}"; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "{sys}\taskkill.exe"; Parameters: "/F /IM {#MyAppExeName}"; Flags: runhidden; RunOnceId: "KillApp"

[Code]
// Close the running app before an upgrade so its exe isn't locked.
function PrepareToInstall(var NeedsRestart: Boolean): String;
var rc: Integer;
begin
  Exec(ExpandConstant('{sys}\taskkill.exe'), '/F /IM {#MyAppExeName}', '', SW_HIDE, ewWaitUntilTerminated, rc);
  Result := '';
end;

// True if any WinPcap-API-compatible capture driver is already present.
function IsPcapInstalled: Boolean;
begin
  Result := RegKeyExists(HKLM, 'SOFTWARE\WOW6432Node\WinPcap') or
            RegKeyExists(HKLM, 'SOFTWARE\WOW6432Node\Npcap') or
            RegKeyExists(HKLM, 'SYSTEM\CurrentControlSet\Services\npcap') or
            RegKeyExists(HKLM, 'SYSTEM\CurrentControlSet\Services\npf');
end;

// True if the Edge WebView2 Evergreen runtime is installed (machine or user).
function IsWebView2Installed: Boolean;
var v: String;
begin
  Result :=
    RegQueryStringValue(HKLM, 'SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}', 'pv', v) or
    RegQueryStringValue(HKCU, 'SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}', 'pv', v);
  if Result and ((v = '') or (v = '0.0.0.0')) then
    Result := False;
end;

// Remove the per-user autostart entry on uninstall.
procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RegDeleteValue(HKCU, 'Software\Microsoft\Windows\CurrentVersion\Run', 'AlbionMarketDataClient');
end;
