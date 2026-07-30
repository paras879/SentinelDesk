"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/hooks/use-auth";
import { updateProfile } from "@/services/auth";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

export default function ProfilePage() {
  const { user } = useAuth();
  
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [isEditingProfile, setIsEditingProfile] = useState(false);
  const [isSavingProfile, setIsSavingProfile] = useState(false);

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isSavingPassword, setIsSavingPassword] = useState(false);

  useEffect(() => {
    if (user) {
      setName(user.name);
      setEmail(user.email);
    }
  }, [user]);

  const handleUpdateProfile = async () => {
    try {
      setIsSavingProfile(true);
      await updateProfile({ name, email });
      alert("Profile updated successfully! Refresh to see changes globally.");
      setIsEditingProfile(false);
    } catch (error) {
      alert("Failed to update profile.");
    } finally {
      setIsSavingProfile(false);
    }
  };

  const handleUpdatePassword = async () => {
    if (password !== confirmPassword) {
      alert("Passwords do not match!");
      return;
    }
    if (!password) {
      alert("Please enter a new password.");
      return;
    }

    try {
      setIsSavingPassword(true);
      await updateProfile({ password });
      alert("Password updated successfully!");
      setPassword("");
      setConfirmPassword("");
    } catch (error: any) {
      alert("Failed to update password: " + (error?.response?.data?.error || error.message));
    } finally {
      setIsSavingPassword(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold tracking-tight">Profile</h1>
        <p className="text-sm text-muted-foreground">Manage your personal information and security preferences.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
            <CardDescription>
              Your personal profile details.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center space-x-4">
              <Avatar className="h-20 w-20">
                <AvatarFallback className="text-2xl">
                  {name?.charAt(0)?.toUpperCase() || "A"}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-1">
                <h3 className="font-medium text-lg leading-none">{name || "Admin User"}</h3>
                <p className="text-sm text-muted-foreground">{email || "admin@sentineldesk.com"}</p>
                <div className="pt-2">
                  <Badge variant="secondary">Administrator</Badge>
                </div>
              </div>
            </div>
            
            <div className="space-y-4 pt-4 border-t">
              <div className="space-y-2">
                <Label htmlFor="name">Full Name</Label>
                <Input 
                  id="name" 
                  value={name} 
                  onChange={(e) => setName(e.target.value)}
                  readOnly={!isEditingProfile} 
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input 
                  id="email" 
                  type="email" 
                  value={email} 
                  onChange={(e) => setEmail(e.target.value)}
                  readOnly={!isEditingProfile} 
                />
              </div>
            </div>
          </CardContent>
          <CardFooter className="flex justify-end gap-2">
            {isEditingProfile ? (
              <>
                <Button variant="outline" onClick={() => {
                  setIsEditingProfile(false);
                  setName(user?.name || "");
                  setEmail(user?.email || "");
                }}>Cancel</Button>
                <Button onClick={handleUpdateProfile} disabled={isSavingProfile}>
                  {isSavingProfile ? "Saving..." : "Save Changes"}
                </Button>
              </>
            ) : (
              <Button variant="outline" onClick={() => setIsEditingProfile(true)}>
                Edit Profile
              </Button>
            )}
          </CardFooter>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Change Password</CardTitle>
            <CardDescription>
              Update your password to keep your account secure.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="new-password">New Password</Label>
              <Input 
                id="new-password" 
                type="password" 
                placeholder="••••••••" 
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirm-password">Confirm Password</Label>
              <Input 
                id="confirm-password" 
                type="password" 
                placeholder="••••••••" 
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>
          </CardContent>
          <CardFooter className="flex justify-end">
            <Button onClick={handleUpdatePassword} disabled={isSavingPassword || !password || password !== confirmPassword}>
              {isSavingPassword ? "Updating..." : "Update Password"}
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
