import React from "react";
import { ChevronRight, ChevronLeft } from "lucide-react";
import { InstallStep } from "./types";

interface NavigationProps {
  step: InstallStep;
  setStep: (step: InstallStep) => void;
  acceptedTos: boolean;
  apiKey: string;
}

export function Navigation({ step, setStep, acceptedTos, apiKey }: NavigationProps) {
  return (
    <div className="flex items-center justify-between border-t border-slate-800/80 pt-4 mt-auto">
      <div>
        {step !== "installing" && step !== "complete" && (
          <button
            onClick={() => setStep("welcome")}
            className="text-xs text-slate-500 hover:text-slate-300 transition-colors"
          >
            Cancel
          </button>
        )}
      </div>

      <div className="flex items-center gap-3">
        {step === "welcome" && (
          <button
            onClick={() => setStep("tos")}
            className="flex items-center gap-1.5 rounded-lg bg-cyan-600 px-5 py-2 text-xs font-medium text-white hover:bg-cyan-500 transition-all shadow-lg shadow-cyan-950/50"
          >
            Next <ChevronRight className="h-4 w-4" />
          </button>
        )}

        {step === "tos" && (
          <>
            <button
              onClick={() => setStep("welcome")}
              className="flex items-center gap-1.5 rounded-lg border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-medium text-slate-200 hover:bg-slate-700 transition-all"
            >
              <ChevronLeft className="h-4 w-4" /> Back
            </button>
            <button
              onClick={() => setStep("config")}
              disabled={!acceptedTos}
              className="flex items-center gap-1.5 rounded-lg bg-cyan-600 px-5 py-2 text-xs font-medium text-white hover:bg-cyan-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-cyan-950/50"
            >
              Accept & Continue <ChevronRight className="h-4 w-4" />
            </button>
          </>
        )}

        {step === "config" && (
          <>
            <button
              onClick={() => setStep("tos")}
              className="flex items-center gap-1.5 rounded-lg border border-slate-700 bg-slate-800 px-4 py-2 text-xs font-medium text-slate-200 hover:bg-slate-700 transition-all"
            >
              <ChevronLeft className="h-4 w-4" /> Back
            </button>
            <button
              onClick={() => setStep("installing")}
              disabled={!apiKey}
              className="flex items-center gap-1.5 rounded-lg bg-cyan-600 px-5 py-2 text-xs font-medium text-white hover:bg-cyan-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-cyan-950/50"
            >
              Install Now <ChevronRight className="h-4 w-4" />
            </button>
          </>
        )}

        {step === "complete" && (
          <button
            onClick={() => alert("Setup closed.")}
            className="flex items-center gap-1.5 rounded-lg bg-emerald-600 px-6 py-2 text-xs font-medium text-white hover:bg-emerald-500 transition-all shadow-lg shadow-emerald-950/50"
          >
            Finish
          </button>
        )}
      </div>
    </div>
  );
}