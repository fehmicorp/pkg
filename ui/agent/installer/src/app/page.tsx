"use client";

import React from "react";
import { useInstaller } from "@/components/useInstaller";
import { TitleBar } from "@/components/TitleBar";
import { Sidebar } from "@/components/Sidebar";
import { StepWelcome } from "@/components/StepWelcome";
import { StepTos } from "@/components/StepTos";
import { StepConfig } from "@/components/StepConfig";
import { StepInstalling } from "@/components/StepInstalling";
import { StepComplete } from "@/components/StepComplete";
import { Navigation } from "@/components/Navigation";

export default function FehmiAgentInstaller() {
  const installer = useInstaller();

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-slate-950 p-4 font-sans text-slate-100 select-none">
      <div className="flex h-[560px] w-[820px] flex-col overflow-hidden rounded-xl border border-slate-800 bg-slate-900 shadow-2xl shadow-cyan-950/20">
        <TitleBar />

        <div className="flex flex-1 overflow-hidden">
          <Sidebar step={installer.step} />

          <div className="flex flex-1 flex-col justify-between p-8 bg-slate-900/60">
            {installer.step === "welcome" && <StepWelcome />}
            {installer.step === "tos" && (
              <StepTos
                accepted={installer.acceptedTos}
                setAccepted={installer.setAcceptedTos}
              />
            )}
            {installer.step === "config" && (
              <StepConfig
                serverUrl={installer.serverUrl}
                setServerUrl={installer.setServerUrl}
                apiKey={installer.apiKey}
                setApiKey={installer.setApiKey}
                installPath={installer.installPath}
                setInstallPath={installer.setInstallPath}
                installAsService={installer.installAsService}
                setInstallAsService={installer.setInstallAsService}
                enableAutostart={installer.enableAutostart}
                setEnableAutostart={installer.setEnableAutostart}
              />
            )}
            {installer.step === "installing" && (
              <StepInstalling
                progress={installer.progress}
                currentAction={installer.currentAction}
                logs={installer.logs}
              />
            )}
            {installer.step === "complete" && <StepComplete />}

            <Navigation
              step={installer.step}
              setStep={installer.setStep}
              acceptedTos={installer.acceptedTos}
              apiKey={installer.apiKey}
            />
          </div>
        </div>
      </div>
    </div>
  );
}