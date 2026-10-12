"use client";

import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Check,
  Sparkles,
  ShieldCheck,
  Zap,
  MessageSquare,
  Mail,
  Server,
  QrCode,
  Copy,
  CheckCircle2,
  ArrowLeft,
  CreditCard,
  Lock,
  ExternalLink,
} from "lucide-react";
import { toast } from "sonner";
import { BILLING_CONFIG } from "@/config/billing";

interface UpgradeModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  currentDevices?: number;
  maxDevices?: number;
}

type BillingCycle = "monthly" | "yearly";
type ModalView = "plans" | "checkout" | "success";

export function UpgradeModal({
  open,
  onOpenChange,
  currentDevices = 1,
  maxDevices = 3,
}: UpgradeModalProps) {
  const [billingCycle, setBillingCycle] = useState<BillingCycle>("monthly");
  const [view, setView] = useState<ModalView>("plans");
  const [selectedPlan, setSelectedPlan] = useState<"pro" | "enterprise">("pro");
  const [utrNumber, setUtrNumber] = useState("");
  const [customerContact, setCustomerContact] = useState("");
  const [copied, setCopied] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const proPrice =
    billingCycle === "monthly"
      ? BILLING_CONFIG.plans.pro.monthlyPrice
      : BILLING_CONFIG.plans.pro.yearlyPrice;

  // Standard UPI URI format accepted by Google Pay, PhonePe, Paytm, Cred, BHIM
  const upiUri = `upi://pay?pa=${BILLING_CONFIG.upiId}&pn=${encodeURIComponent(
    BILLING_CONFIG.merchantName
  )}&am=${proPrice}&cu=INR&tn=${encodeURIComponent(
    `SentinelDesk Pro ${billingCycle.toUpperCase()}`
  )}`;

  const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=220x220&margin=10&data=${encodeURIComponent(
    upiUri
  )}`;

  const handleCopyUpi = () => {
    navigator.clipboard.writeText(BILLING_CONFIG.upiId);
    setCopied(true);
    toast.success("UPI ID copied to clipboard!");
    setTimeout(() => setCopied(false), 2000);
  };

  const handleWhatsAppContact = (planName: string, priceText: string) => {
    const text = encodeURIComponent(
      `Hello! I want to upgrade my SentinelDesk license to *${planName} Plan* (${priceText}).\nCurrent active devices: ${currentDevices}.\nPlease provide activation details.`
    );
    const cleanPhone = BILLING_CONFIG.whatsappNumber.replace(/[^0-9]/g, "");
    window.open(`https://wa.me/${cleanPhone}?text=${text}`, "_blank");
  };

  const handleEmailContact = (planName: string, priceText: string) => {
    const subject = encodeURIComponent(
      `SentinelDesk License Upgrade Request - ${planName}`
    );
    const body = encodeURIComponent(
      `Hello SentinelDesk Sales Team,\n\nI want to upgrade my account to the ${planName} Plan (${priceText}).\nCurrently monitoring ${currentDevices} devices.\n\nPlease share the official invoice and license activation instructions.\n\nThank you!`
    );
    window.open(
      `mailto:${BILLING_CONFIG.supportEmail}?subject=${subject}&body=${body}`,
      "_blank"
    );
  };

  const handleConfirmPayment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!utrNumber.trim()) {
      toast.error("Please enter the 12-digit UTR or Transaction ID");
      return;
    }
    setSubmitting(true);
    setTimeout(() => {
      setSubmitting(false);
      setView("success");
      toast.success("Payment details submitted successfully!");
    }, 1000);
  };

  const handleClose = () => {
    onOpenChange(false);
    setTimeout(() => {
      setView("plans");
      setUtrNumber("");
      setCustomerContact("");
    }, 300);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-4xl w-[96vw] max-h-[88vh] overflow-y-auto overscroll-contain p-0 bg-background border border-border/70 shadow-2xl rounded-2xl">
        {/* VIEW 1: PRICING PLANS */}
        {view === "plans" && (
          <div className="p-4 sm:p-5 md:p-6 space-y-3.5">
            <DialogHeader className="text-center space-y-1">
              <div className="flex items-center justify-center gap-2">
                <span className="p-1.5 rounded-lg bg-primary/10 text-primary ring-1 ring-primary/20">
                  <Zap className="h-4 w-4" />
                </span>
                <DialogTitle className="text-xl sm:text-2xl font-extrabold tracking-tight">
                  Enterprise Endpoint <span className="text-primary">Licensing</span>
                </DialogTitle>
              </div>
              <DialogDescription className="text-xs max-w-md mx-auto text-muted-foreground">
                Secure more workstations, unlock the multi-screen wall, and get 24/7 dedicated support.
              </DialogDescription>

              {/* Billing Toggle (Monthly / Yearly) */}
              <div className="flex items-center justify-center pt-0.5">
                <div className="flex items-center bg-muted/60 p-0.5 rounded-full border border-border/80 text-xs font-semibold">
                  <button
                    type="button"
                    onClick={() => setBillingCycle("monthly")}
                    className={`px-3 py-1 rounded-full transition-all ${
                      billingCycle === "monthly"
                        ? "bg-background text-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    Monthly Billing
                  </button>
                  <button
                    type="button"
                    onClick={() => setBillingCycle("yearly")}
                    className={`px-3 py-1 rounded-full flex items-center gap-1.5 transition-all ${
                      billingCycle === "yearly"
                        ? "bg-background text-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    <span>Annual Billing</span>
                    <span className="bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 text-[10px] px-1.5 py-0.2 rounded-full font-bold">
                      SAVE 20%
                    </span>
                  </button>
                </div>
              </div>
            </DialogHeader>

            {/* Pricing Cards Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3.5">
              {/* Starter Plan */}
              <div className="rounded-2xl border bg-card/40 p-4 flex flex-col justify-between shadow-sm relative">
                <div>
                  <div className="flex justify-between items-center mb-1">
                    <h3 className="font-bold text-base">Starter</h3>
                    <Badge variant="secondary" className="text-[11px] px-2 py-0">
                      Current Plan
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground mb-2">
                    Basic endpoint monitoring for small setups.
                  </p>
                  <div className="mb-2">
                    <span className="text-2xl font-extrabold">₹0</span>
                    <span className="text-xs text-muted-foreground ml-1">/ forever</span>
                  </div>
                  <ul className="space-y-1.5 text-xs text-muted-foreground">
                    {BILLING_CONFIG.plans.starter.features.map((feat, i) => (
                      <li key={i} className="flex items-center gap-1.5">
                        <Check className="h-3.5 w-3.5 text-primary shrink-0" />
                        <span>{feat}</span>
                      </li>
                    ))}
                  </ul>
                </div>
                <Button variant="outline" className="w-full mt-3 text-xs h-8" disabled>
                  Active Plan
                </Button>
              </div>

              {/* Professional Plan */}
              <div className="rounded-2xl border-2 border-primary bg-gradient-to-b from-primary/5 via-card to-card p-4 flex flex-col justify-between shadow-xl relative ring-1 ring-primary/20">
                <div className="absolute -top-2.5 left-1/2 -translate-x-1/2">
                  <Badge className="bg-primary text-primary-foreground text-[10px] px-2.5 py-0 shadow font-semibold flex items-center gap-1">
                    <Sparkles className="h-2.5 w-2.5" /> MOST POPULAR
                  </Badge>
                </div>
                <div>
                  <div className="flex justify-between items-center mb-1">
                    <h3 className="font-bold text-base text-primary">Professional</h3>
                    <Badge className="bg-primary/10 text-primary border-primary/25 text-[11px] px-2 py-0">
                      25 Endpoints
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground mb-2">
                    Complete monitoring & control for corporate teams.
                  </p>
                  <div className="mb-2">
                    <span className="text-2xl font-extrabold">
                      ₹{proPrice.toLocaleString("en-IN")}
                    </span>
                    <span className="text-xs text-muted-foreground ml-1">
                      / {billingCycle === "monthly" ? "month" : "year"}
                    </span>
                  </div>
                  <ul className="space-y-1 text-xs">
                    {BILLING_CONFIG.plans.pro.features.map((feat, i) => (
                      <li key={i} className="flex items-center gap-1.5">
                        <Check className="h-3.5 w-3.5 text-primary shrink-0 font-bold" />
                        <span className="font-medium">{feat}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                {/* Primary & Direct Actions */}
                <div className="space-y-1.5 mt-3">
                  {/* Online Instant Payment */}
                  <Button
                    onClick={() => {
                      setSelectedPlan("pro");
                      setView("checkout");
                    }}
                    className="w-full text-xs font-bold h-8.5 shadow-md bg-gradient-to-r from-amber-500 to-primary hover:from-amber-600 hover:to-primary/90 text-white flex items-center justify-center gap-1.5"
                  >
                    <CreditCard className="h-3.5 w-3.5" /> Pay Online (Instant UPI / Cards)
                  </Button>

                  {/* Secondary WhatsApp & Email Options */}
                  <div className="grid grid-cols-2 gap-1.5">
                    <Button
                      variant="outline"
                      onClick={() =>
                        handleWhatsAppContact(
                          "Professional (25 Devices)",
                          `₹${proPrice}/${billingCycle}`
                        )
                      }
                      className="text-xs h-7.5 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-500/10 border-emerald-500/30 px-2"
                    >
                      <MessageSquare className="h-3.5 w-3.5 mr-1" /> WhatsApp
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() =>
                        handleEmailContact(
                          "Professional (25 Devices)",
                          `₹${proPrice}/${billingCycle}`
                        )
                      }
                      className="text-xs h-7.5 px-2"
                    >
                      <Mail className="h-3.5 w-3.5 mr-1" /> Email Invoice
                    </Button>
                  </div>
                </div>
              </div>

              {/* Enterprise Plan */}
              <div className="rounded-2xl border bg-card/40 p-4 flex flex-col justify-between shadow-sm">
                <div>
                  <div className="flex justify-between items-center mb-1">
                    <h3 className="font-bold text-base">Enterprise</h3>
                    <Badge variant="outline" className="text-[11px] px-2 py-0">
                      Custom
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground mb-2">
                    Dedicated server infrastructure & unlimited seats.
                  </p>
                  <div className="mb-2">
                    <span className="text-2xl font-extrabold">Custom</span>
                    <span className="text-xs text-muted-foreground ml-1">/ deployment</span>
                  </div>
                  <ul className="space-y-1.5 text-xs text-muted-foreground">
                    {BILLING_CONFIG.plans.enterprise.features.map((feat, i) => (
                      <li key={i} className="flex items-center gap-1.5">
                        <Check className="h-3.5 w-3.5 text-primary shrink-0" />
                        <span>{feat}</span>
                      </li>
                    ))}
                  </ul>
                </div>
                <div className="space-y-1.5 mt-3">
                  <Button
                    variant="outline"
                    onClick={() =>
                      handleWhatsAppContact(
                        "Enterprise (Unlimited)",
                        "Custom Corporate Pricing"
                      )
                    }
                    className="w-full text-xs flex items-center justify-center gap-1.5 font-semibold h-8.5"
                  >
                    <Server className="h-3.5 w-3.5" /> Talk to Enterprise Sales
                  </Button>
                  <Button
                    variant="ghost"
                    onClick={() =>
                      handleEmailContact(
                        "Enterprise (Unlimited)",
                        "Custom Corporate Pricing"
                      )
                    }
                    className="w-full text-xs h-7.5 text-muted-foreground"
                  >
                    <Mail className="h-3.5 w-3.5 mr-1" /> Contact via Email
                  </Button>
                </div>
              </div>
            </div>

            {/* Trust Footer */}
            <div className="text-center text-[11px] text-muted-foreground pt-2.5 border-t flex items-center justify-center gap-3">
              <span className="flex items-center gap-1 text-emerald-600 dark:text-emerald-400 font-medium">
                <ShieldCheck className="h-3.5 w-3.5" /> 100% Secure Transaction
              </span>
              <span>•</span>
              <span>Instant Activation</span>
              <span>•</span>
              <span>GST Invoice Provided</span>
            </div>
          </div>
        )}

        {/* VIEW 2: INSTANT CHECKOUT (UPI QR & GATEWAY) */}
        {view === "checkout" && (
          <div className="p-5 sm:p-6 space-y-4">
            <div className="flex items-center justify-between border-b pb-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setView("plans")}
                className="gap-1.5 text-xs font-semibold h-8"
              >
                <ArrowLeft className="h-4 w-4" /> Back to Plans
              </Button>
              <div className="flex items-center gap-1.5">
                <Lock className="h-3.5 w-3.5 text-emerald-500" />
                <span className="text-[11px] font-semibold text-muted-foreground">
                  256-Bit SSL Encrypted Checkout
                </span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5 items-start">
              {/* Left Column: QR Code & UPI Details */}
              <div className="flex flex-col items-center text-center p-4 rounded-2xl border bg-muted/20 space-y-3">
                <div className="space-y-0.5">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-primary">
                    Scan with any UPI App
                  </span>
                  <h4 className="text-sm font-bold">
                    GPay • PhonePe • Paytm • BHIM
                  </h4>
                </div>

                {/* Dynamic Live UPI QR Code */}
                <div className="bg-white p-2.5 rounded-2xl shadow-md border border-border">
                  <img
                    src={qrCodeUrl}
                    alt="UPI Payment QR Code"
                    className="w-40 h-40 rounded-xl object-contain"
                  />
                </div>

                <div className="w-full space-y-2">
                  <div className="flex items-center justify-between p-2.5 px-3 rounded-xl bg-card border text-xs">
                    <span className="text-muted-foreground">UPI ID:</span>
                    <span className="font-mono font-semibold">{BILLING_CONFIG.upiId}</span>
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={handleCopyUpi}
                      className="h-7 px-2 text-xs gap-1"
                    >
                      {copied ? (
                        <Check className="h-3.5 w-3.5 text-emerald-500" />
                      ) : (
                        <Copy className="h-3.5 w-3.5" />
                      )}
                      <span>{copied ? "Copied" : "Copy"}</span>
                    </Button>
                  </div>

                  <div className="text-[11px] text-muted-foreground flex justify-between px-1">
                    <span>Payee: {BILLING_CONFIG.merchantName}</span>
                    <span className="font-bold text-foreground">
                      Amount: ₹{proPrice.toLocaleString("en-IN")}
                    </span>
                  </div>
                </div>

                {/* Direct UPI Intent Button (Mobile & Desktop App) */}
                <a
                  href={upiUri}
                  className="w-full inline-flex items-center justify-center gap-2 h-9 rounded-xl bg-primary text-primary-foreground text-xs font-bold shadow hover:bg-primary/90"
                >
                  <ExternalLink className="h-3.5 w-3.5" /> Open in UPI App
                </a>

                {/* Optional Online Payment Gateway Link (PayUMoney / Cards) */}
                {BILLING_CONFIG.paymentLink && (
                  <a
                    href={BILLING_CONFIG.paymentLink}
                    target="_blank"
                    rel="noreferrer"
                    className="w-full inline-flex items-center justify-center gap-2 h-9 rounded-xl border border-primary/40 bg-primary/10 hover:bg-primary/20 text-primary text-xs font-bold shadow-sm"
                  >
                    <CreditCard className="h-3.5 w-3.5" /> Pay via Card / NetBanking
                  </a>
                )}
              </div>

              {/* Right Column: Order Summary & UTR Submission */}
              <div className="space-y-5">
                {/* Order Summary */}
                <div className="p-4 rounded-2xl border bg-card space-y-3">
                  <div className="flex justify-between items-center text-sm">
                    <span className="font-bold">SentinelDesk Professional</span>
                    <Badge className="text-xs">25 Endpoints</Badge>
                  </div>
                  <div className="text-xs text-muted-foreground space-y-1">
                    <div className="flex justify-between">
                      <span>Billing Term:</span>
                      <span className="font-medium text-foreground capitalize">
                        {billingCycle}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span>Live Multi-Screen Matrix:</span>
                      <span className="text-emerald-500 font-medium">Included</span>
                    </div>
                    <div className="flex justify-between border-t pt-2 text-sm">
                      <span className="font-bold">Total Payable:</span>
                      <span className="font-extrabold text-primary text-base">
                        ₹{proPrice.toLocaleString("en-IN")}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Verification Form */}
                <form onSubmit={handleConfirmPayment} className="space-y-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="utr" className="text-xs font-semibold">
                      UPI Reference / UTR Number (12 digits) *
                    </Label>
                    <Input
                      id="utr"
                      placeholder="e.g. 412356789012"
                      value={utrNumber}
                      onChange={(e) => setUtrNumber(e.target.value)}
                      required
                      className="font-mono text-xs"
                    />
                    <p className="text-[11px] text-muted-foreground">
                      Scan the QR above, pay from your app, and enter the UTR/Reference ID shown in your receipt.
                    </p>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="contact" className="text-xs font-semibold">
                      Your Email or WhatsApp Phone (For License Key) *
                    </Label>
                    <Input
                      id="contact"
                      placeholder="e.g. client@company.com or 9876543210"
                      value={customerContact}
                      onChange={(e) => setCustomerContact(e.target.value)}
                      required
                      className="text-xs"
                    />
                  </div>

                  <Button
                    type="submit"
                    disabled={submitting}
                    className="w-full text-xs font-bold h-10 bg-emerald-600 hover:bg-emerald-700 text-white"
                  >
                    {submitting ? "Verifying..." : "Confirm & Activate License"}
                  </Button>
                </form>

                {/* Instant WhatsApp Support */}
                <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-xs flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <MessageSquare className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
                    <span>Need instant manual verification?</span>
                  </div>
                  <Button
                    size="sm"
                    variant="link"
                    onClick={() =>
                      handleWhatsAppContact(
                        "Pro License UTR Verification",
                        `UTR: ${utrNumber || "Pending"}`
                      )
                    }
                    className="h-auto p-0 text-emerald-600 dark:text-emerald-400 font-bold"
                  >
                    Chat on WhatsApp →
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* VIEW 3: SUCCESS STATE */}
        {view === "success" && (
          <div className="p-8 sm:p-12 text-center space-y-5 max-w-lg mx-auto">
            <div className="w-16 h-16 rounded-full bg-emerald-500/10 text-emerald-500 flex items-center justify-center mx-auto ring-8 ring-emerald-500/5">
              <CheckCircle2 className="h-10 w-10" />
            </div>
            <div className="space-y-2">
              <h3 className="text-2xl font-bold">Payment Request Submitted!</h3>
              <p className="text-xs text-muted-foreground">
                Your payment reference has been recorded. Our automated system is verifying the UTR:
              </p>
              <div className="p-3 bg-muted rounded-xl font-mono text-xs font-bold text-primary">
                UTR: {utrNumber}
              </div>
              <p className="text-xs text-muted-foreground">
                Your plan quota will automatically increase to 25 devices upon receipt. You can also send the screenshot to our WhatsApp support for priority 2-minute activation.
              </p>
            </div>
            <div className="pt-2 flex flex-col sm:flex-row gap-2 justify-center">
              <Button
                onClick={() =>
                  handleWhatsAppContact(
                    "Payment Receipt Submission",
                    `UTR: ${utrNumber}`
                  )
                }
                className="text-xs font-bold bg-emerald-600 hover:bg-emerald-700 text-white gap-1.5"
              >
                <MessageSquare className="h-4 w-4" /> Send Receipt on WhatsApp
              </Button>
              <Button variant="outline" onClick={handleClose} className="text-xs">
                Back to Dashboard
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
