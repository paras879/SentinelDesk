"use client";

import { useState } from "react";
import { useTheme } from "next-themes";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Copy,
  Terminal,
  Moon,
  Sun,
  Monitor,
  UserPlus,
  Check,
  MessageSquare,
  ShieldCheck,
  Key,
  Mail,
  User as UserIcon,
  Loader2,
} from "lucide-react";
import { registerUser } from "@/services/auth";
import { toast } from "sonner";

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();

  // New User Creation State
  const [userName, setUserName] = useState("");
  const [userEmail, setUserEmail] = useState("");
  const [userPassword, setUserPassword] = useState("");
  const [creatingUser, setCreatingUser] = useState(false);
  const [createdCredentials, setCreatedCredentials] = useState<{
    name: string;
    email: string;
    password: string;
  } | null>(null);
  const [copiedCreds, setCopiedCreds] = useState(false);

  const deploymentCommand = `Invoke-WebRequest -Uri "http://localhost:8080/agent/download" -OutFile "sentinel-agent.exe"
.\\sentinel-agent.exe install
.\\sentinel-agent.exe start`;

  const copyToClipboard = () => {
    navigator.clipboard.writeText(deploymentCommand);
    toast.success("Command copied to clipboard!");
  };

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!userName.trim() || !userEmail.trim() || !userPassword.trim()) {
      toast.error("Please fill in all user fields");
      return;
    }
    if (userPassword.length < 6) {
      toast.error("Password should be at least 6 characters");
      return;
    }

    setCreatingUser(true);
    try {
      await registerUser({
        name: userName.trim(),
        email: userEmail.trim(),
        password: userPassword,
      });

      setCreatedCredentials({
        name: userName.trim(),
        email: userEmail.trim(),
        password: userPassword,
      });

      toast.success(`Account for ${userName} created successfully!`);
      setUserName("");
      setUserEmail("");
      setUserPassword("");
    } catch (err: any) {
      const msg = err?.response?.data?.error || err?.message || "Failed to create user";
      toast.error(msg);
    } finally {
      setCreatingUser(false);
    }
  };

  const handleCopyCredentials = () => {
    if (!createdCredentials) return;
    const text = `SentinelDesk Dashboard Credentials:
Portal URL: https://sentinel-desk-olive.vercel.app/login
Name: ${createdCredentials.name}
Email: ${createdCredentials.email}
Password: ${createdCredentials.password}`;
    navigator.clipboard.writeText(text);
    setCopiedCreds(true);
    toast.success("Credentials copied to clipboard!");
    setTimeout(() => setCopiedCreds(false), 2000);
  };

  const handleSendViaWhatsApp = () => {
    if (!createdCredentials) return;
    const text = encodeURIComponent(
      `Hello ${createdCredentials.name}! Here are your SentinelDesk Dashboard login details:\n\n🌐 Portal: https://sentinel-desk-olive.vercel.app/login\n📧 Email: ${createdCredentials.email}\n🔑 Password: ${createdCredentials.password}\n\nPlease keep your credentials safe.`
    );
    window.open(`https://wa.me/?text=${text}`, "_blank");
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground">Manage your workspace settings, user accounts, and deployment options.</p>
      </div>

      <Tabs defaultValue="general" className="w-full">
        <TabsList className="grid w-full grid-cols-3 md:w-[480px]">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="users">User Accounts</TabsTrigger>
          <TabsTrigger value="deployment">Deployment</TabsTrigger>
        </TabsList>

        {/* TAB 1: APPEARANCE */}
        <TabsContent value="general" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Appearance</CardTitle>
              <CardDescription>
                Customize how SentinelDesk looks on your device.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-4">
                <Button
                  variant={theme === "light" ? "default" : "outline"}
                  onClick={() => setTheme("light")}
                  className="w-full sm:w-auto"
                >
                  <Sun className="mr-2 h-4 w-4" />
                  Light
                </Button>
                <Button
                  variant={theme === "dark" ? "default" : "outline"}
                  onClick={() => setTheme("dark")}
                  className="w-full sm:w-auto"
                >
                  <Moon className="mr-2 h-4 w-4" />
                  Dark
                </Button>
                <Button
                  variant={theme === "system" ? "default" : "outline"}
                  onClick={() => setTheme("system")}
                  className="w-full sm:w-auto"
                >
                  <Monitor className="mr-2 h-4 w-4" />
                  System
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* TAB 2: USER ACCOUNTS CREATION */}
        <TabsContent value="users" className="mt-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
            {/* Create User Form */}
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <div className="p-2 rounded-xl bg-primary/10 text-primary">
                    <UserPlus className="h-5 w-5" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Create New User Account</CardTitle>
                    <CardDescription>
                      Create login credentials for your clients, staff, or team members.
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleCreateUser} className="space-y-4">
                  <div className="space-y-1.5">
                    <Label htmlFor="name" className="text-xs font-semibold">
                      Full Name *
                    </Label>
                    <div className="relative">
                      <UserIcon className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="name"
                        placeholder="e.g. Rahul Sharma or Client Name"
                        value={userName}
                        onChange={(e) => setUserName(e.target.value)}
                        required
                        className="pl-9 text-xs"
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="email" className="text-xs font-semibold">
                      Email Address *
                    </Label>
                    <div className="relative">
                      <Mail className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="email"
                        type="email"
                        placeholder="e.g. rahul@company.com"
                        value={userEmail}
                        onChange={(e) => setUserEmail(e.target.value)}
                        required
                        className="pl-9 text-xs"
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="password" className="text-xs font-semibold">
                      Account Password *
                    </Label>
                    <div className="relative">
                      <Key className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="password"
                        type="text"
                        placeholder="e.g. Rahul@2025"
                        value={userPassword}
                        onChange={(e) => setUserPassword(e.target.value)}
                        required
                        className="pl-9 text-xs font-mono"
                      />
                    </div>
                    <p className="text-[11px] text-muted-foreground">
                      Minimum 6 characters. You can copy and give this password to your user.
                    </p>
                  </div>

                  <Button
                    type="submit"
                    disabled={creatingUser}
                    className="w-full text-xs font-semibold h-9"
                  >
                    {creatingUser ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Creating User in Database...
                      </>
                    ) : (
                      <>
                        <UserPlus className="mr-2 h-4 w-4" />
                        Create User Account
                      </>
                    )}
                  </Button>
                </form>
              </CardContent>
            </Card>

            {/* Generated Credentials / Sharing Helper */}
            <div className="space-y-4">
              {createdCredentials ? (
                <Card className="border-emerald-500/40 bg-emerald-500/5 shadow-md">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2 text-emerald-600 dark:text-emerald-400 font-bold text-sm">
                        <ShieldCheck className="h-5 w-5" />
                        <span>Account Created Successfully!</span>
                      </div>
                      <Badge className="bg-emerald-600 text-white text-[10px]">Ready to Share</Badge>
                    </div>
                    <CardDescription className="text-xs">
                      Send these login credentials directly to your user via WhatsApp or Email:
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3 pt-0">
                    <div className="p-3 rounded-xl bg-background border font-mono text-xs space-y-1.5">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Portal:</span>
                        <span className="font-semibold text-primary truncate max-w-[220px]">
                          https://sentinel-desk-olive.vercel.app/login
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Name:</span>
                        <span className="font-bold">{createdCredentials.name}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Email:</span>
                        <span className="font-bold">{createdCredentials.email}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Password:</span>
                        <span className="font-bold text-emerald-600 dark:text-emerald-400">
                          {createdCredentials.password}
                        </span>
                      </div>
                    </div>

                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={handleCopyCredentials}
                        className="flex-1 text-xs h-8.5 gap-1.5 font-medium"
                      >
                        {copiedCreds ? (
                          <Check className="h-3.5 w-3.5 text-emerald-500" />
                        ) : (
                          <Copy className="h-3.5 w-3.5" />
                        )}
                        <span>{copiedCreds ? "Copied!" : "Copy Details"}</span>
                      </Button>
                      <Button
                        size="sm"
                        onClick={handleSendViaWhatsApp}
                        className="flex-1 text-xs h-8.5 gap-1.5 font-bold bg-emerald-600 hover:bg-emerald-700 text-white"
                      >
                        <MessageSquare className="h-3.5 w-3.5" />
                        <span>Send WhatsApp</span>
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ) : (
                <Card className="border-dashed bg-muted/20">
                  <CardHeader>
                    <CardTitle className="text-sm font-semibold">How User Management Works</CardTitle>
                    <CardDescription className="text-xs space-y-2 pt-1">
                      <p>
                        • <strong>Security</strong>: Public sign up is disabled so strangers cannot access your dashboard.
                      </p>
                      <p>
                        • <strong>Full Control</strong>: As the master admin, you can create user accounts for your team or clients directly from here.
                      </p>
                      <p>
                        • <strong>Instant Access</strong>: Users can immediately log in at your portal with their assigned email & password.
                      </p>
                    </CardDescription>
                  </CardHeader>
                </Card>
              )}
            </div>
          </div>
        </TabsContent>

        {/* TAB 3: AGENT DEPLOYMENT */}
        <TabsContent value="deployment" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Agent Deployment</CardTitle>
              <CardDescription>
                Download and install the SentinelDesk agent on your Windows endpoints.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  PowerShell Installation Script
                </label>
                <div className="relative">
                  <div className="absolute left-3 top-3 text-muted-foreground">
                    <Terminal className="h-4 w-4" />
                  </div>
                  <pre className="min-h-[100px] w-full rounded-md border border-input bg-muted/50 px-10 py-3 text-sm font-mono ring-offset-background">
                    {deploymentCommand}
                  </pre>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="absolute right-2 top-2 h-8 w-8"
                    onClick={copyToClipboard}
                    title="Copy to clipboard"
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  Run this command in an elevated PowerShell prompt (Run as Administrator).
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}

