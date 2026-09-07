import React from "react";
import { TitleBar } from "@/components/TitleBar";
import { useAppConfig } from "@/useAppConfig";
import { useInstaller } from "@/components/useInstaller";
import { Sidebar } from "@/components/Sidebar";

function App() {
  const appConfig = useAppConfig();
  const installer = useInstaller();

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-slate-950 font-sans text-slate-100 select-none">
      <div
        style={{
          width: `${appConfig.Width}px`,
          height: `${appConfig.Height}px`,
        }}
        className="flex flex-col overflow-hidden rounded-xl border border-slate-800 bg-slate-900 shadow-2xl shadow-cyan-950/20"
      >
        <TitleBar config={appConfig} />
        <div className="flex flex-1 overflow-hidden">
          <Sidebar step={installer.step} config={appConfig}/>
        </div>
      </div>
    </div>
  );
}

export default App;