/**
 * SentinelDesk Billing & Payment Gateway Configuration
 * Edit your WhatsApp number, UPI ID, and Email here.
 */
export const BILLING_CONFIG = {
  merchantName: process.env.NEXT_PUBLIC_MERCHANT_NAME || "SentinelDesk Security",
  
  // Enter your WhatsApp number with country code (e.g., "919876543210")
  whatsappNumber: process.env.NEXT_PUBLIC_WHATSAPP_NUMBER || "919999999999",

  // Enter your official support/sales email
  supportEmail: process.env.NEXT_PUBLIC_SUPPORT_EMAIL || "parastomar@recenturesoft.com",

  // Enter your UPI ID for receiving payments (e.g., "name@okaxis", "name@upi", "name@paytm")
  upiId: process.env.NEXT_PUBLIC_UPI_ID || "parastomar@upi",

  // Optional direct payment link (e.g. Razorpay Payment Page link / Cashfree / Stripe link)
  paymentLink: process.env.NEXT_PUBLIC_PAYMENT_LINK || "",

  plans: {
    starter: {
      name: "Starter",
      badge: "Free Tier",
      devices: 3,
      monthlyPrice: 0,
      yearlyPrice: 0,
      features: [
        "Up to 3 Active Endpoints",
        "Real-time Telemetry & CPU/RAM",
        "Process & Windows Service Monitor",
        "Single Screen Live View",
        "Community Support",
      ],
    },
    pro: {
      name: "Professional",
      badge: "Recommended ⭐",
      devices: 25,
      monthlyPrice: 1499,
      yearlyPrice: 14990, // 2 months free
      features: [
        "Up to 25 Monitored Endpoints",
        "Multi-Screen Matrix Wall (6 Displays)",
        "Ultra-low Latency Live Screen Stream",
        "Remote Windows Service Restart/Kill",
        "Software Inventory & Audit Tracking",
        "Instant WhatsApp & Priority Email Support",
        "99.9% Uptime SLA",
      ],
    },
    enterprise: {
      name: "Enterprise",
      badge: "Corporate",
      devices: "Unlimited",
      monthlyPrice: 4999,
      features: [
        "Unlimited Endpoint Monitoring",
        "Dedicated Server & Database Isolation",
        "Custom Domain & White-label Branding",
        "Full Security Audit Logs & Retention",
        "Custom Endpoint Deployment Package",
        "24/7 Dedicated Account Manager",
      ],
    },
  },
};
