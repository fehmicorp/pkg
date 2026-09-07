import React from "react";
import { Terminal } from "lucide-react";

interface StepInstallingProps {
  progress: number;
  currentAction: string;
  logs: string[];
}

export function StepInstalling({ progress, currentAction, logs }: StepInstallingProps) {
  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-xl font-bold text-white">Installing Fehmi Agent</h1>
        <p className="text-xs text-slate-400">Please wait while setup deploys binaries and registers services...</p>
      </div>

      <div className="space-y-2">
        <div className="flex justify-between text-xs">
          <span className="text-slate-300 font-medium">{currentAction}</span>
          <span className="text-cyan-400 font-mono">{progress}%</span>
        </div>
        <div className="h-2.5 w-full overflow-hidden rounded-full bg-slate-950 border border-slate-800">
          <div
            className="h-full bg-gradient-to-r from-cyan-500 to-blue-600 transition-all duration-300 ease-out"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      <div className="rounded-lg border border-slate-800 bg-slate-950 p-3">
        <div className="flex items-center gap-2 border-b border-slate-800/80 pb-2 mb-2 text-[11px] font-mono text-slate-400">
          <Terminal className="h-3.5 w-3.5 text-cyan-400" />
          <span>Installation Log</span>
        </div>
        <div className="h-36 overflow-y-auto font-mono text-[11px] text-slate-300 space-y-1">
          {logs.map((log, index) => (
            <div key={index} className="leading-tight text-slate-400">
              {log}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}