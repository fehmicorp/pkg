import React from "react";
import { CheckCircle2 } from "lucide-react";

export function StepComplete() {
  return (
    <div className="flex flex-col items-center justify-center text-center space-y-4 my-auto">
      <div className="flex h-16 w-16 items-center justify-center rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
        <CheckCircle2 className="h-10 w-10" />
      </div>

      <div>
        <h1 className="text-2xl font-bold text-white">Installation Complete</h1>
        <p className="mt-1 text-xs text-slate-400 max-w-md">
          Fehmi Agent has been successfully installed and registered as a background service on this machine.
        </p>
      </div>

      <div className="w-full max-w-sm rounded-lg border border-slate-800 bg-slate-950/60 p-4 text-left text-xs space-y-2">
        <div className="flex justify-between border-b border-slate-800/60 pb-1.5">
          <span className="text-slate-400">Status</span>
          <span className="text-emerald-400 font-medium flex items-center gap-1">
            <span className="h-1.5 w-1.5 rounded-full bg-emerald-400 animate-pulse" /> Running
          </span>
        </div>
        <div className="flex justify-between border-b border-slate-800/60 pb-1.5">
          <span className="text-slate-400">Gateway</span>
          <span className="text-slate-200 font-mono text-[11px]">fehmi.cloud</span>
        </div>
        <div className="flex justify-between">
          <span className="text-slate-400">Service</span>
          <span className="text-slate-200 font-mono text-[11px]">FehmiAgentSvc</span>
        </div>
      </div>
    </div>
  );
}