import { useEffect, useState } from "react";
import { GetConfig, WindowMinimize, WindowToggleMaximize, WindowClose } from "../wailsjs/go/main/App";

interface AppConfig {
  AppName: string;
  Description: string;
  Icon: string;
  Version: string;
  Domain: string;
  installDir: string;
}

export function TitleBar() {
  const [config, setConfig] = useState<AppConfig | null>(null);

  useEffect(() => {
    if (GetConfig) {
      GetConfig()
        .then((cfg: AppConfig) => setConfig(cfg))
        .catch((err) => console.error("Failed to load app config:", err));
    }
  }, []);

  const displayTitle = config ? `${config.AppName} ${config.Version}` : "Fehmi Agent Installer";


  return (
    <div
      style={{ style: "drag" } as any}
      className="flex h-10 w-full items-center justify-between border-b border-slate-800 bg-slate-950/80 px-4 select-none"
    >
      <div className="flex items-center gap-2">
        <img src="images/logo.png" alt="Fehmi Logo" className="h-5 w-5 object-contain" />
        <span className="text-xs font-semibold tracking-wide text-slate-300">
          {displayTitle}
        </span>
      </div>
      <div
        style={{ style: "no-drag" } as any}
        className="flex items-center gap-2"
      >
        <button
          onClick={() => WindowMinimize?.()}
          title="Minimize"
          className="h-3 w-3 rounded-full bg-slate-700 hover:bg-yellow-500 transition-colors cursor-pointer border-0 outline-none"
        />
        <button
          onClick={() => WindowToggleMaximize?.()}
          title="Maximize"
          className="h-3 w-3 rounded-full bg-slate-700 hover:bg-green-500 transition-colors cursor-pointer border-0 outline-none"
        />
        <button
          onClick={() => WindowClose?.()}
          title="Close"
          className="h-3 w-3 rounded-full bg-slate-700 hover:bg-rose-500 transition-colors cursor-pointer border-0 outline-none"
        />
      </div>
    </div>
  );
}