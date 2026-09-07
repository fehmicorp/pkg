import React from "react";
import { AppConfig } from "@/useAppConfig";
import { ActionButtons } from "./ActionButton";
import logo from "../assets/images/logo.png";

interface TitleProps {
  config: AppConfig;
}

export function TitleBar({ config }: TitleProps) {
  return (
    <div
      style={{ "--wails-draggable": "drag" } as React.CSSProperties}
      className="flex h-10 w-full items-center justify-between border-b border-slate-800 px-4 bg-slate-950/80 select-none"
    >
      <div className="flex items-center gap-2">
        <img src={logo} alt="Logo" className="h-5 w-5 object-contain" />
        <span className="text-xs font-semibold tracking-wide text-slate-300">
          {config.AppName}
        </span>
      </div>

      <div
        style={{ "--wails-draggable": "no-drag" } as React.CSSProperties}
        className="flex items-center gap-2"
      >
        <ActionButtons
          allowMinimize={true}
          allowMaximize={false}
          allowQuit={true}
        />
      </div>
    </div>
  );
}