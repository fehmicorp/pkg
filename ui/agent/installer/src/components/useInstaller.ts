"use client";

import { useState, useEffect } from "react";
import { InstallStep } from "./types";

export function useInstaller() {
  const [step, setStep] = useState<InstallStep>("welcome");
  const [acceptedTos, setAcceptedTos] = useState(false);
  const [installPath, setInstallPath] = useState("C:\\Program Files\\FehmiCorp\\FehmiAgent");
  const [serverUrl, setServerUrl] = useState("https://agent-gateway.fehmi.cloud");
  const [apiKey, setApiKey] = useState("");
  const [installAsService, setInstallAsService] = useState(true);
  const [enableAutostart, setEnableAutostart] = useState(true);

  const [progress, setProgress] = useState(0);
  const [currentAction, setCurrentAction] = useState("Initializing installation setup...");
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    if (step === "installing") {
      const steps = [
        { pct: 15, msg: "Creating directory structure...", log: "Directory created: C:\\Program Files\\FehmiCorp\\FehmiAgent" },
        { pct: 35, msg: "Extracting Fehmi Agent binaries...", log: "Unpacking fehmi-agent-win-x64.exe..." },
        { pct: 55, msg: "Writing configuration files...", log: "Generated config.yaml with endpoint and secure credentials." },
        { pct: 75, msg: "Registering Windows Service...", log: "sc.exe create FehmiAgent Service binPath= fehmi-agent.exe start= auto" },
        { pct: 90, msg: "Starting Fehmi Agent Service...", log: "Service 'FehmiAgent' started successfully." },
        { pct: 100, msg: "Installation completed successfully!", log: "Handshake verified with https://agent-gateway.fehmi.cloud" },
      ];

      let currentIndex = 0;
      const interval = setInterval(() => {
        if (currentIndex < steps.length) {
          const item = steps[currentIndex];
          setProgress(item.pct);
          setCurrentAction(item.msg);
          setLogs((prev) => [...prev, `[${new Date().toLocaleTimeString()}] ${item.log}`]);
          currentIndex++;
        } else {
          clearInterval(interval);
          setTimeout(() => setStep("complete"), 800);
        }
      }, 1000);

      return () => clearInterval(interval);
    }
  }, [step]);

  return {
    step,
    setStep,
    acceptedTos,
    setAcceptedTos,
    installPath,
    setInstallPath,
    serverUrl,
    setServerUrl,
    apiKey,
    setApiKey,
    installAsService,
    setInstallAsService,
    enableAutostart,
    setEnableAutostart,
    progress,
    currentAction,
    logs,
  };
}