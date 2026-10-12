"use client";

import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Check, Sparkles, ShieldCheck, Zap, MessageSquare, Mail, Server } from "lucide-react";

interface UpgradeModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  currentDevices?: number;
  maxDevices?: number;
}

export function UpgradeModal({
  open,
  onOpenChange,
  currentDevices = 1,
  maxDevices = 3,
}: UpgradeModalProps) {
  const handleWhatsAppContact = (plan: string) => {
    const text = encodeURIComponent(
      `Hello! I want to upgrade my SentinelDesk license to ${plan} Plan (My current usage: ${currentDevices}/${maxDevices} devices). Please share payment and activation details.`
    );
    window.open(`https://wa.me/?text=${text}`, "_blank");
  };

  const handleEmailContact = (plan: string) => {
    const subject = encodeURIComponent(`SentinelDesk License Upgrade - ${plan} Plan`);
    const body = encodeURIComponent(
      `Hello Sales Team,\n\nI want to upgrade my SentinelDesk license to the ${plan} Plan.\nCurrent active devices: ${currentDevices}.\n\nPlease provide payment and setup instructions.\n\nThank you!`
    );
    window.open(`mailto:support@sentineldesk.com?subject=${subject}&body=${body}`, "_blank");
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl p-6 sm:p-8 bg-gradient-to-b from-background to-muted/20 border border-primary/20">
        <DialogHeader className="text-center space-y-2 pb-4">
          <div className="flex items-center justify-center gap-2">
            <span className="p-2 rounded-xl bg-primary/10 text-primary">
              <Zap className="h-6 w-6" />
            </span>
          </div>
          <DialogTitle className="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Upgrade Your <span className="text-primary">SentinelDesk</span> Plan
          </DialogTitle>
          <DialogDescription className="text-sm sm:text-base max-w-xl mx-auto">
            Scale endpoint monitoring across your enterprise. Connect more laptops, enable multi-screen live walls, and unlock enterprise controls.
          </DialogDescription>
        </DialogHeader>

        {/* Pricing Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-2">
          {/* Starter Plan */}
          <div className="rounded-2xl border bg-card p-5 flex flex-col justify-between shadow-sm relative">
            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-bold text-lg">Starter</h3>
                <Badge variant="secondary" className="text-xs">Current Plan</Badge>
              </div>
              <p className="text-xs text-muted-foreground mb-4">
                Basic monitoring for small setups and testing.
              </p>
              <div className="mb-4">
                <span className="text-3xl font-extrabold">Free</span>
                <span className="text-xs text-muted-foreground ml-1">/ forever</span>
              </div>
              <ul className="space-y-2 text-xs text-muted-foreground">
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Up to <strong>3 Devices</strong></span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Live telemetry & hardware stats</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Process & Service manager</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Single-screen Live View</span>
                </li>
              </ul>
            </div>
            <Button variant="outline" className="w-full mt-6 text-xs" disabled>
              Active Plan
            </Button>
          </div>

          {/* Pro Plan */}
          <div className="rounded-2xl border-2 border-primary bg-card/60 p-5 flex flex-col justify-between shadow-lg relative ring-1 ring-primary/20">
            <div className="absolute -top-3 left-1/2 -translate-x-1/2">
              <Badge className="bg-primary text-primary-foreground text-[10px] px-2.5 py-0.5 shadow font-semibold flex items-center gap-1">
                <Sparkles className="h-3 w-3" /> MOST POPULAR
              </Badge>
            </div>
            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-bold text-lg text-primary">Professional</h3>
                <Badge className="bg-primary/10 text-primary border-primary/25 text-xs">25 Devices</Badge>
              </div>
              <p className="text-xs text-muted-foreground mb-4">
                Complete endpoint security for growing offices & teams.
              </p>
              <div className="mb-4">
                <span className="text-3xl font-extrabold">₹1,499</span>
                <span className="text-xs text-muted-foreground ml-1">/ month</span>
              </div>
              <ul className="space-y-2 text-xs">
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Up to <strong>25 Devices</strong></span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Multi-Screen Matrix View (6 Screens)</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Ultra-low latency streaming</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Remote Service Restart & Control</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Priority WhatsApp & Email Support</span>
                </li>
              </ul>
            </div>
            <div className="space-y-2 mt-6">
              <Button
                onClick={() => handleWhatsAppContact("Pro (25 Devices)")}
                className="w-full text-xs font-semibold shadow-md flex items-center justify-center gap-1.5"
              >
                <MessageSquare className="h-3.5 w-3.5" /> Upgrade via WhatsApp
              </Button>
              <Button
                variant="outline"
                onClick={() => handleEmailContact("Pro (25 Devices)")}
                className="w-full text-xs"
              >
                <Mail className="h-3.5 w-3.5 mr-1" /> Upgrade via Email
              </Button>
            </div>
          </div>

          {/* Enterprise Plan */}
          <div className="rounded-2xl border bg-card p-5 flex flex-col justify-between shadow-sm">
            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-bold text-lg">Enterprise</h3>
                <Badge variant="outline" className="text-xs">Custom</Badge>
              </div>
              <p className="text-xs text-muted-foreground mb-4">
                Dedicated infrastructure for large corporate networks.
              </p>
              <div className="mb-4">
                <span className="text-3xl font-extrabold">Custom</span>
              </div>
              <ul className="space-y-2 text-xs text-muted-foreground">
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span><strong>Unlimited Devices</strong></span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Dedicated Private Server & DB</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Custom White-label Branding</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>Full Audit Logs & Retention</span>
                </li>
                <li className="flex items-center gap-2">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  <span>24/7 Dedicated Account Manager</span>
                </li>
              </ul>
            </div>
            <Button
              variant="outline"
              onClick={() => handleWhatsAppContact("Enterprise (Unlimited)")}
              className="w-full mt-6 text-xs flex items-center justify-center gap-1.5"
            >
              <Server className="h-3.5 w-3.5" /> Contact Sales
            </Button>
          </div>
        </div>

        {/* Footer info */}
        <div className="text-center text-xs text-muted-foreground pt-3 border-t flex items-center justify-center gap-2">
          <ShieldCheck className="h-4 w-4 text-green-500" />
          <span>Licenses are activated instantly upon confirmation without server downtime.</span>
        </div>
      </DialogContent>
    </Dialog>
  );
}
