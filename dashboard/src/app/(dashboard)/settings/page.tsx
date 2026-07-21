"use client";

import { useTheme } from "next-themes";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Copy, Terminal, Moon, Sun, Monitor } from "lucide-react";

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();
  
  const deploymentCommand = `Invoke-WebRequest -Uri "http://localhost:8080/agent/download" -OutFile "sentinel-agent.exe"
.\\sentinel-agent.exe install
.\\sentinel-agent.exe start`;

  const copyToClipboard = () => {
    navigator.clipboard.writeText(deploymentCommand);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground">Manage your workspace settings and deployment options.</p>
      </div>

      <Tabs defaultValue="general" className="w-full">
        <TabsList className="grid w-full grid-cols-2 md:w-[400px]">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="deployment">Deployment</TabsTrigger>
        </TabsList>
        
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
                  variant={theme === 'light' ? 'default' : 'outline'} 
                  onClick={() => setTheme('light')}
                  className="w-full sm:w-auto"
                >
                  <Sun className="mr-2 h-4 w-4" />
                  Light
                </Button>
                <Button 
                  variant={theme === 'dark' ? 'default' : 'outline'} 
                  onClick={() => setTheme('dark')}
                  className="w-full sm:w-auto"
                >
                  <Moon className="mr-2 h-4 w-4" />
                  Dark
                </Button>
                <Button 
                  variant={theme === 'system' ? 'default' : 'outline'} 
                  onClick={() => setTheme('system')}
                  className="w-full sm:w-auto"
                >
                  <Monitor className="mr-2 h-4 w-4" />
                  System
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
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
