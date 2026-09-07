import React from "react";
import { CheckCircle2 } from "lucide-react";

export function StepWelcome() {
  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold text-white">Install Fehmi Agent</h1>
        <p className="mt-1 text-sm text-slate-400">
          Deploy system monitoring, telemetry, and infrastructure automation services on this endpoint.
        </p>
      </div>

      <div className="rounded-lg border border-slate-800 bg-slate-950/60 p-4 space-y-3">
        <h3 className="text-xs font-semibold uppercase tracking-wider text-cyan-400">Prerequisites Check</h3>
        <div className="grid grid-cols-2 gap-3 text-xs">
          <div className="flex items-center gap-2 text-slate-300">
            <CheckCircle2 className="h-4 w-4 text-emerald-400" /> Architecture: x86_64
          </div>
          <div className="flex items-center gap-2 text-slate-300">
            <CheckCircle2 className="h-4 w-4 text-emerald-400" /> Administrator Rights
          </div>
          <div className="flex items-center gap-2 text-slate-300">
            <CheckCircle2 className="h-4 w-4 text-emerald-400" /> PowerShell 5.1+
          </div>
          <div className="flex items-center gap-2 text-slate-300">
            <CheckCircle2 className="h-4 w-4 text-emerald-400" /> Port 443 Outbound
          </div>
        </div>
      </div>

      <div className="rounded-lg border border-slate-800/60 bg-slate-950/30 p-4 text-xs text-slate-400 leading-relaxed">
        By clicking <strong className="text-slate-200">Next</strong>, you agree to allow Fehmi Agent to register as a system service for remote orchestration, diagnostic telemetry, and automated patch distribution.
      </div>
    </div>
  );
}