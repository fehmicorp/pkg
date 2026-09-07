import React from "react";

export function TitleBar() {
  return (
    <div className="flex h-10 w-full items-center justify-between border-b border-slate-800 bg-slate-950/80 px-4 select-none">
      <div className="flex items-center gap-2">
        <div className="flex h-5 w-5 items-center justify-center rounded bg-cyan-600 font-mono text-xs font-bold text-white">
          F
        </div>
        <span className="text-xs font-semibold tracking-wide text-slate-300">
          Fehmi Agent Setup v2.4.0 (x64)
        </span>
      </div>
      <div className="flex items-center gap-2">
        <div className="h-3 w-3 rounded-full bg-slate-700 hover:bg-yellow-500 transition-colors cursor-pointer" />
        <div className="h-3 w-3 rounded-full bg-slate-700 hover:bg-green-500 transition-colors cursor-pointer" />
        <div className="h-3 w-3 rounded-full bg-slate-700 hover:bg-rose-500 transition-colors cursor-pointer" />
      </div>
    </div>
  );
}