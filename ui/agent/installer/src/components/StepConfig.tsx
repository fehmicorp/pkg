import React from "react";
import { Server, Key, Folder } from "lucide-react";

interface StepConfigProps {
  serverUrl: string;
  setServerUrl: (val: string) => void;
  apiKey: string;
  setApiKey: (val: string) => void;
  installPath: string;
  setInstallPath: (val: string) => void;
  installAsService: boolean;
  setInstallAsService: (val: boolean) => void;
  enableAutostart: boolean;
  setEnableAutostart: (val: boolean) => void;
}

export function StepConfig({
  serverUrl,
  setServerUrl,
  apiKey,
  setApiKey,
  installPath,
  setInstallPath,
  installAsService,
  setInstallAsService,
  enableAutostart,
  setEnableAutostart,
}: StepConfigProps) {
  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-xl font-bold text-white">Agent Configuration</h1>
        <p className="text-xs text-slate-400">Configure connection settings to join your infrastructure pool.</p>
      </div>

      <div className="space-y-3">
        <div>
          <label className="block text-xs font-medium text-slate-300 mb-1">Gateway Endpoint URL</label>
          <div className="flex items-center rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 focus-within:border-cyan-500">
            <Server className="h-4 w-4 text-slate-400 mr-2" />
            <input
              type="text"
              value={serverUrl}
              onChange={(e) => setServerUrl(e.target.value)}
              className="w-full bg-transparent text-xs text-white outline-none"
            />
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium text-slate-300 mb-1">Enrollment Key / Auth Token</label>
          <div className="flex items-center rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 focus-within:border-cyan-500">
            <Key className="h-4 w-4 text-slate-400 mr-2" />
            <input
              type="password"
              placeholder="fhm_live_xxxxxxxxxxxxxxxx"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              className="w-full bg-transparent text-xs text-white outline-none"
            />
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium text-slate-300 mb-1">Installation Path</label>
          <div className="flex items-center rounded-lg border border-slate-700 bg-slate-950 px-3 py-2">
            <Folder className="h-4 w-4 text-slate-400 mr-2" />
            <input
              type="text"
              value={installPath}
              onChange={(e) => setInstallPath(e.target.value)}
              className="w-full bg-transparent text-xs text-slate-300 outline-none"
            />
          </div>
        </div>

        <div className="pt-1 space-y-2">
          <label className="flex items-center gap-2 text-xs text-slate-300 cursor-pointer">
            <input
              type="checkbox"
              checked={installAsService}
              onChange={(e) => setInstallAsService(e.target.checked)}
              className="rounded border-slate-700 bg-slate-950 text-cyan-600 focus:ring-0"
            />
            Install as background Windows Service (<code className="text-cyan-400">fehmi-agentd</code>)
          </label>
          <label className="flex items-center gap-2 text-xs text-slate-300 cursor-pointer">
            <input
              type="checkbox"
              checked={enableAutostart}
              onChange={(e) => setEnableAutostart(e.target.checked)}
              className="rounded border-slate-700 bg-slate-950 text-cyan-600 focus:ring-0"
            />
            Enable automatic startup on Windows boot
          </label>
        </div>
      </div>
    </div>
  );
}