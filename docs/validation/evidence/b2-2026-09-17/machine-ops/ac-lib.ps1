# B2 scratch library: unelevated AppContainer launcher (P/Invoke) shared by the B2 batteries.
# Delete at slice closeout.
$ErrorActionPreference = 'Continue'

if (-not ('AcLauncher2' -as [type])) {
    $csharp = @'
using System;
using System.Runtime.InteropServices;

public static class AcLauncher2
{
    [StructLayout(LayoutKind.Sequential)]
    public struct STARTUPINFO
    {
        public int cb; public IntPtr lpReserved; public IntPtr lpDesktop; public IntPtr lpTitle;
        public int dwX; public int dwY; public int dwXSize; public int dwYSize; public int dwXCountChars; public int dwYCountChars;
        public int dwFillAttribute; public int dwFlags; public short wShowWindow; public short cbReserved2;
        public IntPtr lpReserved2; public IntPtr hStdInput; public IntPtr hStdOutput; public IntPtr hStdError;
    }

    [StructLayout(LayoutKind.Sequential)]
    public struct STARTUPINFOEX { public STARTUPINFO StartupInfo; public IntPtr lpAttributeList; }

    [StructLayout(LayoutKind.Sequential)]
    public struct PROCESS_INFORMATION { public IntPtr hProcess; public IntPtr hThread; public int dwProcessId; public int dwThreadId; }

    [StructLayout(LayoutKind.Sequential)]
    public struct SID_AND_ATTRIBUTES { public IntPtr Sid; public uint Attributes; }

    [StructLayout(LayoutKind.Sequential)]
    public struct SECURITY_CAPABILITIES { public IntPtr AppContainerSid; public IntPtr Capabilities; public uint CapabilityCount; public uint Reserved; }

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int CreateAppContainerProfile(string name, string displayName, string description, IntPtr capabilities, uint capabilityCount, out IntPtr sid);

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int DeriveAppContainerSidFromAppContainerName(string name, out IntPtr sid);

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int DeleteAppContainerProfile(string name);

    [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool ConvertSidToStringSid(IntPtr sid, out IntPtr stringSid);

    [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool ConvertStringSidToSid(string sddlSid, out IntPtr sid);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool InitializeProcThreadAttributeList(IntPtr list, int count, int flags, ref IntPtr size);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool UpdateProcThreadAttribute(IntPtr list, uint flags, IntPtr attribute, IntPtr value, IntPtr size, IntPtr previous, IntPtr returned);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern void DeleteProcThreadAttributeList(IntPtr list);

    [DllImport("kernel32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool CreateProcess(string app, string commandLine, IntPtr pa, IntPtr ta, bool inherit, uint flags, IntPtr env, string cwd, ref STARTUPINFOEX si, out PROCESS_INFORMATION pi);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern uint WaitForSingleObject(IntPtr handle, uint ms);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool GetExitCodeProcess(IntPtr process, out uint code);

    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool CloseHandle(IntPtr handle);

    private const uint ATTR_SECURITY_CAPABILITIES = 0x00020009;
    private const uint EXTENDED_STARTUPINFO_PRESENT = 0x00080000;
    private const uint CREATE_NO_WINDOW = 0x08000000;
    private const uint CREATE_UNICODE_ENVIRONMENT = 0x00000400;
    private const uint SE_GROUP_ENABLED = 0x00000004;

    public static string SidToString(IntPtr sid)
    {
        IntPtr s;
        if (!ConvertSidToStringSid(sid, out s)) { return ""; }
        string r = Marshal.PtrToStringUni(s);
        return r == null ? "" : r;
    }

    // Builds a UTF-16 environment block (sorted, double-NUL terminated) from "K=V" pairs so secrets
    // travel through the environment and never through a command line.
    public static IntPtr BuildEnvironmentBlock(string[] pairs)
    {
        string[] sorted = (string[])pairs.Clone();
        Array.Sort(sorted, StringComparer.OrdinalIgnoreCase);
        System.Text.StringBuilder builder = new System.Text.StringBuilder();
        foreach (string pair in sorted)
        {
            builder.Append(pair);
            builder.Append('\0');
        }
        builder.Append('\0');
        return Marshal.StringToHGlobalUni(builder.ToString());
    }

    [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool CredRead(string target, uint type, uint flags, out IntPtr credential);

    [DllImport("advapi32.dll", SetLastError = true)]
    public static extern void CredFree(IntPtr buffer);

    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    private struct CREDENTIAL
    {
        public uint Flags;
        public uint Type;
        public IntPtr TargetName;
        public IntPtr Comment;
        public System.Runtime.InteropServices.ComTypes.FILETIME LastWritten;
        public uint CredentialBlobSize;
        public IntPtr CredentialBlob;
        public uint Persist;
        public uint AttributeCount;
        public IntPtr Attributes;
        public IntPtr TargetAlias;
        public IntPtr UserName;
    }

    // Reads a generic credential secret into a string. The caller owns the secret and must never log it.
    public static string ReadGenericCredentialSecret(string target, out int blobSize)
    {
        blobSize = 0;
        IntPtr handle;
        if (!CredRead(target, 1 /*CRED_TYPE_GENERIC*/, 0, out handle)) { return null; }
        try
        {
            CREDENTIAL credential = (CREDENTIAL)Marshal.PtrToStructure(handle, typeof(CREDENTIAL));
            blobSize = (int)credential.CredentialBlobSize;
            if (credential.CredentialBlob == IntPtr.Zero || credential.CredentialBlobSize == 0) { return ""; }
            byte[] bytes = new byte[credential.CredentialBlobSize];
            Marshal.Copy(credential.CredentialBlob, bytes, 0, bytes.Length);
            // credential blobs are UTF-16LE for generic credentials written by tooling
            string text = System.Text.Encoding.Unicode.GetString(bytes).TrimEnd('\0');
            return text;
        }
        finally
        {
            CredFree(handle);
        }
    }

    public static int LaunchWithEnv(string containerSid, string[] capabilitySids, string commandLine, string cwd, string[] envPairs, out int lastError, out int createdPid)
    {
        lastError = 0; createdPid = 0;
        IntPtr container;
        if (!ConvertStringSidToSid(containerSid, out container)) { lastError = Marshal.GetLastWin32Error(); return -1; }

        IntPtr capArray = IntPtr.Zero;
        int capCount = capabilitySids == null ? 0 : capabilitySids.Length;
        if (capCount > 0)
        {
            int elem = Marshal.SizeOf(typeof(SID_AND_ATTRIBUTES));
            capArray = Marshal.AllocHGlobal(elem * capCount);
            for (int i = 0; i < capCount; i++)
            {
                IntPtr capSid;
                if (!ConvertStringSidToSid(capabilitySids[i], out capSid)) { lastError = Marshal.GetLastWin32Error(); return -2; }
                SID_AND_ATTRIBUTES e = new SID_AND_ATTRIBUTES();
                e.Sid = capSid; e.Attributes = SE_GROUP_ENABLED;
                Marshal.StructureToPtr(e, (IntPtr)(capArray.ToInt64() + i * elem), false);
            }
        }

        SECURITY_CAPABILITIES caps = new SECURITY_CAPABILITIES();
        caps.AppContainerSid = container; caps.Capabilities = capArray; caps.CapabilityCount = (uint)capCount; caps.Reserved = 0;
        int capSize = Marshal.SizeOf(typeof(SECURITY_CAPABILITIES));
        IntPtr capsPtr = Marshal.AllocHGlobal(capSize);
        Marshal.StructureToPtr(caps, capsPtr, false);

        IntPtr size = IntPtr.Zero;
        InitializeProcThreadAttributeList(IntPtr.Zero, 1, 0, ref size);
        IntPtr list = Marshal.AllocHGlobal(size);
        if (!InitializeProcThreadAttributeList(list, 1, 0, ref size)) { lastError = Marshal.GetLastWin32Error(); return -3; }
        if (!UpdateProcThreadAttribute(list, 0, (IntPtr)ATTR_SECURITY_CAPABILITIES, capsPtr, (IntPtr)capSize, IntPtr.Zero, IntPtr.Zero))
        {
            lastError = Marshal.GetLastWin32Error();
            DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); return -4;
        }

        IntPtr envBlock = IntPtr.Zero;
        if (envPairs != null && envPairs.Length > 0) { envBlock = BuildEnvironmentBlock(envPairs); }

        STARTUPINFOEX si = new STARTUPINFOEX();
        si.StartupInfo.cb = Marshal.SizeOf(typeof(STARTUPINFOEX));
        si.lpAttributeList = list;

        PROCESS_INFORMATION pi;
        uint flags = CREATE_NO_WINDOW | EXTENDED_STARTUPINFO_PRESENT;
        if (envBlock != IntPtr.Zero) { flags |= CREATE_UNICODE_ENVIRONMENT; }
        if (!CreateProcess(null, commandLine, IntPtr.Zero, IntPtr.Zero, false, flags, envBlock, cwd, ref si, out pi))
        {
            lastError = Marshal.GetLastWin32Error();
            DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); Marshal.FreeHGlobal(capsPtr);
            if (envBlock != IntPtr.Zero) { Marshal.FreeHGlobal(envBlock); }
            return -5;
        }

        createdPid = pi.dwProcessId;
        WaitForSingleObject(pi.hProcess, 0xFFFFFFFF);
        uint code;
        GetExitCodeProcess(pi.hProcess, out code);
        CloseHandle(pi.hThread); CloseHandle(pi.hProcess);
        DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); Marshal.FreeHGlobal(capsPtr);
        if (capArray != IntPtr.Zero) { Marshal.FreeHGlobal(capArray); }
        if (envBlock != IntPtr.Zero) { Marshal.FreeHGlobal(envBlock); }
        return (int)code;
    }

    public static int Launch(string containerSid, string[] capabilitySids, string commandLine, string cwd, out int lastError, out int createdPid)
    {
        lastError = 0; createdPid = 0;
        IntPtr container;
        if (!ConvertStringSidToSid(containerSid, out container)) { lastError = Marshal.GetLastWin32Error(); return -1; }

        IntPtr capArray = IntPtr.Zero;
        int capCount = capabilitySids == null ? 0 : capabilitySids.Length;
        if (capCount > 0)
        {
            int elem = Marshal.SizeOf(typeof(SID_AND_ATTRIBUTES));
            capArray = Marshal.AllocHGlobal(elem * capCount);
            for (int i = 0; i < capCount; i++)
            {
                IntPtr capSid;
                if (!ConvertStringSidToSid(capabilitySids[i], out capSid)) { lastError = Marshal.GetLastWin32Error(); return -2; }
                SID_AND_ATTRIBUTES e = new SID_AND_ATTRIBUTES();
                e.Sid = capSid; e.Attributes = SE_GROUP_ENABLED;
                Marshal.StructureToPtr(e, (IntPtr)(capArray.ToInt64() + i * elem), false);
            }
        }

        SECURITY_CAPABILITIES caps = new SECURITY_CAPABILITIES();
        caps.AppContainerSid = container; caps.Capabilities = capArray; caps.CapabilityCount = (uint)capCount; caps.Reserved = 0;
        int capSize = Marshal.SizeOf(typeof(SECURITY_CAPABILITIES));
        IntPtr capsPtr = Marshal.AllocHGlobal(capSize);
        Marshal.StructureToPtr(caps, capsPtr, false);

        IntPtr size = IntPtr.Zero;
        InitializeProcThreadAttributeList(IntPtr.Zero, 1, 0, ref size);
        IntPtr list = Marshal.AllocHGlobal(size);
        if (!InitializeProcThreadAttributeList(list, 1, 0, ref size)) { lastError = Marshal.GetLastWin32Error(); return -3; }
        if (!UpdateProcThreadAttribute(list, 0, (IntPtr)ATTR_SECURITY_CAPABILITIES, capsPtr, (IntPtr)capSize, IntPtr.Zero, IntPtr.Zero))
        {
            lastError = Marshal.GetLastWin32Error();
            DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); return -4;
        }

        STARTUPINFOEX si = new STARTUPINFOEX();
        si.StartupInfo.cb = Marshal.SizeOf(typeof(STARTUPINFOEX));
        si.lpAttributeList = list;

        PROCESS_INFORMATION pi;
        if (!CreateProcess(null, commandLine, IntPtr.Zero, IntPtr.Zero, false, CREATE_NO_WINDOW | EXTENDED_STARTUPINFO_PRESENT, IntPtr.Zero, cwd, ref si, out pi))
        {
            lastError = Marshal.GetLastWin32Error();
            DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); Marshal.FreeHGlobal(capsPtr); return -5;
        }

        createdPid = pi.dwProcessId;
        WaitForSingleObject(pi.hProcess, 0xFFFFFFFF);
        uint code;
        GetExitCodeProcess(pi.hProcess, out code);
        CloseHandle(pi.hThread); CloseHandle(pi.hProcess);
        DeleteProcThreadAttributeList(list); Marshal.FreeHGlobal(list); Marshal.FreeHGlobal(capsPtr);
        if (capArray != IntPtr.Zero) { Marshal.FreeHGlobal(capArray); }
        return (int)code;
    }
}
'@
    Add-Type -TypeDefinition $csharp -Language CSharp
}

# NOTE: this library intentionally no longer creates the AppContainer profile. `apply-b2-loopback.ps1` is
# the single entry point allowed to create machine state; every harness resolves the SID derive-only via
# `Get-B2ContainerSidExisting` (b2-loopback-lib.ps1) and fails with a clear message when apply has not run.
function Get-B2ContainerSid([string]$Name) {
    $sidPtr = [IntPtr]::Zero
    $hr = [AcLauncher2]::DeriveAppContainerSidFromAppContainerName($Name, [ref]$sidPtr)
    if ($hr -ne 0 -or $sidPtr -eq [IntPtr]::Zero) {
        throw ("AppContainer profile '{0}' does not exist (derive failed 0x{1:X8}); run apply-b2-loopback.ps1 first" -f $Name, $hr)
    }
    return [AcLauncher2]::SidToString($sidPtr)
}
