import React from "react";
import { ShieldCheck, Monitor, CheckCircle2 } from "lucide-react";
import { InstallStep } from "./types";
import { AppConfig } from "@/useAppConfig";


interface SidebarProps {
  step: InstallStep;
  config: AppConfig;
}

export function Sidebar({ step, config }: SidebarProps) {
  return (
    <div className="flex w-60 flex-col justify-between border-r border-slate-800 bg-slate-950/40 p-6">
      <div>
        <div className="mb-8 flex items-center gap-3">
          <ShieldCheck className="h-8 w-8 text-cyan-400" />
          <div>
            <h2 className="text-sm font-bold text-white leading-tight">{config.Title}</h2>
            <p className="text-[11px] text-slate-400">{config.Tagline}</p>
          </div>
        </div>

        <div className="space-y-4">
          <StepItem
            active={step === "welcome"}
            completed={step === "config" || step === "installing" || step === "complete"}
            title="1. Welcome"
            subtitle="Overview & EULA"
          />
          <StepItem
            active={step === "config"}
            completed={step === "installing" || step === "complete"}
            title="2. Configuration"
            subtitle="Server & API Setup"
          />
          <StepItem
            active={step === "installing"}
            completed={step === "complete"}
            title="3. Installation"
            subtitle="Deploying Binaries"
          />
          <StepItem
            active={step === "complete"}
            completed={step === "complete"}
            title="4. Finish"
            subtitle="Complete Setup"
          />
        </div>
      </div>

      <div className="rounded-lg border border-slate-800/80 bg-slate-900/50 p-3 text-[11px] text-slate-400">
        <div className="flex items-center gap-2 mb-1 text-slate-300 font-medium">
          <Monitor className="h-3.5 w-3.5 text-cyan-400" /> Target OS
        </div>
        <span>Windows 10/11 / Server 2019+</span>
      </div>
    </div>
  );
}

function StepItem({
  title,
  subtitle,
  active,
  completed,
}: {
  title: string;
  subtitle: string;
  active: boolean;
  completed: boolean;
}) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5">
        {completed && !active ? (
          <CheckCircle2 className="h-4 w-4 text-emerald-400" />
        ) : active ? (
          <div className="h-4 w-4 rounded-full border-2 border-cyan-400 bg-cyan-950 flex items-center justify-center">
            <div className="h-1.5 w-1.5 rounded-full bg-cyan-400" />
          </div>
        ) : (
          <div className="h-4 w-4 rounded-full border border-slate-700 bg-slate-900" />
        )}
      </div>
      <div>
        <p
          className={`text-xs font-semibold leading-none ${
            active ? "text-cyan-400" : completed ? "text-slate-200" : "text-slate-500"
          }`}
        >
          {title}
        </p>
        <p className="text-[10px] text-slate-500 mt-0.5">{subtitle}</p>
      </div>
    </div>
  );
}